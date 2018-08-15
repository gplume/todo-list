package boltmapper

import (
	"encoding/binary"
	"encoding/json"
	"errors"
	"fmt"
	"sort"

	"github.com/boltdb/bolt"
	"github.com/gplume/todo-list/models"
)

type boltDB struct {
	*bolt.DB
	todos []byte
}

// func init() {
// 	fmt.Println("BoltMapper package started")
// }

// NewBoltDatamapper INIT FUNCTION
func NewBoltDatamapper(db *bolt.DB) (models.DataMapper, error) {
	var err error
	_, err = getBucket(db, "todos")
	if err != nil {
		return nil, err
	}

	return &boltDB{
		db,
		[]byte("todos"),
	}, nil
}

// used by "defer" in main goroutine
func (db *boltDB) Closing() {
	db.Close()
}

func (db *boltDB) Db() interface{} {
	return db
}

/***************************************
*************** UTILITIES **************
***************************************/

func getBucket(db *bolt.DB, bucketName string) (b *bolt.Bucket, err error) {
	err = db.Update(func(tx *bolt.Tx) error {
		var err error
		b, err = tx.CreateBucketIfNotExists([]byte(bucketName))
		return err
	})
	return
}

func itob(v int) []byte {
	b := make([]byte, 8)
	binary.BigEndian.PutUint64(b, uint64(v))
	return b
}

/***************************************
*************** METHODS ****************
***************************************/

// saveTodo persist bytes to todos bucket.
func (db *boltDB) SaveTodo(td *models.Todo) error {
	if td == nil {
		return errors.New("todo is nil")
	}
	return db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket(db.todos)
		if td.ID == 0 {
			id, _ := b.NextSequence()
			td.ID = int(id)
		}
		buf, err := json.Marshal(td)
		if err != nil {
			return fmt.Errorf("todo cannot be properly encoded: %v", err)
		}
		return b.Put(itob(td.ID), buf)
	})
}

// listTodos reteive all todos in db and sort them by sorting parameter
func (db *boltDB) ListTodos(sorting string) ([]*models.Todo, error) {
	todos := make([]*models.Todo, 0)
	err := db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket(db.todos)
		return b.ForEach(func(k, v []byte) error {
			if v != nil {
				var td *models.Todo
				if err := json.Unmarshal(v, &td); err != nil {
					return fmt.Errorf("error parsing todos: %v", err)
				}
				if v != nil {
					todos = append(todos, td)
				}
			}
			return nil
		})

	})
	switch sorting {
	case "desc":
		sort.Slice(todos, func(i, j int) bool {
			return todos[i].Deadline.After(todos[j].Deadline)
		})
	case "priority":
		sort.Slice(todos, func(i, j int) bool {
			return todos[i].Priority < todos[j].Priority
		})
	default: // "asc"
		sort.Slice(todos, func(i, j int) bool {
			return todos[i].Deadline.Before(todos[j].Deadline)
		})
	}
	return todos, err
}

// getTodo reteive simgle todo defined by todoKey (id)
func (db *boltDB) GetTodo(todoKey int) (*models.Todo, error) {
	var todo *models.Todo
	return todo, db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket(db.todos)
		td := b.Get(itob(todoKey))
		if td == nil {
			return fmt.Errorf("There's no todo with that ID! (%v)", todoKey)
		}
		return json.Unmarshal(td, &todo)
	})
}

// updateTodo is the same as saveTodo() but stays here for compatibility reason
// and eventual others db systems.
// boldDB key/value technology obviously don't need that
// an additionnal ID check is done for data integrity
func (db *boltDB) UpdateTodo(td *models.Todo) error {
	return db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket(db.todos)
		if td.ID == 0 {
			return errors.New("cannot update an ID==0 todo")
		}
		buf, err := json.Marshal(td)
		if err != nil {
			return fmt.Errorf("todo cannot be properly encoded: %v", err)
		}
		return b.Put(itob(td.ID), buf)
	})
}

// deleteTodo persistently delete record by key (id)
func (db *boltDB) DeleteTodo(todoKey int) error {
	return db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket(db.todos)
		if err := b.Delete(itob(todoKey)); err != nil {
			return fmt.Errorf("todo cannot be properly deleted: %v", err)
		}
		return nil
	})
}
