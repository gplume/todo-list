package main

import (
	"os"
	"testing"

	"github.com/gplume/todo-list/pkg/engine"
)

func TestInitTestAppEngine(t *testing.T) {
	var err error
	_, err = engine.NewApp(false)
	if err != nil {
		t.Fatalf("Error initializing Application : %v", err)
	}
	engine.App.Datamapper.Closing()
	_, err = engine.NewApp(true)
	if err != nil {
		t.Fatalf("Error initializing Application : %v", err)
	}
	engine.App.Datamapper.Closing()
}

// func initTestAppEngine() {
// 	var err error
// 	_, err = engine.NewApp(true)
// 	if err != nil {
// 		log.Fatalf("Error initializing Application : %v", err)
// 	}
// 	engine.App.Datamapper.Closing()
// }

// TestMain is a trick to initialize app before running tests
func TestMain(m *testing.M) {
	runTests := m.Run()
	// engine.App.Datamapper.Closing()
	// switch engine.App.Cfg.DBType {
	// case "bolt":
	// 	err := os.Remove(fmt.Sprintf("%s/%s", engine.App.Cfg.DBDirectory, engine.App.Cfg.DBTestName))
	// 	if err != nil {
	// 		log.Fatalf("could not remove bolt database file properly: %v", err)
	// 	}
	// }
	os.Exit(runTests)
}
