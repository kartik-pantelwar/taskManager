package taskservice

import (
	"task_service/src/internal/adaptors/persistance"
	"task_service/src/internal/core/task"
)

type TaskService struct {
	taskRepo persistance.TaskRepo
}

func NewTaskService(taskRepo persistance.TaskRepo) TaskService {
	return TaskService{taskRepo: taskRepo}
}

func (t *TaskService) CreateTask(task task.Task) (task.Task, error) {
	newTask, err := t.taskRepo.CreateNewTask(task)
	return newTask, err
}

func (t *TaskService) UpdateTask(task task.Task) (task.Task, error) {
	newTask, err := t.taskRepo.UpdateOldTask(task)
	return newTask, err
}

func (t *TaskService) GetAllTask(user_id int) ([]task.Task, error) {
	allTask, err := t.taskRepo.GetAllTaskDb(user_id)
	return allTask, err
}

func (t *TaskService) DeleteTask(task task.Task) (task.Task, error) {
	newTask, err := t.taskRepo.DeleteThisTask(task)
	return newTask, err
}
