package main

import "time"

type Task struct {
	TaskID    string    `json:"task_id"`
	UserID    int       `json:"user_id"`
	Title     string    `json:"title"`
	Body      string    `json:"body"`
	Status    bool      `json:"status"`
	Timestamp time.Time `json:"timestamp"`
}

type Customer struct {
	UserID   int    `json:"user_id"`
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type CreateCustomer struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type CreateTask struct {
	TaskID string `json:"task_id"`
	UserID int    `json:"user_id"`
	Title  string `json:"title"`
	Body   string `json:"body"`
	Status bool   `json:"status"`
}
