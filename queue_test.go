package main

import (
	"encoding/json"
	"fmt"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"math/rand"
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestNewQueue(t *testing.T) {
	t.Run("It should load the db file and parse the items", func(t *testing.T) {
		dbContent := `
		[
			{
				"value": "test123",
				"created_at": "2006-01-02T15:04:05Z07:00"
			}
		]
		`
		// Creates the tmp file for the db file
		dbFile, err := os.CreateTemp(os.TempDir(), "")
		assert.NoError(t, err)

		// Writes the db content
		_, err = dbFile.WriteString(dbContent)
		assert.NoError(t, err)

		// Close it to not mess with the test
		err = dbFile.Close()
		assert.NoError(t, err)

		queue, err := NewQueue(dbFile.Name())
		assert.NoError(t, err)

		expectedItems := []item{
			{
				Value:     "test123",
				CreatedAt: "2006-01-02T15:04:05Z07:00",
			},
		}

		assert.Equal(t, expectedItems, queue.items)

		// Remove temp file
		_ = os.Remove(dbFile.Name())
	})
	t.Run("It should create a db file when it does not exist", func(t *testing.T) {
		testDBFile := filepath.Join(os.TempDir(), "db-file.json")
		_, err := NewQueue(testDBFile)
		assert.NoError(t, err)

		stat, err := os.Stat(testDBFile)
		assert.NoError(t, err)

		assert.True(t, !stat.IsDir())

		_ = os.Remove(testDBFile)
	})
}

func TestQueue_Add(t *testing.T) {
	t.Run("It should add an item in an empty queue", func(t *testing.T) {
		testDBFile := filepath.Join(os.TempDir(), "db-file.json")
		queue, err := NewQueue(testDBFile)
		assert.NoError(t, err)

		expectedValue := "test"
		err = queue.Add(expectedValue)
		assert.NoError(t, err)

		dbContent, err := ioutil.ReadFile(testDBFile)
		assert.NoError(t, err)

		var items []item
		err = json.Unmarshal(dbContent, &items)
		assert.NoError(t, err)

		assert.Equal(t, len(items), 1)
		assert.Equal(t, items[0].Value, expectedValue)

		_ = os.Remove(testDBFile)
	})
	t.Run("It should add an item in the queue respecting FIFO", func(t *testing.T) {
		var items []item
		items = append(items, item{
			Value: "previous-value",
		})

		dbContent, err := json.Marshal(items)
		assert.NoError(t, err)

		dbPath := filepath.Join(os.TempDir(), "db-file.json")
		dbFile, err := os.Create(dbPath)
		assert.NoError(t, err)

		_, err = dbFile.Write(dbContent)
		assert.NoError(t, err)

		queue, err := NewQueue(dbPath)
		assert.NoError(t, err)

		expectedValue := "new-value"
		err = queue.Add(expectedValue)
		assert.NoError(t, err)

		dbUpdatedContent, err := ioutil.ReadFile(dbPath)
		assert.NoError(t, err)

		var updatedItems []item
		err = json.Unmarshal(dbUpdatedContent, &updatedItems)
		assert.NoError(t, err)

		assert.Equal(t, len(updatedItems), 2)
		assert.Equal(t, updatedItems[1].Value, expectedValue)

		_ = os.Remove(dbPath)
	})
}

func TestQueue_Pop(t *testing.T) {
	t.Run("It should pop the first item in the queue", func(t *testing.T) {
		var items []item
		expectedValue := "test"
		expectedItem := item{
			Value: expectedValue,
		}
		items = append(items, expectedItem)
		items = append(items, item{
			Value: "one-more",
		})
		items = append(items, item{
			Value: "two-more",
		})

		dbContent, err := json.Marshal(items)
		assert.NoError(t, err)

		dbPath := filepath.Join(os.TempDir(), "db-file.json")
		dbFile, err := os.Create(dbPath)
		assert.NoError(t, err)

		_, err = dbFile.Write(dbContent)
		assert.NoError(t, err)

		err = dbFile.Close()
		assert.NoError(t, err)

		queue, err := NewQueue(dbPath)
		assert.NoError(t, err)

		gotValue, err := queue.Pop()
		assert.NoError(t, err)

		assert.Equal(t, expectedValue, gotValue)
		assert.Equal(t, len(queue.items), 2)
	})
	t.Run("It should return no value when the queue is empty", func(t *testing.T) {
		rand.Seed(time.Now().UnixNano())
		dbPath := filepath.Join(os.TempDir(), fmt.Sprintf("%d", rand.Int()))

		queue, err := NewQueue(dbPath)
		assert.NoError(t, err)

		gotValue, err := queue.Pop()

		assert.Equal(t, nil, gotValue)

		_ = os.Remove(dbPath)
	})
}
