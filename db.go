package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

type PostgresDB struct {
	db *sql.DB
}

func NewPostgreSQLStorage() (*PostgresDB, error) {

	db, err := sql.Open("postgres", GetConnectionString())
	if err != nil {
		log.Printf("Problem while connecting to db: %s", err)
		return nil, err
	}

	if err = db.Ping(); err != nil {
		return nil, err
	}

	log.Println("DB: successfully connected!")

	return &PostgresDB{
		db,
	}, nil
}

func GetConnectionString() string {
	godotenv.Load()
	psUSER := os.Getenv("PS_USER")
	psPSW := os.Getenv("PS_PSW")
	psNAME := os.Getenv("PS_NAME")

	return fmt.Sprintf("user=%s password=%s dbname=%s sslmode=disable", psUSER, psPSW, psNAME)
}

func (s PostgresDB) CreateCustomerTable() error {
	query := `CREATE TABLE customer (
		user_id SERIAL NOT NULL PRIMARY KEY,
		username VARCHAR(255),
		email VARCHAR(255),
		password VARCHAR(255)
		);`
	_, err := s.db.Exec(query)
	if err != nil {
		return err
	}
	return nil
}

func (s PostgresDB) CreateTaskTable() error {
	query := `CREATE TABLE taks(
		task_id VARCHAR(255),
		user_id int,
		title VARCHAR(255),
		body VARCHAR(255),
		status BOOLEAN,
		created_at timestamp
		);`
	_, err := s.db.Exec(query)
	if err != nil {
		return err
	}
	return nil
}

func (s PostgresDB) CreateCustomer(c CreateCustomer) error {
	query := `INSERT INTO customer (username,email,password)
		Values
		($1,$2,$3)`

	_, err := s.db.Query(query, c.Username, c.Email, c.Password)
	if err != nil {
		return err
	}
	return nil
}

func (s PostgresDB) CreateTask(t CreateTask) error {
	query := `INSERT INTO task (task_id,user_id,title,body,status) VALUES
		($1,$2,$3,$4,$5)`

	_, err := s.db.Query(query, t.TaskID, t.UserID, t.Title, t.Body, t.Status)
	if err != nil {
		return err
	}
	return nil
}

func (s PostgresDB) GetTaskByID(taskID string) (*Task, error) {
	rows, err := s.db.Query(`SELECT * FROM task WHERE task_id=$1`, taskID)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		return convertToTask(rows)
	}

	return nil, fmt.Errorf("account %s not found", taskID)
}

func (s PostgresDB) GetTasks() ([]*Task, error) {
	rows, err := s.db.Query(`SELECT * FROM task`)
	if err != nil {
		return nil, err
	}
	var tasks []*Task
	for rows.Next() {
		task, err := convertToTask(rows)
		if err != nil {
			return nil, err
		}
		tasks = append(tasks, task)
	}
	return tasks, nil
}

func (s PostgresDB) GetTaskByUserID(userID int) ([]*Task, error) {
	rows, err := s.db.Query(`SELECT * FROM task WHERE user_id=$1`, userID)
	if err != nil {
		return nil, err
	}
	var tasks []*Task
	for rows.Next() {
		task, err := convertToTask(rows)
		if err != nil {
			return nil, err
		}
		tasks = append(tasks, task)
	}
	return tasks, nil
}

func (s PostgresDB) DeleteByTaskID(taskID string) error {
	_, err := s.db.Exec("DELETE FROM task WHERE task_id=$1", taskID)
	if err != nil {
		return err
	}
	return nil
}

func convertToTask(rows *sql.Rows) (*Task, error) {
	task := new(Task)
	err := rows.Scan(
		&task.TaskID,
		&task.UserID,
		&task.Title,
		&task.Body,
		&task.Status,
		&task.Timestamp,
	)
	return task, err
}