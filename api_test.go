package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strconv"
	"testing"
	"time"

	"github.com/gplume/todo-list/pkg/models"
	"github.com/stretchr/testify/assert"
)

func TestCRUDfromEndPoints(t *testing.T) {

	assert := assert.New(t)
	// require := require.New(t)
	newTodo := &models.Todo{
		Deadline:    time.Now().AddDate(0, 0, 1), // +1 day
		Title:       "New task",
		Description: "Here's the description of the new task...",
		Priority:    models.High,
	}
	// CREATE: POST NEW TODO
	t.Run("POST:/todo", func(t *testing.T) {
		rec1 := httptest.NewRecorder()
		buf, err := json.Marshal(newTodo)
		assert.Nil(err)
		body := bytes.NewBuffer(buf)
		req, err := http.NewRequest("POST", "/todo", body)
		assert.Nil(err)
		app.router.ServeHTTP(rec1, req)
		assert.Equal(http.StatusCreated, rec1.Code)
		assert.NotNil(rec1.Body)
	})

	// RETREIVE LIST
	t.Run("GET:/todo", func(t *testing.T) {
		rec2 := httptest.NewRecorder()
		req2, err := http.NewRequest("GET", "/todo", nil)
		assert.Nil(err)
		app.router.ServeHTTP(rec2, req2)
		assert.Equal(http.StatusOK, rec2.Code)
		assert.NotNil(rec2.Body)

		getbody, err := ioutil.ReadAll(rec2.Body)
		assert.Nil(err)
		e := make([]models.Todo, 0)
		err = json.Unmarshal(getbody, &e)
		assert.Nil(err)
		if assert.NotNil(e) && len(e) > 0 {
			newTodo.ID = e[0].ID
			assert.Equal("models.Todo", fmt.Sprintf("%T", e[0]))
		}
	})

	// RETREIVE SINGLE
	t.Run("GET:/todo/:id", func(t *testing.T) {
		rec := httptest.NewRecorder()
		req, err := http.NewRequest("GET", fmt.Sprintf("/todo/%d", newTodo.ID), nil)
		assert.Nil(err)
		app.router.ServeHTTP(rec, req)
		assert.Equal(http.StatusOK, rec.Code)
		assert.NotNil(rec.Body)

		getbody, err := ioutil.ReadAll(rec.Body)
		assert.Nil(err)
		var firstTodo *models.Todo
		err = json.Unmarshal(getbody, &firstTodo)
		assert.Nil(err)
		if assert.NotNil(firstTodo) {
			assert.Equal("*models.Todo", fmt.Sprintf("%T", firstTodo))
		}
	})

	// UPDATE
	t.Run("PUT:/todo", func(t *testing.T) {
		newTodo.Title = "Updated task title"
		newTodo.Description = "Updated task description"
		newTodo.Priority = models.Low
		buf3, err := json.Marshal(newTodo)
		assert.Nil(err)
		body3 := bytes.NewBuffer(buf3)

		rec3 := httptest.NewRecorder()
		req3, err := http.NewRequest("PUT", "/todo", body3)
		assert.Nil(err)
		app.router.ServeHTTP(rec3, req3)
		assert.Equal(http.StatusOK, rec3.Code)

		updatedBody, err := ioutil.ReadAll(rec3.Body)
		assert.Nil(err)
		var updated *models.Todo
		err = json.Unmarshal(updatedBody, &updated)
		assert.Nil(err)
		assert.Equal(newTodo.Creation.Format(time.RFC3339), updated.Creation.Format(time.RFC3339))
		assert.Equal(newTodo.Deadline.Format(time.RFC3339), updated.Deadline.Format(time.RFC3339))
		assert.Equal(newTodo.Title, updated.Title)
		assert.Equal(newTodo.Description, updated.Description)
		assert.Equal(newTodo.Priority, updated.Priority)
	})

	// DELETE
	t.Run("DELETE:/todo/:id", func(t *testing.T) {
		rec4 := httptest.NewRecorder()
		todoKey := strconv.Itoa(newTodo.ID)
		url, err := url.Parse("/todo/" + todoKey)
		assert.Nil(err)
		req4, err := http.NewRequest("DELETE", url.EscapedPath(), nil)
		assert.Nil(err)
		app.router.ServeHTTP(rec4, req4)
		assert.Equal(http.StatusNoContent, rec4.Code)
	})
}
