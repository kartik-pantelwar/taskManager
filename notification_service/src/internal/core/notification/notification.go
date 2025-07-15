package notification

import "time"

type Notification struct {
	ID         string    `json:"id"`
	TaskID     int       `json:"task_id"`
	Action     string    `json:"action"`
	TaskName   string    `json:"task_name"`
	UserID     int       `json:"user_id"`     // User who receives the notification (assigned_to)
	AssignedBy int       `json:"assigned_by"` // User who assigned the task (logged-in user)
	Message    string    `json:"message"`
	Timestamp  time.Time `json:"timestamp"`
}
