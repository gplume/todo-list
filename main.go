package main

import (
	"errors"
	"fmt"
	"log"
	"os"
	"path"

	"github.com/boltdb/bolt"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

var app *application

type application struct {
	cfg        *config
	datamapper dataMapper
}

func newApp() (*application, error) {

	var err error
	app = &application{}

	app.cfg, err = newConfig()
	if err != nil {
		return app, fmt.Errorf("error in config.env: %v", err)
	}

	if err := os.MkdirAll(path.Join(getDefaultBaseDir(), app.cfg.DBDirectory), 0777); err != nil && !os.IsExist(err) {
		return nil, err
	}

	switch app.cfg.DBType {
	// there is a switch here if one should try another database system
	// the switch is done via the config.env file (DB_TYPE) so only at the application start for now
	// for example add "pgsql" here and write all the corresponding methods in a pgsl_datamapper file
	case "bolt":
		var db *bolt.DB
		db, err = bolt.Open(fmt.Sprintf("%s/%s", app.cfg.DBDirectory, app.cfg.DBName), 0660, nil)
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

	return app, nil
}

func main() {
	fmt.Println("*****************************************************************")
	fmt.Println("*****************  Welcome to ToDoList API!  ********************")
	fmt.Println("*****************************************************************")

	gin.DisableConsoleColor()
	// gin.DefaultWriter = colorable.NewColorableStdout() // for windows git bash especially

	err := godotenv.Load("config/config.env")
	if err != nil {
		log.Fatalf("Error loading config.env file: %v", err)
	}

	app, err = newApp()
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

	r := gin.Default()

	bind := fmt.Sprintf(":%d", app.cfg.ServerPort)

	if app.cfg.SSLEnabled {
		err = r.RunTLS(bind, app.cfg.SSLPub, app.cfg.SSLKey)
		if err != nil {
			log.Fatal("ListenAndServe: ", err)
		}
	} else {
		r.Run(bind)
	}
}
