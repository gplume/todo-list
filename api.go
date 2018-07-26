package main

import (
	"fmt"
	"log"
	"net/http"
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
	high
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
	if td.Priority > 3 {
		errors["priority"] = "Please set a correct priority for the task"
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
	app.prom.postCount.Inc()
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
	app.prom.updateCount.Inc()
	c.JSON(http.StatusOK, todo)
}

func getTodo(c *gin.Context) {

	todoKey := c.Param("id")
	intKey, err := strconv.Atoi(todoKey)
	if err != nil {
		msg := fmt.Sprintf("error parsing todo key: %v", err)
		log.Println(msg)
		c.JSON(http.StatusBadRequest, gin.H{"error": msg})
		return
	}
	todo, err := app.datamapper.getTodo(intKey)
	if err != nil {
		log.Println("app.datamapper.getTodo() error: ", err)
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	app.prom.getCount.Inc()
	c.JSON(http.StatusOK, todo)
}

func deleteTodo(c *gin.Context) {

	todoKey := c.Param("id")
	intKey, err := strconv.Atoi(todoKey)
	if err != nil || intKey == 0 {
		msg := fmt.Sprintf("error parsing todo key: %v", err)
		log.Println(msg)
		c.JSON(http.StatusBadRequest, gin.H{"error": msg})
		return
	}
	// here we use getTodo because BoltDB returns no error if key doesn't exists!!
	if _, err := app.datamapper.getTodo(intKey); err != nil {
		log.Println("app.datamapper.deleteTodo() error: ", err)
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	if err := app.datamapper.deleteTodo(intKey); err != nil {
		log.Println("app.datamapper.deleteTodo() error: ", err)
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	app.prom.deleteCount.Inc()
	c.JSON(http.StatusNoContent, gin.MIMEJSON)
}

func listTodos(c *gin.Context) {

	sorting := c.DefaultQuery("sort", "asc")
	todos, err := app.datamapper.listTodos(sorting)
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	app.prom.listCount.Inc()
	c.JSON(http.StatusOK, todos)
}
