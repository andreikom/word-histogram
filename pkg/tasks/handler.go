package tasks

import "word-histrogram/pkg/storage"

type Handler struct {
	storage.WordStore
}

// Task shall be executed by a worker pool
type Task interface {
	Execute()
}
