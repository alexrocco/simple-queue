package main

import (
	"encoding/json"
	"fmt"
	"github.com/sirupsen/logrus"
	"io/ioutil"
	"net/http"
)

type httpMethods struct {
	queue Queue
	log *logrus.Logger
}

func (h * httpMethods) Add(w http.ResponseWriter, r *http.Request) {
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

	h.log.Infof("Body request: %s", string(body))

	var data interface{}
	err = json.Unmarshal(body, &data)
	if err != nil {
		w.WriteHeader(http.StatusUnprocessableEntity)
		_, _ = w.Write([]byte(fmt.Sprintf("error parsing the request body to JSON: %s", err.Error())))

		return
	}

	err = h.queue.Add(data)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write([]byte(fmt.Sprintf("error adding the message to the queue: %s", err.Error())))

		return
	}

	w.WriteHeader(http.StatusOK)
}

func (h * httpMethods) Pop(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		_, _ = w.Write([]byte("it should be a POST request"))

		return
	}

	value, err := h.queue.Pop()
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
}