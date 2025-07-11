package taskhandler

import (
	"encoding/json"
	"fmt"
	"log"
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

func (t *TaskHandler) CreateHandler(w http.ResponseWriter, r *http.Request) {
	userId, ok := r.Context().Value("user_id").(int)
	if !ok {
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(map[string]interface{}{"error": "user not found in context"})
		return
	}
	fmt.Println("user id=", userId)
	var task task.Task
	if err := json.NewDecoder(r.Body).Decode(&task); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}
	task.AssignedBy = userId

	createdTask, err := t.taskService.CreateTask(task)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}
	task = createdTask
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(task)

}

func (t *TaskHandler) UpdateHandler(w http.ResponseWriter, r *http.Request) {
	userId, ok := r.Context().Value("user_id").(int)
	if !ok {
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(map[string]interface{}{"error": "user not found in context"})
		return
	}
	var task task.Task
	if err := json.NewDecoder(r.Body).Decode(&task); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}
	task.AssignedBy = userId
	createdTask, err := t.taskService.UpdateTask(task)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}
	task = createdTask
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(task)

}

func (t *TaskHandler) GetMyHandler(w http.ResponseWriter, r *http.Request) {
	userId, ok := r.Context().Value("user_id").(int)
	if !ok {
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(map[string]interface{}{"error": "user not found in context"})
		return
	}
	myTask, err := t.taskService.GetAllTask(userId)
	if err != nil {
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	for _, i := range myTask {
		json.NewEncoder(w).Encode(i)
	}
}

func (t *TaskHandler) DeleteHandler(w http.ResponseWriter, r *http.Request) {
	task_idStr := chi.URLParam(r, "task-id")
	task_id, err := strconv.Atoi(task_idStr)
	if err != nil {
		log.Fatalf("Failed to get Task ID: %v", err)
	}
	userId, ok := r.Context().Value("user_id").(int)
	if !ok {
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(map[string]interface{}{"error": "user not found in context"})
		return
	}
	task := task.Task{Id: task_id, AssignedBy: userId}
	deletedTask, err := t.taskService.DeleteTask(task)
	if err != nil {
		log.Fatalf("Failed to delete task: %v", err)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(deletedTask)

}
