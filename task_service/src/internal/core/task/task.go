package task

import "time"

type Task struct {
	Id          int       `json:"id"`
	Name        string    `json:"name"`
	AssignedTo  int       `json:"assigned_to"`
	Description string    `json:"description"`
	TaskStatus  string    `json:"task_status"`
	CreatedAt   time.Time `json:"created_at"`
	Priority    int       `json:"priority"`
}
