package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

type APIServer struct {
	addr string
	db   *PostgresDB
}

func NewAPIServer(addr string, db *PostgresDB) *APIServer {
	return &APIServer{
		addr: addr,
		db:   db,
	}
}

func (s *APIServer) Run() error {
	router := mux.NewRouter()

	router.HandleFunc("/tasks/{id}/", s.getTaskById).Methods("GET")
	router.HandleFunc("/tasks/{id}/", s.deleteTaskById).Methods("DELETE")
	router.HandleFunc("/tasks/user/{id}/", s.getTaskForUser)
	router.HandleFunc("/tasks/", s.getTasks).Methods("GET")
	router.HandleFunc("/tasks/", s.createTask).Methods("POST")

	log.Printf("API server running on port: %s", s.addr)
	return http.ListenAndServe(s.addr, router)
}

func (s *APIServer) getTaskById(w http.ResponseWriter, r *http.Request) {
	taskID := mux.Vars(r)["id"]
	task, err := s.db.GetTaskByID(taskID)
	if err != nil {
		log.Println(err)
		return
	}
	json.NewEncoder(w).Encode(task)
}

func (s *APIServer) getTasks(w http.ResponseWriter, r *http.Request) {
	tasks, err := s.db.GetTasks()
	if err != nil {
		log.Println(err)
		return
	}
	json.NewEncoder(w).Encode(tasks)
}
func (s *APIServer) getTaskForUser(w http.ResponseWriter, r *http.Request) {
	userID := mux.Vars(r)["id"]

	id, err := strconv.Atoi(userID)

	if err != nil {
		log.Println(err)
		return
	}
	tasks, err := s.db.GetTaskByUserID(id)
	if err != nil {
		log.Println(err)
		return
	}
	json.NewEncoder(w).Encode(tasks)
}

func (s *APIServer) createTask(w http.ResponseWriter, r *http.Request) {
	requestBody := r.Body

	bs, err := io.ReadAll(requestBody)
	if err != nil {
		log.Println(err)
		return
	}

	var newCreateTask CreateTask

	err = json.Unmarshal(bs, &newCreateTask)

	if err != nil {
		log.Println(err)
		return
	}

	fmt.Println(newCreateTask)
	err = s.db.CreateTask(newCreateTask)
	if err != nil {
		log.Println(err)
		return
	}
	json.NewEncoder(w).Encode(newCreateTask)
}

func (s *APIServer) deleteTaskById(w http.ResponseWriter, r *http.Request) {
	taskID := mux.Vars(r)["id"]

	err := s.db.DeleteByTaskID(taskID)
	if err != nil {
		log.Println(err)
		return
	}
	json.NewEncoder(w).Encode(map[string]string{"message": "Task deleted successfully"})
}
