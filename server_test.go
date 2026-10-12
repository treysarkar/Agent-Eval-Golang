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

// Test for create-todo

func TestCreate(t *testing.T){
	s:=NewServer()
	for i, title := range []string{"Learn Go", "Write a test"}{
		body, err:=json.Marshal(map[string]string{"title": title})
		if err !=nil{
			t.Fatal(err)
		}
		rec:=request(t,s,"POST","/todos", string(body))
		expectStatus(t,rec,http.StatusCreated)
		if rec.Header().Get("Content-Type") != "application/json" {
			t.Fatal("create response must be JSON")
		}

		todo:=readTodo(t,rec)
		if todo.ID !=i+1 || todo.Title !=title || todo.Done{
			t.Fatalf("Unexpected todo : %+v", todo)
		}

	}
}

// Title validation testing

func TestValidation(t *testing.T){
	s:= NewServer()
	 invalid := []string{
        "", "{", "null", "[]", "{}",
        `{"title":""}`, `{"title":"   "}`, `{"title":42}`,
        `{"title":"Go","extra":true}`,
        `{"title":"Go"} {"title":"Again"}`,
        `{"title":"Go"} trailing`,
        `{"title":"` + strings.Repeat("x", (1<<20)+1) + `"}`,
    }

	for i,body :=range invalid{
		rec:=request(t,s,"POST","/todos",body)
		if rec.Code != http.StatusBadRequest {
			t.Fatalf("invalid case %d: got %d, want 400", i, rec.Code)
		}
	}

	rec:=request(t,s,"POST","/todos",`{"title":"  Learn Go. "}`)
	expectStatus(t, rec, http.StatusCreated)
	todo:=readTodo(t,rec)
	if todo.ID != 1 || todo.Title != "Learn Go." || todo.Done {
		t.Fatalf("Unexpected todo : %+v", todo)
	}
}
