package main

import (
	"fmt"
	"log"
	"os"
	"testing"
)

func getTestAppEngine() *application {
	var err error
	app, err = newApp(true)
	if err != nil {
		log.Fatalf("Error initializing Application : %v", err)
	}
	router := mainEngineAndRoutes()
	app.router = router
	return app
}

func TestMain(m *testing.M) {
	app := getTestAppEngine()
	runTests := m.Run()
	app.datamapper.close()
	switch app.cfg.DBType {
	case "bolt":
		err := os.Remove(fmt.Sprintf("%s/%s", app.cfg.DBDirectory, app.cfg.DBTestName))
		if err != nil {
			log.Fatalf("could not remove bolt database file properly: %v", err)
		}
	}
	os.Exit(runTests)
}
