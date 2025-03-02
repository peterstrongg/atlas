package main

import (
	"bufio"
	"errors"
	"fmt"
	"os"

	"github.com/dgraph-io/badger/v4"
)

type KV struct {
	db *badger.DB
}

func main() {
	db, _ := NewBadgerDb("./mac_vendors")
	defer db.Close()

	file, _ := os.Open("mac-vendor.txt")
	defer file.Close()

	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanLines)

	for scanner.Scan() {
		db.Set(scanner.Text()[:6], scanner.Text()[7:])
	}
}

func NewBadgerDb(pathToDb string) (*KV, error) {
	opts := badger.DefaultOptions(pathToDb)

	opts.Logger = nil
	badgerInstance, err := badger.Open(opts)
	if err != nil {
		return nil, fmt.Errorf("opening kv: %w", err)
	}

	return &KV{db: badgerInstance}, nil
}

func (k *KV) Close() error {
	return k.db.Close()
}

// nolint:wrapcheck
func (k *KV) Exists(key string) (bool, error) {
	var exists bool
	err := k.db.View(
		func(tx *badger.Txn) error {
			if val, err := tx.Get([]byte(key)); err != nil {
				return err
			} else if val != nil {
				exists = true
			}
			return nil
		})
	if errors.Is(err, badger.ErrKeyNotFound) {
		err = nil
	}
	return exists, err
}

func (k *KV) Get(key string) (string, error) {
	var value string
	return value, k.db.View(
		func(tx *badger.Txn) error {
			item, err := tx.Get([]byte(key))
			if err != nil {
				return fmt.Errorf("getting value: %w", err)
			}
			valCopy, err := item.ValueCopy(nil)
			if err != nil {
				return fmt.Errorf("copying value: %w", err)
			}
			value = string(valCopy)
			return nil
		})
}

func (k *KV) Set(key, value string) error {
	return k.db.Update(
		func(txn *badger.Txn) error {
			return txn.Set([]byte(key), []byte(value))
		})
}

func (k *KV) Delete(key string) error {
	return k.db.Update(
		func(txn *badger.Txn) error {
			return txn.Delete([]byte(key))
		})
}
