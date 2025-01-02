package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"testing"
)

func TestE2ETodoService(t *testing.T) {
	// Set up the test server
	server := setupTestServer()

	// Define helper functions
	client := &http.Client{}

	// 1. Test Create a To-Do
	t.Run("CreateTodo", func(t *testing.T) {
		todoUrl := createAndFetchTodo(t, client, server)
		if todoUrl == "" {
			t.Fatalf("expected a todo URL, got %v", todoUrl)
		}
	})

	// 2. Test Get All To-Dos
	t.Run("GetTodos", func(t *testing.T) {
		resp, err := client.Get(server.URL + "/todos")
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		if resp.StatusCode != http.StatusOK {
			t.Errorf("expected status %v, got %v", http.StatusOK, resp.StatusCode)
		}

		var todos []map[string]interface{}
		err = json.NewDecoder(resp.Body).Decode(&todos)
		if err != nil {
			t.Errorf("Error while decoding get todos response, got %v", err)
		}
		if len(todos) < 1 {
			t.Errorf("expected at least one todo, got %v", len(todos))
		}
	})

	// 3. Test Update a To-Do
	t.Run("UpdateTodo", func(t *testing.T) {
		todoURL := createAndFetchTodo(t, client, server)

		todoUpdate := map[string]interface{}{
			"title":     "Updated Todo",
			"completed": true,
		}
		body, _ := json.Marshal(todoUpdate)
		req, _ := http.NewRequest(http.MethodPut, todoURL, bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		resp, err := client.Do(req)
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		if resp.StatusCode != http.StatusOK {
			t.Errorf("expected status %v, got %v", http.StatusOK, resp.StatusCode)
		}

		var updatedTodo map[string]interface{}
		err = json.NewDecoder(resp.Body).Decode(&updatedTodo)
		if err != nil {
			t.Errorf("Error while decoding update todo response, got %v", err)
		}
		if updatedTodo["title"] != "Updated Todo" {
			t.Errorf("expected title %v, got %v", "Updated Todo", updatedTodo["title"])
		}
		if !updatedTodo["completed"].(bool) {
			t.Errorf("expected completed to be true, got %v", updatedTodo["completed"].(bool))
		}
	})

	// 4. Test Delete a To-Do
	t.Run("DeleteTodo", func(t *testing.T) {
		todoURL := createAndFetchTodo(t, client, server)
		req, _ := http.NewRequest(http.MethodDelete, todoURL, nil)
		resp, err := client.Do(req)
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		if resp.StatusCode != http.StatusOK {
			t.Errorf("expected status %v, got %v", http.StatusOK, resp.StatusCode)
		}

		// Verify deletion
		resp, _ = client.Get(todoURL)
		if resp.StatusCode != http.StatusNotFound {
			t.Errorf("expected status %v, got %v", http.StatusNotFound, resp.StatusCode)
		}
	})
}

func createAndFetchTodo(t *testing.T, client *http.Client, server *httptest.Server) string {
	// Setup: Create a To-Do item to update
	todo := map[string]interface{}{
		"title":     "Initial Todo",
		"completed": false,
	}
	body, _ := json.Marshal(todo)
	resp, err := client.Post(server.URL+"/todos", "application/json", bytes.NewBuffer(body))
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if resp.StatusCode != http.StatusCreated {
		t.Fatalf("expected status %v, got %v", http.StatusCreated, resp.StatusCode)
	}

	err = json.NewDecoder(resp.Body).Decode(&todo)
	if err != nil {
		t.Errorf("Error while decoding create todo response, got %v", err)
	}
	if todo["title"] != "Initial Todo" {
		t.Errorf("expected title %v, got %v", "Test Todo", todo["title"])
	}
	if todo["completed"].(bool) {
		t.Errorf("expected completed to be false, got %v", todo["completed"].(bool))
	}
	id := strconv.Itoa(int(todo["id"].(float64)))
	return strings.Join([]string{server.URL, "/todos/", id}, "")
}

func setupTestServer() *httptest.Server {
	setupDatabase()
	engine := setupServer()
	return httptest.NewServer(engine)
}
