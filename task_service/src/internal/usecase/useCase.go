package task

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"strconv"
	"task_service/src/internal/adaptors/persistance"
	"task_service/src/internal/adaptors/redis/notification"
	"task_service/src/internal/core/task"
	pb "task_service/src/internal/interfaces/input/grpc/generated/generated"
	"time"
)

type TaskService struct {
	taskRepo            persistance.TaskRepo
	notificationService *notification.NotificationService
	grpcClient          pb.SessionValidatorClient
}

// Constructor with notification service and gRPC client
func NewTaskService(taskRepo persistance.TaskRepo, notificationService *notification.NotificationService, grpcClient pb.SessionValidatorClient) TaskService {
	return TaskService{
		taskRepo:            taskRepo,
		notificationService: notificationService,
		grpcClient:          grpcClient,
	}
}

// CreateTask + notification
func (t *TaskService) CreateTask(ctx context.Context, taskData task.Task, userID int) (task.Task, int, error) {
	// Validate if the assigned_to user exists before creating the task
	if taskData.AssignedTo != 0 {
		userExistsReq := &pb.ValidateUserRequest{
			UserId: strconv.Itoa(taskData.AssignedTo),
		}

		userExistsResp, err := t.grpcClient.ValidateUser(ctx, userExistsReq)
		if err != nil {
			return task.Task{}, 0, fmt.Errorf("failed to validate user: %v", err)
		}

		if !userExistsResp.Status {
			return task.Task{}, 0, fmt.Errorf("User does not exist")
		}
	}

	createdTask, count, err := t.taskRepo.CreateNewTask(taskData)
	if err != nil {
		return task.Task{}, count, err
	}

	t.publishTaskEvent("task_created", createdTask, userID)
	return createdTask, count, nil
}

// UpdateTask + notification
func (t *TaskService) UpdateTask(ctx context.Context, taskData task.Task, userID int) (task.Task, error) {
	// Validate if the assigned_to user exists before updating the task
	if taskData.AssignedTo != 0 {
		userExistsReq := &pb.ValidateUserRequest{
			UserId: strconv.Itoa(taskData.AssignedTo),
		}

		userExistsResp, err := t.grpcClient.ValidateUser(ctx, userExistsReq)
		if err != nil {
			return task.Task{}, fmt.Errorf("failed to validate user: %v", err)
		}

		if !userExistsResp.Status {
			return task.Task{}, fmt.Errorf("User does not exist")
		}
	}

	// updation
	updatedTask, err := t.taskRepo.UpdateOldTask(taskData)
	if err != nil {
		return task.Task{}, err
	}

	// Send notification
	t.publishTaskEvent("task_updated", updatedTask, userID)
	return updatedTask, nil
}

// DeleteTask + notification
func (t *TaskService) DeleteTask(ctx context.Context, taskID int, userID int) error {
	// getting task details for notifications
	taskData, err := t.taskRepo.GetTaskByID(taskID)
	if err != nil {
		return err
	}

	// deleting task
	err = t.taskRepo.DeleteTask(taskID)
	if err != nil {
		return err
	}

	// Send notification
	t.publishTaskEvent("task_deleted", taskData, userID)
	return nil
}

// created but not used
func (t *TaskService) GetAllTasks() ([]task.Task, error) {
	return t.taskRepo.GetAllTask()
}

func (t *TaskService) GetTasksByUserID(userID int) ([]task.Task, error) {
	return t.taskRepo.GetTasksByUserID(userID)
}

func (t *TaskService) publishTaskEvent(eventType string, task1 task.Task, userID int) {
	if t.notificationService == nil {
		log.Printf("Notification service not available, skipping event publication")
		return
	}

	event := task.TaskEvent{
		EventType:  eventType,
		TaskID:     task1.Id,
		TaskName:   task1.Name,
		AssignedTo: task1.AssignedTo,
		AssignedBy: userID,
		Timestamp:  time.Now(),
	}

	eventJSON, err := json.Marshal(event)
	if err != nil {
		log.Printf("Failed to marshal task Event: %v", err)
		return
	}

	// Use the notification service to publish the event
	err = t.notificationService.PublishEvent(context.Background(), "task_events", eventJSON)
	if err != nil {
		log.Printf("Failed to publish task event: %v", err)
	} else {
		log.Printf("Published %s event for task %d", eventType, task1.Id)
	}
}

func (t *TaskService) GetUserTasks(taskStatus task.TaskStatus) (int, task.TaskStatus, error) {
	var newStatus task.TaskStatus
	count, newStatus, err := t.taskRepo.GetUserTaskDb(taskStatus)
	return count, newStatus, err
}
