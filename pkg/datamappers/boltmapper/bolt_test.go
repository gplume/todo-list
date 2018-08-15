package boltmapper

import (
	"io/ioutil"
	"log"
	"os"
	"testing"
	"time"

	"github.com/boltdb/bolt"
	datamapper "github.com/gplume/todo-list/pkg/datamappers"
	"github.com/gplume/todo-list/pkg/models"
)

func Test_boltDB_SaveTodo(t *testing.T) {
	// int boltDB
	tmpfile, err := ioutil.TempFile("", "tests")
	if err != nil {
		t.Errorf("%s", err)
	}
	defer os.Remove(tmpfile.Name())
	bdb, err := bolt.Open(tmpfile.Name(), 0660, nil)
	if err != nil {
		t.Errorf("%s", err)
	}
	defer bdb.Close()
	dtmp, err := NewBoltDatamapper(bdb)
	if err != nil {
		t.Errorf("%s", err)
	}

	// testing
	type args struct {
		td *models.Todo
	}
	tests := []struct {
		name    string
		db      datamapper.DataMapper
		args    args
		wantErr bool
	}{
		{
			name: "test1_true",
			db:   dtmp,
			args: args{
				td: &models.Todo{
					Title:       "First task title",
					Description: "First Task description",
					Deadline:    time.Now().AddDate(0, 0, 1),
					Priority:    models.Low,
				}},
			wantErr: false,
		},
		{
			name: "test2_true",
			db:   dtmp,
			args: args{
				td: &models.Todo{
					Title:       "Second task title",
					Description: "Second Task description",
					Deadline:    time.Now().AddDate(0, 0, 1),
					Priority:    models.High,
				}},
			wantErr: false,
		},
		{name: "test3_false", db: dtmp, args: args{td: nil}, wantErr: true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			log.Printf("%+v", tt.args.td)
			if err := tt.db.SaveTodo(tt.args.td); (err != nil) != tt.wantErr {
				t.Errorf("boltDB.SaveTodo() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
