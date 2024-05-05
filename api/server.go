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

	router.HandleFunc("/tasks/{id}/", makeHTTPHandleFunc(s.getTaskById)).Methods("GET")
	router.HandleFunc("/tasks/{id}/", makeHTTPHandleFunc(s.deleteTaskById)).Methods("DELETE")
	router.HandleFunc("/tasks/user/{id}/", makeHTTPHandleFunc(s.getTaskForUser))
	router.HandleFunc("/tasks/", makeHTTPHandleFunc(s.getTasks)).Methods("GET")
	router.HandleFunc("/tasks/", makeHTTPHandleFunc(s.createTask)).Methods("POST")
	router.HandleFunc("/tasks/{id}/", makeHTTPHandleFunc(s.updateTask)).Methods("PUT")
	router.HandleFunc("/register/", makeHTTPHandleFunc(s.registerUser)).Methods("POST")
	router.HandleFunc("/login/", makeHTTPHandleFunc(s.loginCustomer)).Methods("POST")

	log.Printf("API server running on port: %s", s.addr)
	return http.ListenAndServe(s.addr, router)
}

func (s *APIServer) getTaskById(w http.ResponseWriter, r *http.Request) error {
	taskID, err := getID("id", r)
	if err != nil {
		log.Println(err)
		return err
	}

	task, err := s.db.GetTaskByID(taskID)
	if err != nil {
		log.Println(err)
		return err
	}
	return writeJSON(w, http.StatusOK, task)
}

func (s *APIServer) getTasks(w http.ResponseWriter, r *http.Request) error {
	tasks, err := s.db.GetTasks()
	if err != nil {
		log.Println(err)
		return err
	}
	return writeJSON(w, http.StatusOK, tasks)
}
func (s *APIServer) getTaskForUser(w http.ResponseWriter, r *http.Request) error {
	userID := mux.Vars(r)["id"]

	id, err := strconv.Atoi(userID)

	if err != nil {
		log.Println(err)
		return err
	}
	tasks, err := s.db.GetTaskByUserID(id)
	if err != nil {
		log.Println(err)
		return err
	}
	return writeJSON(w, http.StatusOK, tasks)
}

func (s *APIServer) createTask(w http.ResponseWriter, r *http.Request) error {
	requestBody := r.Body

	bs, err := io.ReadAll(requestBody)
	if err != nil {
		log.Println(err)
		return err
	}

	var newCreateTask types.CreateTask

	tID := uuid.NewV4()
	newCreateTask.TaskID = tID.String()

	err = json.Unmarshal(bs, &newCreateTask)

	if err != nil {
		log.Println(err)
		return err
	}

	err = s.db.CreateTask(newCreateTask)
	if err != nil {
		log.Println(err)
		return err
	}
	return writeJSON(w, http.StatusCreated, newCreateTask)
}

func (s *APIServer) deleteTaskById(w http.ResponseWriter, r *http.Request) error {
	taskID, err := getID("id", r)
	if err != nil {
		log.Println(err)
		return err
	}

	err = s.db.DeleteByTaskID(taskID)
	if err != nil {
		log.Println(err)
		return err
	}
	return writeJSON(w, http.StatusOK, map[string]string{"message": "Task deleted successfully"})
}

func (s *APIServer) updateTask(w http.ResponseWriter, r *http.Request) error {
	taskID, err := getID("id", r)
	if err != nil {
		log.Println(err)
		return err
	}

	bs, err := io.ReadAll(r.Body)

	if err != nil {
		log.Println(err)
		return err
	}

	var updateValues types.UpdateTask
	err = json.Unmarshal(bs, &updateValues)

	if err != nil {
		log.Println(err)
		return err
	}

	task, err := s.db.UpdateTask(taskID, updateValues)

	if err != nil {
		log.Println(err)
		return err
	}
	return writeJSON(w, http.StatusOK, task)
}

func (s *APIServer) registerUser(w http.ResponseWriter, r *http.Request) error {
	bs, err := io.ReadAll(r.Body)

	if err != nil {
		log.Println(err)
		return err
	}

	var newCustomer types.RegisterRequest

	err = json.Unmarshal(bs, &newCustomer)
	if err != nil {
		log.Println(err)
		return err
	}

	err = newCustomer.CryptPassword()
	if err != nil {
		log.Println(err)
		return err
	}
	err = s.db.CreateCustomer(newCustomer)
	if err != nil {
		log.Println(err)
		return err
	}
	return writeJSON(w, http.StatusCreated, map[string]string{"message": "User created"})
}

func (s *APIServer) loginCustomer(w http.ResponseWriter, r *http.Request) error {
	bs, err := io.ReadAll(r.Body)
	if err != nil {
		log.Println(err)
		return err
	}

	var loginRequest types.LoginRequest

	err = json.Unmarshal(bs, &loginRequest)

	if err != nil {
		log.Println(err)
		return err
	}
	customer, err := s.db.GetCustomerByUsername(loginRequest.Username)

	if err != nil {
		log.Println(err)
		return err
	}

	if err = bcrypt.CompareHashAndPassword([]byte(customer.Password), []byte(loginRequest.Password)); err != nil {
		log.Println(err)
		return err
	}

	log.Println("Login successfull")
	return writeJSON(w, http.StatusCreated, map[string]string{"sessionid": uuid.NewV4().String()})
}

func getID(name string, r *http.Request) (string, error) {
	if key, ok := mux.Vars(r)[name]; ok {
		return key, nil
	}
	return "", fmt.Errorf("key was not found: %s", name)
}

func writeJSON(w http.ResponseWriter, status int, v any) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	return json.NewEncoder(w).Encode(v)
}

type apiFunc func(w http.ResponseWriter, r *http.Request) error

func makeHTTPHandleFunc(f apiFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := f(w, r); err != nil {
			writeJSON(w, http.StatusBadRequest, err)
		}
	}
}
