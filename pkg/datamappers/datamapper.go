package datamapper

import "github.com/gplume/todo-list/pkg/models"

// DataMapper holds the active database system methods
type DataMapper interface {
	Db() interface{}
	Closing()
	models.ToDoMapper
}
