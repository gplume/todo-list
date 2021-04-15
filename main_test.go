package main

import (
	"fmt"
	"log"
	"os"
	"path"
	"testing"
	"time"

	"github.com/gplume/todo-list/pkg/engine"
)

func initTestAppEngine() error {
	var err error
	_, err = engine.NewApp(true)
	if err != nil {
		return fmt.Errorf("Error initializing Application : %v", err)
	}
	return nil
}

// TestMain is a trick to initialize app before running tests
func TestMain(m *testing.M) {
	if err := initTestAppEngine(); err != nil {
		log.Fatalf("could not initialize Tas App engine: %v", err)
	}

	runTests := m.Run()
	engine.App.Datamapper.Closing()
	time.Sleep(3 * time.Second) // because it's so fast it missed the db file creation in fs
	switch engine.App.Cfg.DBType {
	case "bolt":
		if err := os.Remove(path.Join(engine.App.Cfg.BaseDir, fmt.Sprintf("%s/%s", engine.App.Cfg.DBDirectory, engine.App.Cfg.DBTestName))); err != nil {
			log.Fatalf("could not remove bolt database file properly: %v", err)
		}
	}
	os.Exit(runTests)
}
