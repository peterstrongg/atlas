package routes

import (
	"errors"
	"fmt"
	"strings"

	"github.com/dgraph-io/badger/v4"
	"github.com/gin-gonic/gin"
)

type KV struct {
	db *badger.DB
}

type Host struct {
	IpAddress  string
	MacAddress string
	Dynamic    bool
}

func ReportScan(c *gin.Context) {
	var hosts []Host
	err := c.BindJSON(&hosts)
	if err != nil {
		// TODO: Handle errors
	}

	db, err := NewBadgerDb("./mac_vendors")
	if err != nil {
		panic(err)
	}
	defer db.Close()

	for _, host := range hosts {
		mac := strings.ReplaceAll(host.MacAddress, "-", ":")
		dbKey := strings.ToUpper(strings.ReplaceAll(mac, ":", "")[:6])
		keyExists, _ := db.Exists(dbKey)
		if keyExists {
			value, _ := db.Get(dbKey)
			fmt.Println(value)
		}
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

func (k *KV) Close() error {
	return k.db.Close()
}
