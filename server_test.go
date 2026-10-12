// This is Test helper

package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"strings"
)

// A seperate wire type will check the public JSOn contract independently

type wireTodo struct {
	ID    int    `json:"id"`
	Title string `json:"title"`
	Done  bool   `json:"done"`
}

func request(t *testing.T, h http.Handler, method, path, body string) *httptest.ResponseRecorder {
	t.Helper()
	req := httptest.NewRequest(method, path, strings.NewReader(body))
	req.Header.Set("Contetn-Type", "application/json")
	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, req)
	return rec
}

func expectStatus(t *testing.T, rec *httptest.ResponseRecorder, want int) {
	t.Helper()
	if rec.Code != want {
		t.Fatalf("status: got %d, want %d; body=%s", rec.Code, want, rec.Body.String())
	}
}

func readTodo(t *testing.T, rec *httptest.ResponseRecorder) wireTodo {
	t.Helper()
	var todo wireTodo
	if err := json.NewDecoder(rec.Body).Decode(&todo); err != nil {
		t.Fatalf("failed to read todo: %v", err)
	}
	return todo
}

func TestHealth(t *testing.T) {
	rec := request(t, NewServer(), "GET", "/health", "")
	expectStatus(t, rec, http.StatusOK)
	if rec.Header().Get("Content-Type") != "application/json" {
		t.Fatal("health response must be JSON")
	}
	var body map[string]string
	if err := json.Unmarshal(rec.Body.Bytes(), &body); err != nil {
		t.Fatal(err)
	}
	if body["status"] != "ok" {
		t.Fatalf("health body: got %v", body)
	}
}
