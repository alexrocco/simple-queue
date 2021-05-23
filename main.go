package main

import (
	"encoding/json"
	"fmt"
	"github.com/sirupsen/logrus"
	"io/ioutil"
	"net/http"
	"os"
)

func main() {
	queueDBPath := os.Getenv("QUEUE_DB_PATH")
	log := logrus.New()

	queue, err := NewQueue(queueDBPath)
	if err != nil {
		log.Fatal(err)
	}

	http.HandleFunc("/add", func(w http.ResponseWriter, r *http.Request) {
		defer r.Body.Close()

		if r.Method != http.MethodPost {
			w.WriteHeader(http.StatusMethodNotAllowed)
			_, _ = w.Write([]byte("it should be a POST request"))

			return
		}

		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			_, _ = w.Write([]byte(fmt.Sprintf("error reading the request body: %s", err.Error())))

			return
		}

		log.Infof("Body request: %s", string(body))

		var data interface{}
		err = json.Unmarshal(body, &data)
		if err != nil {
			w.WriteHeader(http.StatusUnprocessableEntity)
			_, _ = w.Write([]byte(fmt.Sprintf("error parsing the request body to JSON: %s", err.Error())))

			return
		}

		err = queue.Add(data)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			_, _ = w.Write([]byte(fmt.Sprintf("error adding the message to the queue: %s", err.Error())))

			return
		}

		w.WriteHeader(http.StatusOK)
	})

	http.HandleFunc("/pop", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			w.WriteHeader(http.StatusMethodNotAllowed)
			_, _ = w.Write([]byte("it should be a POST request"))

			return
		}

		r.

		value, err := queue.Pop()
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			_, _ = w.Write([]byte(fmt.Sprintf("error poping the message from the queue: %s", err.Error())))

			return
		}

		bytes, err := json.Marshal(value)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			_, _ = w.Write([]byte(fmt.Sprintf("error parsing the message from the queue: %s", err.Error())))

			return
		}

		w.WriteHeader(http.StatusOK)
		_, _ = w.Write(bytes)
	})

	log.Fatal(http.ListenAndServe(":8080", nil))
}
