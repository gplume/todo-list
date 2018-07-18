package main

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

type todo struct {
	Creation    time.Time  `json:"creationDate"`
	Deadline    time.Time  `json:"deadlineDate"` // mandatory at insert
	Title       string     `json:"title"`        // mandatory at insert
	Description string     `json:"description"`
	Priority    priorities `json:"priority"` // mandatory at insert
}

type priorities int

const (
	_ priorities = iota
	urgent
	medium
	low
)

func (td *todo) validator() (bool, map[string]string) {
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
	return len(errors) == 0, errors
}

func addTodo(c *gin.Context) {
	var todo *todo
	if err := c.BindJSON(&todo); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if ok, errors := todo.validator(); !ok {
		c.JSON(http.StatusBadRequest, errors)
		return
	}
	if err := app.datamapper.createTodo(todo); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, todo)
}

func updateTodo(c *gin.Context) {
	var todo *todo
	if err := c.BindJSON(&todo); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if ok, errors := todo.validator(); !ok {
		c.JSON(http.StatusBadRequest, errors)
		return
	}
	if err := app.datamapper.updateTodo(todo); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, todo)
}

func deleteTodo(c *gin.Context) {
	var todo *todo
	if err := c.BindJSON(&todo); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := app.datamapper.deleteTodo(todo); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, todo)
}

func listTodos(c *gin.Context) {
	todos, err := app.datamapper.listTodos()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, todos)
}
