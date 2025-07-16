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
type TaskEvent struct {
	EventType  string    `json:"event_type"`
	TaskID     int       `json:"task_id"`
	TaskName   string    `json:"task_name"`
	AssignedTo int       `json:"assigned_to"`
	AssignedBy int       `json:"assigned_by"`
	Timestamp  time.Time `json:"timestamp"`
}

type TaskStatus struct {
	Id       int       `json:"user_id"`
	Timeline time.Time `json:"timeline"`
}

type TaskCreate struct {
	Name        string `json:"name"`
	AssignedBy  int
	AssignedTo  int       `json:"assigned_to"`
	Description string    `json:"description"`
	Deadline    time.Time `json:"deadline"`
	Priority    int       `json:"priority"`
}
