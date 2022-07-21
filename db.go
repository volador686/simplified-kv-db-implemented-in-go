package db

import (
	"errors"
	"fmt"
	"log"

	"github.com/boltdb/bolt"
	"github.com/idoubi/goz"
)

var defaultBucket = []byte("default")
var replicaBucket = []byte("replica")

// Database -boltdb(open source)
type Database struct {
	db *bolt.DB
}

// returns an instance of a database that we can work with
func NewDatabase(dbPath string) (db *Database, closeFunc func() error, err error) {
	boltDb, err := bolt.Open(dbPath, 0600, nil)
	if err != nil {
		return nil, nil, err
	}

	closeFunc = boltDb.Close
	db = &Database{
		db: boltDb,
	}

	// the main bucket contains all the data
	if err := db.createDefaultBucket(); err != nil {
		closeFunc()
		return nil, nil, fmt.Errorf("creating default bucket: %w", err)
	}

	// add a replica bucket to store the key-value pairs
	if err := db.createReplicaBucket(); err != nil {
		closeFunc()
		return nil, nil, fmt.Errorf("creating replica bucket: %w", err)
	}

	return db, closeFunc, nil
}

func (d *Database) createDefaultBucket() error {
	return d.db.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists(defaultBucket)
		return err
	})
}

func (d *Database) createReplicaBucket() error {
	return d.db.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists(replicaBucket)
		return err
	})
}

// SetKey sets the key to the requested value into the default database or returns an error
func (d *Database) SetKey(key string, value []byte) error {
	// original version: add key-value pair to the default bucket
	// return d.db.Update(func(tx *bolt.Tx) error {
	// 	b := tx.Bucket(defaultBucket)
	// 	return b.Put([]byte(key), value)
	// })

	// new version: add key-value pair to the default/replica bucket
	return d.db.Update(func(tx *bolt.Tx) error {
		if err := (tx.Bucket(defaultBucket)).Put([]byte(key), value); err != nil {
			return errors.New("default bucket not available")
		}
		return tx.Bucket(replicaBucket).Put([]byte(key), value)
	})
}

// GetKey gets the value of the requested key from a default database
// as for get, there is no need to change
// because our replication is not for get/set use, it is only uesd to prevent database failure
func (d *Database) GetKey(key string) ([]byte, error) {
	var result []byte
	err := d.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket(defaultBucket)
		result = b.Get([]byte(key))
		return nil
	})

	if err == nil {
		return result, nil
	}

	return nil, err
}

func (d *Database) SendReplica(replica_addr string) (string, string) {
	var res1 []byte
	var res2 []byte
	d.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket(replicaBucket)
		c := b.Cursor()
		res1, res2 = c.First()
		for key, value := c.First(); key != nil; key, value = c.Next() {
			log.Printf("key = %v, value = %v\n", string(key), string(value))
			cli := goz.NewClient()
			url := "http://" + replica_addr + "/set"
			resp, err := cli.Post(string(url), goz.Options{
				FormParams: map[string]interface{}{
					"key":   string(key),
					"value": string(value),
				},
			})
			if err != nil {
				log.Printf("replication error:%v\n", err)
			}
			body, _ := resp.GetBody()
			log.Println(body)
		}
		// for key, _ := c.First(); key != nil; key, _ = c.Next() {
		// 	b := tx.Bucket(replicaBucket)
		// 	b.Delete(key)
		// }
		return nil
	})
	return string(res1), string(res2)
}
