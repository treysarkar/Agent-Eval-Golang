package main

import (
	"encoding/json"
	"net/http"
	"sync"
	"io"
	"strings"
)

type Todo struct {
	ID    int    `json:"id`
	Title string `json:"title"`
	Done  bool   `json:"done"`
}

type Server struct {
	mu     sync.Mutex
	todos  []Todo
	nextID int
	mux    *http.ServeMux
}

func NewServer() *Server {
	s := &Server{todos: []Todo{}, nextID: 1, mux: http.NewServeMux()}
	s.mux.HandleFunc("GET /health", s.health)
	s.mux.HandleFunc("POST /todos", s.createTodo)
	s.mux.HandleFunc("GET /todos", s.listTodos)
	s.mux.HandleFunc("PATCH /todos/{id}", s.completeTodo)
	return s
}

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.mux.ServeHTTP(w, r)
}

func writeJSON(w http.ResponseWriter, status int, value any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(value)
}

func (s *Server) health(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, map[string]string{"status": "ok"})
}

func (s *Server) createTodo(w http.ResponseWriter, r *http.Request) {
	var input struct{
		Title string `json:"title"`
	}

	r.Body = http.MaxBytesReader(w,r.Body, 1<<20)
	decoder:= json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	if err:= decoder.Decode(&input); err!= nil{
		http.Error(w,"invalid JSON",http.StatusBadRequest)
		return
	}

	if err := decoder.Decode(&struct{}{}); err != io.EOF {
        http.Error(w, "expected one JSON object", http.StatusBadRequest)
        return
    }
    input.Title = strings.TrimSpace(input.Title)
    if input.Title == "" {
        http.Error(w, "title is required", http.StatusBadRequest)
        return
    }
    s.mu.Lock()
    todo := Todo{ID: s.nextID, Title: input.Title, Done: false}
    s.nextID++
    s.todos = append(s.todos, todo)
    s.mu.Unlock()
    writeJSON(w, http.StatusCreated, todo)
}


func (s *Server) listTodos(w http.ResponseWriter, r *http.Request) {
	s.mu.Lock()
    todos := make([]Todo, len(s.todos))
    copy(todos, s.todos)
    s.mu.Unlock()
    writeJSON(w, http.StatusOK, todos)
}

func (s *Server) completeTodo(w http.ResponseWriter, r *http.Request) {
	http.Error(w, "not implemented", http.StatusNotImplemented)
}
