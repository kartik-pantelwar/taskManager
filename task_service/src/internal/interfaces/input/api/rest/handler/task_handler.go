package handler

import (
	"context"
	"encoding/json"
	"net/http"
	"strconv"
	"task_service/src/internal/core/task"
	taskservice "task_service/src/internal/usecase"
	errorhandling "task_service/src/pkg/error_handling"
	pkgresponse "task_service/src/pkg/response"

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
		errorhandling.HandleError(w, "User Not Found in Context", http.StatusUnauthorized)
		return
	}

	var taskData task.TaskCreate
	err := json.NewDecoder(r.Body).Decode(&taskData)
	if err != nil {
		errorhandling.HandleError(w, "Invalid Request Body", http.StatusBadRequest)
		return
	}

	// Set the assigned_by field to the current user
	taskData.AssignedBy = userId

	// Pass userId to CreateTask for notifications
	createdTask, count, err := t.taskService.CreateTask(context.Background(), taskData, userId)
	if err != nil {
		errorhandling.HandleError(w, err.Error(), http.StatusInternalServerError)
		return
	}

	response := pkgresponse.StandardResponse{
		Status:  "SUCCESS",
		Message: "Task Created Successfully",
		Data: map[string]interface{}{
			"task":  createdTask,
			"count": count + 1,
		},
	}
	pkgresponse.WriteResponse(w, http.StatusOK, response)
}

func (t *TaskHandler) Update(w http.ResponseWriter, r *http.Request) {

	userId, ok := r.Context().Value("user_id").(int)
	if !ok {
		errorhandling.HandleError(w, "User Not Found in Context", http.StatusUnauthorized)
		return
	}

	var taskData task.Task
	err := json.NewDecoder(r.Body).Decode(&taskData)
	if err != nil {
		errorhandling.HandleError(w, "Invalid Request Body", http.StatusBadRequest)
		return
	}

	// Set the assigned_by field to verify authorization
	taskData.AssignedBy = userId

	// ONLY CHANGE: Pass userId to UpdateTask for notifications
	updatedTask, err := t.taskService.UpdateTask(context.Background(), taskData, userId)
	if err != nil {
		errorhandling.HandleError(w, err.Error(), http.StatusInternalServerError)
		return
	}

	response := pkgresponse.StandardResponse{
		Status:  "SUCCESS",
		Message: "Task Updated Successfully",
		Data:    updatedTask,
	}
	pkgresponse.WriteResponse(w, http.StatusOK, response)
}

func (t *TaskHandler) GetMy(w http.ResponseWriter, r *http.Request) {
	userId, ok := r.Context().Value("user_id").(int)
	if !ok {
		errorhandling.HandleError(w, "User Not Found in Context", http.StatusUnauthorized)
		return
	}

	tasks, err := t.taskService.GetTasksByUserID(userId)
	if err != nil {
		errorhandling.HandleError(w, err.Error(), http.StatusInternalServerError)
		return
	}

	response := pkgresponse.StandardResponse{
		Status:  "SUCCESS",
		Message: "User Tasks Retrieved Successfully",
		Data: map[string]interface{}{
			"tasks": tasks,
			"count": len(tasks),
		},
	}
	pkgresponse.WriteResponse(w, http.StatusOK, response)
}

func (t *TaskHandler) Delete(w http.ResponseWriter, r *http.Request) {
	userId, ok := r.Context().Value("user_id").(int)
	if !ok {
		errorhandling.HandleError(w, "User Not Found in Context", http.StatusUnauthorized)
		return
	}

	taskIDStr := chi.URLParam(r, "id")
	taskID, err := strconv.Atoi(taskIDStr)
	if err != nil {
		errorhandling.HandleError(w, "Invalid Task ID", http.StatusBadRequest)
		return
	}

	// ONLY CHANGE: Pass userId to DeleteTask for notifications
	err = t.taskService.DeleteTask(context.Background(), taskID, userId)
	if err != nil {
		errorhandling.HandleError(w, err.Error(), http.StatusInternalServerError)
		return
	}

	response := pkgresponse.StandardResponse{
		Status:  "SUCCESS",
		Message: "Task Deleted Successfully",
	}
	pkgresponse.WriteResponse(w, http.StatusOK, response)
}

func (t *TaskHandler) GetStatus(w http.ResponseWriter, r *http.Request) {
	var taskStatus task.TaskStatus
	err := json.NewDecoder(r.Body).Decode(&taskStatus)
	if err != nil {
		errorhandling.HandleError(w, "Invalid Request Body", http.StatusBadRequest)
		return
	}

	taskCount, newStatus, err := t.taskService.GetUserTasks(taskStatus)
	if err != nil {
		errorhandling.HandleError(w, "Failed to Get Tasks Status", http.StatusInternalServerError)
		return
	}

	response := pkgresponse.StandardResponse{
		Status:  "SUCCESS",
		Message: "Task Status Retrieved Successfully",
		Data: map[string]interface{}{
			"task_count": taskCount,
			"status":     newStatus,
		},
	}
	pkgresponse.WriteResponse(w, http.StatusOK, response)
}
