package handler

import (
	"context"
	"encoding/json"
	"net/http"
	"strconv"
	"task_service/src/internal/core/task"
	taskservice "task_service/src/internal/usecase"

	"github.com/go-chi/chi/v5"
)

type TaskHandler struct {
	taskService taskservice.TaskService
}

func NewTaskHandler(taskcase taskservice.TaskService) TaskHandler {
	return TaskHandler{
		taskService: taskcase,
	}
}

func (t *TaskHandler) Create(w http.ResponseWriter, r *http.Request) {
	userId, ok := r.Context().Value("user_id").(int)
	if !ok {
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(map[string]interface{}{"error": "user not found in context"})
		return
	}

	var taskData task.Task
	err := json.NewDecoder(r.Body).Decode(&taskData)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]interface{}{"error": "invalid request body"})
		return
	}

	// Set the assigned_by field to the current user
	taskData.AssignedBy = userId

	// ONLY CHANGE: Pass userId to CreateTask for notifications
	createdTask, count, err := t.taskService.CreateTask(context.Background(), taskData, userId)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]interface{}{"error": err.Error()})
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"task":  createdTask,
		"count": count+1,
	})
}

func (t *TaskHandler) Update(w http.ResponseWriter, r *http.Request) {
	// Your original code (back to int)
	userId, ok := r.Context().Value("user_id").(int)
	if !ok {
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(map[string]interface{}{"error": "user not found in context"})
		return
	}

	var taskData task.Task
	err := json.NewDecoder(r.Body).Decode(&taskData)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]interface{}{"error": "invalid request body"})
		return
	}

	// Set the assigned_by field to verify authorization
	taskData.AssignedBy = userId

	// ONLY CHANGE: Pass userId to UpdateTask for notifications
	updatedTask, err := t.taskService.UpdateTask(context.Background(), taskData, userId)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]interface{}{"error": err.Error()})
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"task": updatedTask,
	})
}

func (t *TaskHandler) GetMy(w http.ResponseWriter, r *http.Request) {
	// Your original code (back to int)
	userId, ok := r.Context().Value("user_id").(int)
	if !ok {
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(map[string]interface{}{"error": "user not found in context"})
		return
	}

	tasks, err := t.taskService.GetTasksByUserID(userId)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]interface{}{"error": err.Error()})
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"tasks": tasks,
		"count": len(tasks),
	})
}

func (t *TaskHandler) Delete(w http.ResponseWriter, r *http.Request) {
	// Your original code (back to int)
	userId, ok := r.Context().Value("user_id").(int)
	if !ok {
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(map[string]interface{}{"error": "user not found in context"})
		return
	}

	// Your original URL parameter extraction
	taskIDStr := chi.URLParam(r, "id")
	taskID, err := strconv.Atoi(taskIDStr)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]interface{}{"error": "invalid task ID"})
		return
	}

	// ONLY CHANGE: Pass userId to DeleteTask for notifications
	err = t.taskService.DeleteTask(context.Background(), taskID, userId)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]interface{}{"error": err.Error()})
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"message": "Task deleted successfully",
	})
}

func (t *TaskHandler) GetStatus(w http.ResponseWriter, r *http.Request) {
	var taskStatus task.TaskStatus
	json.NewDecoder(r.Body).Decode(&taskStatus)
	taskCount, newStatus, err := t.taskService.GetUserTasks(taskStatus)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]interface{}{"error": "Failed to get tasks status"})
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"No. Of Tasks": taskCount,
	})
	json.NewEncoder(w).Encode(newStatus)
}
