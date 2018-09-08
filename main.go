package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gplume/todo-list/pkg/engine"
	"github.com/gplume/todo-list/pkg/router"
)

func main() {
	fmt.Println("*****************************************************************")
	fmt.Println("*****************  Welcome to ToDoList API!  ********************")
	fmt.Println("*****************************************************************")

	// gin.DisableConsoleColor()
	// gin.DefaultWriter = colorable.NewColorableStdout() // for windows git bash especially

	var err error
	_, err = engine.NewApp(false)
	if err != nil {
		log.Fatalf("Error initializing Application : %v", err)
	}
	defer engine.App.Datamapper.Closing()

	if engine.App.Cfg.UsageMode != "dev" {
		fmt.Println("******************  APP IN PRODUCTION MODE  *********************")
	} else {
		fmt.Println("*******************  app in dev mode  ***************************")

	}

	switch engine.App.Cfg.FrameWorkMode {
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

	bind := fmt.Sprintf(":%d", engine.App.Cfg.ServerPort)

	srv := &http.Server{
		Addr:    bind,
		Handler: router.Engine,
		// ReadTimeout:    10 * time.Second,
		// WriteTimeout:   10 * time.Second,
		// MaxHeaderBytes: 1 << 20,
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
