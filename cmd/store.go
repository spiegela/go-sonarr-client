package main

import (
	"fmt"
	"os"

	"github.com/dgraph-io/badger"
)

var isVerbose bool

type store struct {
	db       *badger.DB
	isClosed bool
	keys     storeKeys
	secret   []byte
}

type storeKeys struct {
	appSecret []byte
	sonarrKey []byte
	sonarrURL []byte
}

func initDataStore(dirName string) (store, error) {
	var db store

	if isVerbose {
		fmt.Println("checking if our database exists in the home directory at:", dirName)
	}

	// create a directory for our database
	if _, err := os.Stat(dirName); os.IsNotExist(err) {
		if isVerbose {
			fmt.Println("creating directory because it doesn't exist")
		}

		if err := os.Mkdir(dirName, os.ModePerm); err != nil {
			return db, err
		}
	}

	options := badger.DefaultOptions

	options.Dir = dirName
	options.ValueDir = dirName

	kvStore, err := badger.Open(options)

	if err != nil {
		return db, err
	}

	if isVerbose {
		fmt.Println("successfully opened data store")
	}

	db.db = kvStore
	db.keys = storeKeys{
		sonarrKey: []byte("sonarr-key"),
		sonarrURL: []byte("sonarr-url"),
	}

	return db, nil
}

func (s store) Close() {
	if s.isClosed {
		fmt.Println("data store already closed")
		return
	}

	if err := s.db.Close(); err != nil {
		fmt.Printf("data store failed to closed: %v\n", err)
	}

	s.isClosed = true
}

func (s store) getSonarrKey() (string, error) {
	var sonarrKey string

	if err := s.db.View(func(txn *badger.Txn) error {
		item, err := txn.Get(s.keys.sonarrKey)

		if err != nil {
			return err
		}

		serializedSonarrURL, err := item.Value()

		if err != nil {
			return err
		}

		sonarrKey = string(serializedSonarrURL)

		return nil
	}); err != nil {
		return sonarrKey, err
	}

	if isVerbose {
		fmt.Printf("Your sonarr key is %s\n", sonarrKey)
	}

	return sonarrKey, nil
}

func (s store) saveSonarrKey(key string) error {
	if isVerbose {
		fmt.Printf("your sonarr key: %s\n", string(key))
	}

	if err := s.db.Update(func(txn *badger.Txn) error {
		return txn.Set(s.keys.sonarrKey, []byte(key), 0x00)
	}); err != nil {
		return err
	}

	if isVerbose {
		fmt.Println("saved key to store")
	}

	return nil
}

func (s store) getSonarrURL() (string, error) {
	var sonarrURL string

	err := s.db.View(func(txn *badger.Txn) error {
		item, err := txn.Get(s.keys.sonarrURL)

		if err != nil {
			return err
		}

		serializedServer, err := item.Value()

		if err != nil {
			return err
		}

		sonarrURL = string(serializedServer)

		return nil
	})

	return sonarrURL, err
}

func (s store) saveSonarrURL(sonarrURL string) error {
	return s.db.Update(func(txn *badger.Txn) error {
		return txn.Set(s.keys.sonarrURL, []byte(sonarrURL), 0)
	})
}
