package task

import "time"

type Task struct {
	Id          int       `json:"id"`
	Name        string    `json:"name"`
	AssignedBy  int       `json:"assigned_by"`
	AssignedTo  int       `json:"assigned_to"`
	Description string    `json:"description"`
	TaskStatus  string    `json:"task_status"`
	CreatedAt   time.Time `json:"created_at"`
	Deadline    time.Time `json:"deadline"`
	Priority    int       `json:"priority"`
}
