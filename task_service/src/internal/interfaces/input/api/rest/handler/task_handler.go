package taskhandler

import (
	"encoding/json"
	"net/http"
	"task_service/src/internal/core/task"
	taskservice "task_service/src/internal/usecase"
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
	var task task.Task

	if err := json.NewDecoder(r.Body).Decode(&task); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}

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

func (t *TaskHandler) Update(w http.ResponseWriter, r *http.Request) {
	var task task.Task

	if err := json.NewDecoder(r.Body).Decode(&task); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}

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

func (t *TaskHandler) GetAll(w http.ResponseWriter, r *http.Request) {
	var allTask []task.Task
	allTask, err := t.taskService.GetAllTask()
	if err != nil {
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	for _, i := range allTask {
		json.NewEncoder(w).Encode(i)
	}
}
