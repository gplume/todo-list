package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"strconv"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func getTestAppEngine() (*application, *gin.Engine) {
	var err error
	app, err = newApp(true)
	if err != nil {
		log.Fatalf("Error initializing Application : %v", err)
	}
	router := mainEngineAndRoutes()

	return app, router
}

func TestMainEngineAndRoutes(t *testing.T) {
	app, router := getTestAppEngine()
	assert := assert.New(t)
	assert.NotNil(app)
	assert.NotNil(router)

	// CREATE: /post-todo
	w1 := httptest.NewRecorder()
	newTodo := &todo{
		Deadline:    time.Now().AddDate(0, 0, 1), // +1 day
		Title:       "New task",
		Description: "Here's the description of the new task...",
		Priority:    urgent,
	}
	buf, err := json.Marshal(newTodo)
	assert.Nil(err)
	body := bytes.NewBuffer(buf)
	req, err := http.NewRequest("POST", "/todo", body)
	assert.Nil(err)
	router.ServeHTTP(w1, req)
	assert.Equal(http.StatusCreated, w1.Code)
	assert.NotNil(w1.Body)

	// RETREIVE: /list
	w2 := httptest.NewRecorder()
	req2, err := http.NewRequest("GET", "/todo", nil)
	assert.Nil(err)
	router.ServeHTTP(w2, req2)
	assert.Equal(http.StatusOK, w2.Code)
	assert.NotNil(w2.Body)

	getbody, err := ioutil.ReadAll(w2.Body)
	assert.Nil(err)
	e := make([]todo, 0)
	err = json.Unmarshal(getbody, &e)
	assert.Nil(err)
	firstTodo := e[0]
	newTodo.ID = firstTodo.ID
	assert.Equal(firstTodo.Deadline.Format(time.RFC3339), newTodo.Deadline.Format(time.RFC3339))
	assert.Equal(firstTodo.Title, newTodo.Title)
	assert.Equal(firstTodo.Description, newTodo.Description)
	assert.Equal(firstTodo.Priority, newTodo.Priority)

	// UPDATE: /update-todo
	newTodo.Title = "Updated task title"
	newTodo.Description = "Updated task description"
	newTodo.Priority = low
	buf3, err := json.Marshal(newTodo)
	assert.Nil(err)
	body3 := bytes.NewBuffer(buf3)

	w3 := httptest.NewRecorder()
	req3, err := http.NewRequest("PUT", "/todo", body3)
	assert.Nil(err)
	router.ServeHTTP(w3, req3)
	assert.Equal(http.StatusOK, w3.Code)

	updatedBody, err := ioutil.ReadAll(w3.Body)
	assert.Nil(err)
	var updated *todo
	err = json.Unmarshal(updatedBody, &updated)
	assert.Nil(err)
	assert.Equal(updated.Creation.Format(time.RFC3339), newTodo.Creation.Format(time.RFC3339))
	assert.Equal(updated.Deadline.Format(time.RFC3339), newTodo.Deadline.Format(time.RFC3339))
	assert.Equal(updated.Title, newTodo.Title)
	assert.Equal(updated.Description, newTodo.Description)
	assert.Equal(updated.Priority, newTodo.Priority)

	// DELETE: /delete-todo
	w4 := httptest.NewRecorder()
	todoKey := strconv.Itoa(updated.ID)
	url, err := url.Parse("/todo/" + todoKey)
	assert.Nil(err)
	req4, err := http.NewRequest("DELETE", url.EscapedPath(), nil)
	assert.Nil(err)
	router.ServeHTTP(w4, req4)
	assert.Equal(http.StatusNoContent, w4.Code)

	// Close DB and eventually delete the associated file
	app.datamapper.close()
	switch app.cfg.DBType {
	case "bolt":
		err = os.Remove(fmt.Sprintf("%s/%s", app.cfg.DBDirectory, app.cfg.DBTestName))
		assert.Nil(err)
	}

}
