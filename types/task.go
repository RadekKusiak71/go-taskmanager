package types

import "time"

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
