package main

import (
	"bytes"
	"fmt"
	"github.com/alexrocco/simple-queue/mock"
	"github.com/golang/mock/gomock"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
)

func Test_httpMethods_Add(t *testing.T) {
	t.Run("It should return 200 and add a message to the queue", func(t *testing.T) {
		msg := "test"
		ctrl := gomock.NewController(t)
		mockQueue := mock.NewMockQueue(ctrl)
		payload := fmt.Sprintf("%q", msg)

		mockQueue.EXPECT().Add(msg).Times(1)

		log := logrus.New()
		log.Out = ioutil.Discard

		httpMethods := httpMethods{
			queue: mockQueue,
			log:   log,
		}

		buf := bytes.NewBufferString(payload)

		req := httptest.NewRequest(http.MethodPost, "http://localhost:8080", buf)
		w := httptest.NewRecorder()

		httpMethods.Add(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
	})
	t.Run("It should return 405 when method is different then post", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		mockQueue := mock.NewMockQueue(ctrl)

		log := logrus.New()
		log.Out = ioutil.Discard

		httpMethods := httpMethods{
			queue: mockQueue,
			log:   log,
		}

		req := httptest.NewRequest(http.MethodGet, "http://localhost:8080", nil)
		w := httptest.NewRecorder()

		httpMethods.Add(w, req)

		assert.Equal(t, http.StatusMethodNotAllowed, w.Code)
	})
	t.Run("It should return 422 when no payload is provided", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		mockQueue := mock.NewMockQueue(ctrl)

		log := logrus.New()
		log.Out = ioutil.Discard

		httpMethods := httpMethods{
			queue: mockQueue,
			log:   log,
		}

		req := httptest.NewRequest(http.MethodPost, "http://localhost:8080", nil)
		w := httptest.NewRecorder()

		httpMethods.Add(w, req)

		assert.Equal(t, http.StatusUnprocessableEntity, w.Code)
	})
	t.Run("It should return 500 when add method returns an error", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		mockQueue := mock.NewMockQueue(ctrl)
		mockQueue.EXPECT().Add(gomock.Any()).Return(errors.New("any error"))

		log := logrus.New()
		log.Out = ioutil.Discard

		httpMethods := httpMethods{
			queue: mockQueue,
			log:   log,
		}
		body := bytes.NewBufferString(`"test"`)

		req := httptest.NewRequest(http.MethodPost, "http://localhost:8080", body)
		w := httptest.NewRecorder()

		httpMethods.Add(w, req)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
	})
}

func Test_httpMethods_Pop(t *testing.T) {
	t.Run("It should return 200 and pop a message from the queue", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		mockQueue := mock.NewMockQueue(ctrl)
		mockQueue.EXPECT().Pop().Times(1)

		log := logrus.New()
		log.Out = ioutil.Discard

		httpMethods := httpMethods{
			queue: mockQueue,
			log:   log,
		}

		req := httptest.NewRequest(http.MethodGet, "http://localhost:8080", nil)
		w := httptest.NewRecorder()

		httpMethods.Pop(w, req)
	})
	t.Run("It should return 405 when method is different then get", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		mockQueue := mock.NewMockQueue(ctrl)

		log := logrus.New()
		log.Out = ioutil.Discard

		httpMethods := httpMethods{
			queue: mockQueue,
			log:   log,
		}

		req := httptest.NewRequest(http.MethodPost, "http://localhost:8080", nil)
		w := httptest.NewRecorder()

		httpMethods.Pop(w, req)

		assert.Equal(t, http.StatusMethodNotAllowed, w.Code)
	})
	t.Run("It should return 500 when pop method returns an error", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		mockQueue := mock.NewMockQueue(ctrl)
		mockQueue.EXPECT().Pop().Return(nil, errors.New("any error"))

		log := logrus.New()
		log.Out = ioutil.Discard

		httpMethods := httpMethods{
			queue: mockQueue,
			log:   log,
		}

		req := httptest.NewRequest(http.MethodGet, "http://localhost:8080", nil)
		w := httptest.NewRecorder()

		httpMethods.Pop(w, req)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
	})
}
