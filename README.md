# Simple Queue

[![codecov](https://codecov.io/gh/alexrocco/simple-queue/branch/master/graph/badge.svg?token=7GS1BOLHQV)](https://codecov.io/gh/alexrocco/simple-queue)

A simple queue app that respects FIFO and saves the current queue in a file (JSON), and exposes the operations by an HTTP server.
It also ensures that the operations, add and pop, are thread-safe by creating a shared lock between them.

It was created to be used in small and home projects and should not be used in production.

## HTTP methods

An HTTP server exposes the following methods on port 8080

### POST /add

Adds the payload (body) from the request to the end of the queue.
The only contract is that the payload should be able to be parsed in JSON.

### GET /pop

Gets the first element from the queue, removes it and returns it as a JSON.
