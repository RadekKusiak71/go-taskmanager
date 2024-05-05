package main

import (
	"log"
	"time"

	"golang.org/x/crypto/bcrypt"
)

type Task struct {
	TaskID    string    `json:"task_id"`
	UserID    int       `json:"user_id"`
	Title     string    `json:"title"`
	Body      string    `json:"body"`
	Status    bool      `json:"status"`
	Timestamp time.Time `json:"timestamp"`
}
type CreateTask struct {
	TaskID string `json:"task_id"`
	UserID int    `json:"user_id"`
	Title  string `json:"title"`
	Body   string `json:"body"`
	Status bool   `json:"status"`
}

type UpdateTask struct {
	Title  string `json:"title"`
	Body   string `json:"body"`
	Status bool   `json:"status"`
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

func (c *CreateCustomer) cryptPassword() error {
	cryptedPS, err := bcrypt.GenerateFromPassword([]byte(c.Password), bcrypt.DefaultCost)

	if err != nil {
		log.Println(err)
		return err
	}
	c.Password = string(cryptedPS)
	return nil
}
