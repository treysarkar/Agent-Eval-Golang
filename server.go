package main

import (
	"encoding/json"
	"net/http"
	"sync"
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
	http.Error(w, "Not Implemented", http.StatusNotImplemented)
}

func (s *Server) createTodo(w http.ResponseWriter, r *http.Request) {
	http.Error(w, "not implemented", http.StatusNotImplemented)
}

func (s *Server) listTodos(w http.ResponseWriter, r *http.Request) {
	http.Error(w, "not implemented", http.StatusNotImplemented)
}

func (s *Server) completeTodo(w http.ResponseWriter, r *http.Request) {
	http.Error(w, "not implemented", http.StatusNotImplemented)
}
