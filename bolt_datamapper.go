package main

import (
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

func (db *boltDB) createTodo(*todo) error {

	return nil
}

func (db *boltDB) listAllTodos() ([]*todo, error) {

	return nil, nil
}

func (db *boltDB) updateTodo(td *todo) error {

	return nil
}

func (db *boltDB) deleteTodo(td *todo) error {

	return nil
}
