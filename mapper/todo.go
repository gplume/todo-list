package mapper

import (
	"time"
)

// Todo Type is exported because on datamappers....
type Todo struct {
	ID          int        `json:"id"`
	Creation    time.Time  `json:"creationDate"`
	Deadline    time.Time  `json:"deadlineDate"` // mandatory at insert
	Title       string     `json:"title"`        // mandatory at insert
	Description string     `json:"description"`
	Priority    Priorities `json:"priority"` // mandatory at insert
}

// Priorities are int value iota form high to low
type Priorities int

// High, Medium, Low...
const (
	_ Priorities = iota
	High
	Medium
	Low
)

// Validator verify todo structure at input
func (td *Todo) Validator() (bool, map[string]string) {

	errors := make(map[string]string)
	if td.Deadline.IsZero() {
		errors["deadline"] = "Please set a deadline date for the new task"
	}
	if td.Title == "" {
		errors["title"] = "Please insert a task Title"
	}
	if td.Priority == 0 {
		errors["priority"] = "Please set a priority for the task"
	}
	if td.Priority > 3 {
		errors["priority"] = "Please set a correct priority for the task"
	}
	return len(errors) == 0, errors
}

// DataMapper holds the active database system methods
type DataMapper interface {
	Db() interface{}
	Closing()
	SaveTodo(*Todo) error
	ListTodos(string) ([]*Todo, error)
	GetTodo(int) (*Todo, error)
	UpdateTodo(*Todo) error
	DeleteTodo(int) error
}
