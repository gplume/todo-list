package main

import (
	"fmt"
	"log"
	"net/http"
	"sort"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

type todo struct {
	ID          int        `json:"id"`
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
		log.Println(err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	todo.Creation = time.Now()
	if ok, errors := todo.validator(); !ok {
		log.Println(errors)
		c.JSON(http.StatusBadRequest, errors)
		return
	}
	if err := app.datamapper.saveTodo(todo); err != nil {
		log.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, todo)
}

func updateTodo(c *gin.Context) {
	var todo *todo
	if err := c.BindJSON(&todo); err != nil {
		log.Println(err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if ok, errors := todo.validator(); !ok {
		log.Println(errors)
		c.JSON(http.StatusBadRequest, errors)
		return
	}
	if err := app.datamapper.updateTodo(todo); err != nil {
		log.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, todo)
}

func getTodo(c *gin.Context) {
	todoKey := c.Param("key")
	intKey, err := strconv.Atoi(todoKey)
	if err != nil {
		msg := fmt.Sprintf("error parsing todo key: %v", err)
		log.Println(msg)
		c.JSON(http.StatusBadRequest, gin.H{"error": msg})
	}
	todo, err := app.datamapper.getTodo(intKey)
	if err != nil {
		log.Println("app.datamapper.getTodo() error: ", err)
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, todo)
}

func deleteTodo(c *gin.Context) {
	todoKey := c.Param("key")
	intKey, err := strconv.Atoi(todoKey)
	if err != nil {
		msg := fmt.Sprintf("error parsing todo key: %v", err)
		log.Println(msg)
		c.JSON(http.StatusBadRequest, gin.H{"error": msg})
	}
	if err := app.datamapper.deleteTodo(intKey); err != nil {
		log.Println("app.datamapper.deleteTodo() error: ", err)
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusNoContent, nil)
}

func listTodos(c *gin.Context) {
	sorting := c.DefaultQuery("sort", "desc")
	todos, err := app.datamapper.listTodos()
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	switch sorting {
	case "asc":
		sort.Slice(todos, func(i, j int) bool {
			return todos[i].Deadline.Before(todos[j].Deadline)
		})
	case "priority":
		sort.Slice(todos, func(i, j int) bool {
			return todos[i].Priority < todos[j].Priority
		})
	default: // "desc"
		sort.Slice(todos, func(i, j int) bool {
			return todos[i].Deadline.After(todos[j].Deadline)
		})
	}
	c.JSON(http.StatusOK, todos)
}
