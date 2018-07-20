package main

import (
	"encoding/json"
	"fmt"
	"time"

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

/***************************************
*************** METHODS ****************
***************************************/

func (db *boltDB) createTodo(td *todo) error {
	return db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket(db.todos)
		buf, err := json.Marshal(td)
		if err != nil {
			return fmt.Errorf("todo cannot be properly encoded: %v", err)
		}
		// Persist bytes to todos bucket.
		// The deadline date is used as key so that is automatically sorted
		// thanks to bolt database btree pagination system
		// One major drawback of that is if the frontend API allow to change this date
		// we should cache it in the updated object so we can find the correct item to replace
		// (delete old one then create new one, limitation of the choosen ordered date key for value in DB)
		// as we cannot modify a key...
		return b.Put([]byte(td.Deadline.Format(time.RFC3339)), buf)
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

func (db *boltDB) updateTodo(td *todo) error {
	return db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket(db.todos)
		buf, err := json.Marshal(td)
		if err != nil {
			return fmt.Errorf("todo cannot be properly encoded: %v", err)
		}
		return b.Put([]byte(td.Deadline.Format(time.RFC3339)), buf)
	})
}

func (db *boltDB) deleteTodo(todoKey string) error {
	return db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket(db.todos)
		if err := b.Delete([]byte(todoKey)); err != nil {
			return fmt.Errorf("todo cannot be properly deleted: %v", err)
		}
		return nil
	})
}
