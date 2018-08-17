package api

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gplume/todo-list/pkg/engine"
	"github.com/gplume/todo-list/pkg/models"
	prome "github.com/gplume/todo-list/pkg/prometheus"
)

// AddTodo insert a complet Todo structure in DB
func AddTodo(c *gin.Context) {

	var todo *models.Todo
	if err := c.BindJSON(&todo); err != nil {
		log.Println(err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	todo.Creation = time.Now()
	if ok, errors := todo.Validator(); !ok {
		log.Println(errors)
		c.JSON(http.StatusBadRequest, errors)
		return
	}
	if err := engine.App.Datamapper.SaveTodo(todo); err != nil {
		log.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	prome.Vars.PostCount.Inc()
	c.JSON(http.StatusCreated, todo)
}

// UpdateTodo update a task by full structure
func UpdateTodo(c *gin.Context) {

	var todo *models.Todo
	if err := c.BindJSON(&todo); err != nil {
		log.Println(err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if ok, errors := todo.Validator(); !ok {
		log.Println(errors)
		c.JSON(http.StatusBadRequest, errors)
		return
	}
	if err := engine.App.Datamapper.UpdateTodo(todo); err != nil {
		log.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	prome.Vars.UpdateCount.Inc()
	c.JSON(http.StatusOK, todo)
}

// GetTodo retreive a task by ID
func GetTodo(c *gin.Context) {

	todoKey := c.Param("id")
	intKey, err := strconv.Atoi(todoKey)
	if err != nil {
		msg := fmt.Sprintf("error parsing todo key: %v", err)
		log.Println(msg)
		c.JSON(http.StatusBadRequest, gin.H{"error": msg})
		return
	}
	todo, err := engine.App.Datamapper.GetTodo(intKey)
	if err != nil {
		log.Println("engine.App.Datamapper.getTodo() error: ", err)
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	prome.Vars.GetCount.Inc()
	c.JSON(http.StatusOK, todo)
}

// DeleteTodo delete a task by ID
func DeleteTodo(c *gin.Context) {

	todoKey := c.Param("id")
	intKey, err := strconv.Atoi(todoKey)
	if err != nil || intKey == 0 {
		msg := fmt.Sprintf("error parsing todo key: %v", err)
		log.Println(msg)
		c.JSON(http.StatusBadRequest, gin.H{"error": msg})
		return
	}
	// here we use getTodo because BoltDB returns no error if key doesn't exists!!
	if _, err := engine.App.Datamapper.GetTodo(intKey); err != nil {
		log.Println("engine.App.Datamapper.deleteTodo() error: ", err)
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	if err := engine.App.Datamapper.DeleteTodo(intKey); err != nil {
		log.Println("engine.App.Datamapper.deleteTodo() error: ", err)
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	prome.Vars.DeleteCount.Inc()
	c.JSON(http.StatusNoContent, gin.MIMEJSON)
}

// ListTodos retreive a list of tasks
func ListTodos(c *gin.Context) {

	sorting := c.DefaultQuery("sort", "asc")
	todos, err := engine.App.Datamapper.ListTodos(sorting)
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	prome.Vars.ListCount.Inc()
	c.JSON(http.StatusOK, todos)
}
