package main

import (
	"github.com/sirupsen/logrus"
	"net/http"
	"os"
)

func main() {
	queueDBPath := os.Getenv("QUEUE_DB_PATH")
	log := logrus.New()

	queue, err := NewFileQueue(queueDBPath)
	if err != nil {
		log.Fatal(err)
	}

	httpMethods := httpMethods{
		queue: queue,
		log: log,
	}

	http.HandleFunc("/add", httpMethods.Add)
	http.HandleFunc("/pop", httpMethods.Pop)

	log.Fatal(http.ListenAndServe(":8080", nil))
}
