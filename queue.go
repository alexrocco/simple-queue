package main

import (
	"encoding/json"
	"github.com/pkg/errors"
	"os"
	"sync"
	"time"
)

type item struct {
	Value     interface{} `json:"value"`
	CreatedAt string      `json:"created_at"`
}

type Queue struct {
	items  []item
	lock   sync.Mutex
	dbPath string
}

//NewQueue creates a new Queue by parsing an existent db file
func NewQueue(dbPath string) (*Queue, error) {
	dbContent, err := os.ReadFile(dbPath)

	items := make([]item, 0, 10)

	switch {
	// File does not exit, create it
	case os.IsNotExist(err):
		dbFile, err := os.Create(dbPath)
		if err != nil {
			return nil, errors.Wrap(err, "error creating db file")
		}

		err = dbFile.Close()
		if err != nil {
			return nil, errors.Wrap(err, "error closing db file")
		}
	// File exist but got some error, permission maybes
	case err != nil && os.IsExist(err):
		return nil, errors.Wrap(err, "error reading db file")
	// File exist, so parse it
	default:
		err = json.Unmarshal(dbContent, &items)
		if err != nil {
			return nil, errors.Wrap(err, "error parsing db file")
		}
	}

	queue := &Queue{
		items:  items,
		lock:   sync.Mutex{},
		dbPath: dbPath,
	}

	return queue, nil
}

// Add adds a new value to the bottom of the slice, to respect the FIFO data structure
func (q *Queue) Add(value interface{}) error {
	q.lock.Lock()

	q.items = append(q.items, item{
		Value:     value,
		CreatedAt: time.Now().Format(time.RFC3339),
	})

	err := q.updateDB()
	if err != nil {
		return err
	}

	q.lock.Unlock()

	return nil
}

// Pop pops the first element in the slice and remove it, respecting the FIFO data structure
func (q *Queue) Pop() (interface{}, error) {
	q.lock.Lock()

	if len(q.items) == 0 {
		q.lock.Unlock()

		return nil, nil
	}

	// Gets the first item (FIFO)
	item := q.items[0]

	// Remove it from the list
	q.items = q.items[1:]

	err := q.updateDB()
	if err != nil {
		return nil, err
	}

	q.lock.Unlock()

	return item.Value, nil
}

// updateDB updates the DB file (JSON) after the slice has been updated
func (q *Queue) updateDB() error {
	dbContent, err := json.Marshal(q.items)
	if err != nil {
		return err
	}

	err = os.WriteFile(q.dbPath, dbContent, os.ModePerm)
	if err != nil {
		return err
	}

	return nil
}
