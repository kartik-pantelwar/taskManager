package task

import "time"

type TaskEvent struct {
	EventType  string    `json:"event_type"`
	TaskID     int       `json:"task_id"`
	TaskName   string    `json:"task_name"`
	AssignedTo int       `json:"assigned_to"`
	AssignedBy int       `json:"assigned_by"`
	Timestamp  time.Time `json:"timestamp"`
}
