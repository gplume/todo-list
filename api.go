package main

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gplume/todo-list/models"
)

func addTodo(c *gin.Context) {

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
	if err := app.datamapper.SaveTodo(todo); err != nil {
		log.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	app.prom.PostCount.Inc()
	c.JSON(http.StatusCreated, todo)
}

func updateTodo(c *gin.Context) {

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
	if err := app.datamapper.UpdateTodo(todo); err != nil {
		log.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	app.prom.UpdateCount.Inc()
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
	todo, err := app.datamapper.GetTodo(intKey)
	if err != nil {
		log.Println("app.datamapper.getTodo() error: ", err)
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	app.prom.GetCount.Inc()
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
	if _, err := app.datamapper.GetTodo(intKey); err != nil {
		log.Println("app.datamapper.deleteTodo() error: ", err)
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	if err := app.datamapper.DeleteTodo(intKey); err != nil {
		log.Println("app.datamapper.deleteTodo() error: ", err)
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	app.prom.DeleteCount.Inc()
	c.JSON(http.StatusNoContent, gin.MIMEJSON)
}

func listTodos(c *gin.Context) {

	sorting := c.DefaultQuery("sort", "asc")
	todos, err := app.datamapper.ListTodos(sorting)
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	app.prom.ListCount.Inc()
	c.JSON(http.StatusOK, todos)
}
