package main

import (
	"fmt"
	"testing"
	"time"
)

func Test_boltDB_getTodo(t *testing.T) {
	type args struct {
		todoKey int
	}
	tests := []struct {
		name    string
		db      *boltDB
		args    args
		want    *todo
		wantErr bool
	}{
		{"todo1", app.datamapper.db().(*boltDB), args{1}, &todo{Title: "First Task title", Deadline: time.Now().AddDate(0, 0, 1), Priority: 1}, false},
		{"todo1", app.datamapper.db().(*boltDB), args{99}, nil, true},
	}
	if err := app.datamapper.saveTodo(tests[0].want); err != nil {
		t.Errorf("boldDB.saveTodo() error = %v", err)
	}
	tests[0].args.todoKey = tests[0].want.ID

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.db.getTodo(tt.args.todoKey)
			if (err != nil) != tt.wantErr {
				t.Errorf("boltDB.getTodo() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if fmt.Sprintf("%T", got) != fmt.Sprintf("%T", tt.want) {
				t.Errorf("boltDB.getTodo() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_boltDB_saveTodo(t *testing.T) {
	type args struct {
		td *todo
	}
	tests := []struct {
		name    string
		db      *boltDB
		args    args
		wantErr bool
	}{
		{
			"saveTodo Unity test",
			app.datamapper.db().(*boltDB),
			args{&todo{Title: "New task title", Description: "New task description", Deadline: time.Now().AddDate(0, 0, 2), Priority: low}},
			false,
		},
		{
			"saveTodo Unity test",
			app.datamapper.db().(*boltDB),
			args{nil},
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.db.saveTodo(tt.args.td); (err != nil) != tt.wantErr {
				t.Errorf("boltDB.saveTodo() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
