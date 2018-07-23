package main

import (
	"encoding/binary"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/boltdb/bolt"
)

type boltDB struct {
	*bolt.DB
	todos []byte
}

/*********** INIT FUNCTION ************/
func newBoltDatamapper(db *bolt.DB, cfg *config) (dataMapper, error) {
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
func (db *boltDB) close() {
	db.Close()
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

func (db *boltDB) saveTodo(td *todo) error {
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
		// Persist bytes to todos bucket.

		return b.Put(itob(td.ID), buf)
	})
}

func (db *boltDB) listTodos() ([]*todo, error) {
	todos := make([]*todo, 0)
	return todos, db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket(db.todos)
		return b.ForEach(func(k, v []byte) error {
			if v != nil {
				var td *todo
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
}

func (db *boltDB) getTodo(todoKey int) (*todo, error) {
	var todo *todo
	return todo, db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket(db.todos)
		td := b.Get(itob(todoKey))
		if td == nil {
			return fmt.Errorf("There's no todo with that ID! (%v)", todoKey)
		}
		return json.Unmarshal(td, &todo)
	})
}

// same as saveTodo() but stays here for compatibility reason and others db systems.
// boldDB key/value technology obviously don't need that
// an additionnal ID check is done for data integrity
func (db *boltDB) updateTodo(td *todo) error {
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

func (db *boltDB) deleteTodo(todoKey int) error {
	return db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket(db.todos)
		if err := b.Delete(itob(todoKey)); err != nil {
			return fmt.Errorf("todo cannot be properly deleted: %v", err)
		}
		return nil
	})
}
