package apperrors

import "errors"

var ErrInvalidCommand = errors.New("Invalid command")
var ErrKeyNotFound = errors.New("key not found")
var ErrKeyExists = errors.New("key already exists")
var ErrQueueEmpty = errors.New("queue is empty")
var ErrQueueTimeout = errors.New("queue timeout")
