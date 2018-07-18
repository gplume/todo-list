package main

import "time"

type todo struct {
	Creation    time.Time `json:"creationDate"`
	Deadline    time.Time `json:"deadlineDate"` // mandatory at insert
	Title       string    `json:"title"`        // mandatory at insert
	Description string    `json:"description"`
	Status      int       `json:"status"` // mandatory at insert
}

type dataMapper interface {
	close()
	todosMapper
}

type todosMapper interface {
	createTodo(*todo) error
	listAllTodos() ([]*todo, error)
	updateTodo(*todo) error
	deleteTodo(*todo) error
}
