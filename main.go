package main

import (
	"context"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/signal"
	"path"
	"time"

	"github.com/boltdb/bolt"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var app *application

type application struct {
	cfg        *config
	router     *gin.Engine
	datamapper dataMapper
	prom       *promeVars
}

func newApp(testing bool) (*application, error) {

	var err error
	app = &application{}

	err = godotenv.Load("config/config.env")
	if err != nil {
		return nil, fmt.Errorf("error loading config.env file: %v", err)
	}

	app.cfg, err = newConfig()
	if err != nil {
		return app, fmt.Errorf("error in config.env: %v", err)
	}

	if err := os.MkdirAll(path.Join(getDefaultBaseDir(), app.cfg.DBDirectory), 0777); err != nil && !os.IsExist(err) {
		return nil, err
	}
	databaseName := app.cfg.DBName
	if testing {
		databaseName = app.cfg.DBTestName
	}

	switch app.cfg.DBType {
	// there is a switch here if one should try another database system
	// the switch is done via the config.env file (DB_TYPE) so only at the application start for now
	// for example add "pgsql" and write all the corresponding methods in a datamapper_pgsql.go file
	case "bolt":
		var db *bolt.DB
		db, err = bolt.Open(fmt.Sprintf("%s/%s", app.cfg.DBDirectory, databaseName), 0660, nil)
		if err != nil {
			return nil, err
		}
		app.datamapper, err = newBoltDatamapper(db, app.cfg)
		if err != nil {
			log.Fatalf("ERROR INITIALIZING BoltDataMapper: %v", err)
		}
	default:
		return nil, errors.New("wrong database type provided in config.env file")
	}

	app.prom = initPromeVars()

	return app, nil
}

func mainEngineAndRoutes() *gin.Engine {
	// Creates a router without any middleware by default
	r := gin.New()

	// Logger middleware will write the logs to gin.DefaultWriter even if you set with GIN_MODE=release.
	// By default gin.DefaultWriter = os.Stdout
	r.Use(gin.Logger())

	// Recovery middleware recovers from any panics and writes a 500 if there was one.
	r.Use(gin.Recovery())

	// group (will add /todo to all routes endpoints below:
	api := r.Group("/todo")

	// statsMiddleWare() used by prometheus to stats requests latencies
	api.Use(statsMiddleWare())

	// API endpoints (group suffix is added automatically)
	api.GET("", listTodos)
	api.GET("/:id", getTodo)
	api.POST("", addTodo)
	api.PUT("", updateTodo)
	api.DELETE("/:id", deleteTodo)

	// Prometheus metrics by convention:
	r.GET("/metrics", gin.WrapH(promhttp.Handler()))

	return r
}

func main() {
	fmt.Println("*****************************************************************")
	fmt.Println("*****************  Welcome to ToDoList API!  ********************")
	fmt.Println("*****************************************************************")

	// gin.DisableConsoleColor()
	// gin.DefaultWriter = colorable.NewColorableStdout() // for windows git bash especially

	// logs: MultiWriter to Stout and file
	logFile, _ := os.Create("./logs/server.log")
	gin.DefaultWriter = io.MultiWriter(logFile, os.Stdout)

	var err error
	app, err = newApp(false)
	if err != nil {
		log.Fatalf("Error initializing Application : %v", err)
	}
	defer app.datamapper.close()

	if app.cfg.UsageMode != "dev" {
		fmt.Println("******************  APP IN PRODUCTION MODE  *********************")
	} else {
		fmt.Println("*******************  app in dev mode  ***************************")

	}

	switch app.cfg.FrameWorkMode {
	case "release":
		gin.SetMode(gin.ReleaseMode)
		fmt.Println("----> gin framework is in RELEASE MODE")
	case "debug":
		gin.SetMode(gin.DebugMode)
		fmt.Println("----> gin framework is in DEBUG mode")
	case "test":
		gin.SetMode(gin.TestMode)
		fmt.Println("----> gin framework is in TEST MODE")
	default:
		gin.SetMode(gin.DebugMode)
		fmt.Println("----> gin framework is in DEBUG MODE")
	}

	registerPrometheusVars()   // placed here so metrics won't be incremented by tests
	r := mainEngineAndRoutes() // gets router and endpoints
	app.router = r

	bind := fmt.Sprintf(":%d", app.cfg.ServerPort)
	log.Printf("----> API is Running on: %s ", bind)

	switch app.cfg.SSLEnabled {
	case true:
		err = r.RunTLS(bind, app.cfg.SSLPub, app.cfg.SSLKey)
		if err != nil {
			log.Fatal("RunTLS: ", err)
		}
	default:
		srv := &http.Server{
			Addr:           bind,
			Handler:        r,
			ReadTimeout:    10 * time.Second,
			WriteTimeout:   10 * time.Second,
			MaxHeaderBytes: 1 << 20,
		}

		go func() {
			// service connections
			if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
				log.Fatalf("listen: %s\n", err)
			}
		}()

		// Wait for interrupt signal to gracefully shutdown the server with
		// a timeout of 5 seconds.
		quit := make(chan os.Signal)
		signal.Notify(quit, os.Interrupt)
		<-quit
		log.Println("Shutdown Server ...")

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		if err := srv.Shutdown(ctx); err != nil {
			log.Fatal("Server Shutdown:", err)
		}
		log.Println("Server exiting")
	}
}
