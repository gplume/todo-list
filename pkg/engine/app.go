package engine

import (
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"path"
	"path/filepath"

	"github.com/boltdb/bolt"
	"github.com/gin-gonic/gin"
	"github.com/gplume/todo-list/pkg/config"
	datamapper "github.com/gplume/todo-list/pkg/datamappers"
	"github.com/gplume/todo-list/pkg/datamappers/boltmapper"
	"github.com/gplume/todo-list/pkg/utils"
	"github.com/joho/godotenv"
)

// App contains application structure
var App *Application

// Application contains application main variables
type Application struct {
	Cfg        *config.Config
	Datamapper datamapper.DataMapper
}

// NewApp sets *Application structure
func NewApp(testing bool) (*Application, error) {

	var err error
	app := &Application{}
	configDir := "conf"
	confPath := configDir + "/config.env"
	if testing {
		exPath := utils.GetDefaultBaseDir()
		for !utils.SearchDir(exPath, configDir) {
			var file string
			exPath, file = filepath.Split(exPath)
			if file == "" {
				return nil, err
			}
			exPath = filepath.Clean(exPath)
		}
		confPath = filepath.Join(exPath, confPath)
		confPath = filepath.Clean(confPath)
	}
	err = godotenv.Load(confPath)
	if err != nil {
		return nil, fmt.Errorf("error loading config.env file: %v", err)
	}
	// log.Println("---> ENV config OK")

	app.Cfg, err = config.New()
	if err != nil {
		return app, fmt.Errorf("error in config.env: %v", err)
	}

	// logs: MultiWriter to Stout and file
	logFile, _ := os.Create("./logs/server.log")
	gin.DefaultWriter = io.MultiWriter(logFile, os.Stdout)

	// DB
	databaseName := app.Cfg.DBName
	// databaseDir := app.Cfg.DBDirectory
	if testing {
		databaseName = app.Cfg.DBTestName

	}

	switch app.Cfg.DBType {
	// there is a switch here if one should try another database system
	// the switch is done via the config.env file (DB_TYPE) so only at the application start for now
	// for example add "pgsql" and write all the corresponding methods in a datamapper_pgsql.go file
	case "bolt":
		var db *bolt.DB
		var err error
		switch testing {
		case true:
			// init boltDB
			tmpfile, err := ioutil.TempFile("", databaseName)
			if err != nil {
				return nil, fmt.Errorf("%s", err)
			}
			defer os.Remove(tmpfile.Name())
			db, err = bolt.Open(tmpfile.Name(), 0660, nil)
			if err != nil {
				return nil, fmt.Errorf("%s", err)
			}
		default:
			if err := os.MkdirAll(path.Join(utils.GetDefaultBaseDir(), app.Cfg.DBDirectory), 0777); err != nil && !os.IsExist(err) {
				return nil, err
			}
			db, err = bolt.Open(fmt.Sprintf("%s/%s", app.Cfg.DBDirectory, databaseName), 0660, nil)
			if err != nil {
				return nil, err
			}
		}
		app.Datamapper, err = boltmapper.NewBoltDatamapper(db)
		if err != nil {
			log.Fatalf("ERROR INITIALIZING BoltDataMapper: %v", err)
		}
	default:
		return nil, errors.New("wrong database type provided in config.env file")
	}

	// Finally
	App = app
	return app, nil
}
