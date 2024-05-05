package api

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"

	"github.com/RadekKusiak71/taskmanager/db"
	"github.com/RadekKusiak71/taskmanager/types"
	"golang.org/x/crypto/bcrypt"

	"github.com/gorilla/mux"
	uuid "github.com/satori/go.uuid"
)

type APIServer struct {
	addr string
	db   *db.PostgresDB
}

func NewAPIServer(addr string, db *db.PostgresDB) *APIServer {
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
	router.HandleFunc("/tasks/{id}/", s.updateTask).Methods("PUT")
	router.HandleFunc("/register/", s.registerUser).Methods("POST")
	router.HandleFunc("/login/", s.loginCustomer).Methods("POST")

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

	var newCreateTask types.CreateTask

	tID := uuid.NewV4()
	newCreateTask.TaskID = tID.String()

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

func (s *APIServer) updateTask(w http.ResponseWriter, r *http.Request) {
	taskID := mux.Vars(r)["id"]
	bs, err := io.ReadAll(r.Body)

	if err != nil {
		log.Println(err)
		return
	}

	var updateValues types.UpdateTask
	err = json.Unmarshal(bs, &updateValues)

	if err != nil {
		log.Println(err)
		return
	}

	task, err := s.db.UpdateTask(taskID, updateValues)

	if err != nil {
		log.Println(err)
		return
	}
	json.NewEncoder(w).Encode(task)
}

func (s *APIServer) registerUser(w http.ResponseWriter, r *http.Request) {
	bs, err := io.ReadAll(r.Body)

	if err != nil {
		log.Println(err)
		return
	}

	var newCustomer types.RegisterRequest

	err = json.Unmarshal(bs, &newCustomer)
	if err != nil {
		log.Println(err)
		return
	}

	err = newCustomer.CryptPassword()
	if err != nil {
		log.Println(err)
		return
	}
	err = s.db.CreateCustomer(newCustomer)
	if err != nil {
		log.Println(err)
		return
	}
	json.NewEncoder(w).Encode(map[string]string{"message": "User created"})
}

func (s *APIServer) loginCustomer(w http.ResponseWriter, r *http.Request) {
	bs, err := io.ReadAll(r.Body)
	if err != nil {
		log.Println(err)
		return
	}

	var loginRequest types.LoginRequest

	err = json.Unmarshal(bs, &loginRequest)

	if err != nil {
		log.Println(err)
		return
	}
	customer, err := s.db.GetCustomerByUsername(loginRequest.Username)

	if err != nil {
		log.Println(err)
		return
	}

	if err = bcrypt.CompareHashAndPassword([]byte(customer.Password), []byte(loginRequest.Password)); err != nil {
		log.Println(err)
		return
	}

	log.Println("Login successfull")
	json.NewEncoder(w).Encode(map[string]string{"sessionid": uuid.NewV4().String()})
}
