package persistance

import (
	"fmt"
	"task_service/src/internal/core/task"
)

type TaskRepo struct {
	db *Database
}

func NewTaskRepo(d *Database) TaskRepo {
	return TaskRepo{db: d}
}

var emptyTask task.Task

func (t *TaskRepo) CreateNewTask(task1 task.TaskCreate) (task.Task, int, error) {
	var count int
	var createdTask task.Task
	tx, err := t.db.db.Begin()
	if err != nil {
		return emptyTask, count, err
	}
	defer tx.Rollback()
	query := `select count(*) from tasks where assigned_to=$1 and created_at < current_timestamp and current_timestamp < deadline`
	err = tx.QueryRow(query, task1.AssignedTo).Scan(&count)
	if err != nil {
		return emptyTask, count, err
	}
	fmt.Println("count = ", count)
	if count >= 3 {
		return emptyTask, count, fmt.Errorf("User Already have more than 3 tasks within the deadline")
	}
	query = `insert into tasks(name,assigned_to,description,priority,assigned_by,deadline) values($1,$2,$3,$4,$5,$6) returning id, name, assigned_by, assigned_to, description, task_status, created_at, deadline, priority`
	err = tx.QueryRow(query, task1.Name, task1.AssignedTo, task1.Description, task1.Priority, task1.AssignedBy, task1.Deadline).Scan(
		&createdTask.Id,
		&createdTask.Name,
		&createdTask.AssignedBy,
		&createdTask.AssignedTo,
		&createdTask.Description,
		&createdTask.TaskStatus,
		&createdTask.CreatedAt,
		&createdTask.Deadline,
		&createdTask.Priority,
	)
	if err != nil {
		return emptyTask, count, err
	}
	err = tx.Commit()
	if err != nil {
		return emptyTask, count, err
	}

	return createdTask, count, nil
}

func (t *TaskRepo) UpdateOldTask(task1 task.Task) (task.Task, error) {
	// we are sending assigned by to verify that the same user who assigned the task, can update the task, no other can can update, only the user who assigned the task, can update it.
	var existingTask task.Task
	query1 := `select name, description, task_status, priority, deadline from tasks where assigned_by=$1 and id=$2`
	err := t.db.db.QueryRow(query1, task1.AssignedBy, task1.Id).Scan(
		&existingTask.Name,
		&existingTask.Description,
		&existingTask.TaskStatus,
		&existingTask.Priority,
		&existingTask.Deadline,
	)
	if err != nil {
		return emptyTask, err
	}
	if task1.Name == "" {
		task1.Name = existingTask.Name
	}
	if task1.Description == "" {
		task1.Description = existingTask.Description
	}
	if task1.TaskStatus == "" {
		task1.TaskStatus = existingTask.TaskStatus
	}
	if task1.Priority == 0 {
		task1.Priority = existingTask.Priority
	}
	if task1.Deadline.IsZero() {
		task1.Deadline = existingTask.Deadline
	}
	query := `update tasks set name=$1, description=$2, task_status=$3, priority=$4, deadline=$5 where id=$6`
	//for this, before writing query, check in task struct, if any field is empty, then do not add it into the query.
	_, err = t.db.db.Exec(query, task1.Name, task1.Description, task1.TaskStatus, task1.Priority, task1.Deadline, task1.Id)
	if err != nil {
		return emptyTask, err
	}
	return task1, nil
}

func (t *TaskRepo) GetAllTaskDb(user_id int) ([]task.Task, error) {
	var tasks []task.Task
	query := "select * from tasks where assigned_to=$1"
	rows, err := t.db.db.Query(query, user_id)
	if err != nil {
		return []task.Task{}, err
	}
	defer rows.Close()
	for rows.Next() {
		var currentTask task.Task
		err := rows.Scan(&currentTask.Id, &currentTask.Name, &currentTask.AssignedTo, &currentTask.Description, &currentTask.TaskStatus, &currentTask.CreatedAt, &currentTask.Priority)
		if err != nil {
			return []task.Task{}, err
		}
		tasks = append(tasks, currentTask)
	}
	return tasks, nil
}

func (t *TaskRepo) DeleteThisTask(task1 task.Task) (task.Task, error) {
	var deleleted task.Task
	query := "delete from tasks where id=$1 and assigned_by=$2 returning id, name, assigned_by, assigned_to, description, task_status, created_at, deadline, priority"
	err := t.db.db.QueryRow(query, task1.Id, task1.AssignedBy).Scan(
		&deleleted.Id,
		&deleleted.Name,
		&deleleted.AssignedBy,
		&deleleted.AssignedTo,
		&deleleted.Description,
		&deleleted.TaskStatus,
		&deleleted.CreatedAt,
		&deleleted.Deadline,
		&deleleted.Priority,
	)
	if err != nil {
		return emptyTask, err
	}
	return deleleted, nil
}

// Get task by ID for notifications
func (t *TaskRepo) GetTaskByID(taskID int) (task.Task, error) {
	var taskData task.Task
	query := `SELECT id, name, assigned_to, description, task_status, created_at, priority, assigned_by, deadline 
			  FROM tasks WHERE id = $1`

	err := t.db.db.QueryRow(query, taskID).Scan(
		&taskData.Id,
		&taskData.Name,
		&taskData.AssignedTo,
		&taskData.Description,
		&taskData.TaskStatus,
		&taskData.CreatedAt,
		&taskData.Priority,
		&taskData.AssignedBy,
		&taskData.Deadline,
	)

	if err != nil {
		return task.Task{}, fmt.Errorf("task not found: %v", err)
	}

	return taskData, nil
}

func (t *TaskRepo) DeleteTask(taskID int) error {
	query := `DELETE FROM tasks WHERE id = $1`

	result, err := t.db.db.Exec(query, taskID)
	if err != nil {
		return fmt.Errorf("failed to delete task: %v", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %v", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("task not found")
	}

	return nil
}

// Get tasks by user ID (tasks assigned to or created by user)
func (t *TaskRepo) GetTasksByUserID(userID int) ([]task.Task, error) {
	query := `SELECT id, name, assigned_to, description, task_status, created_at, priority, assigned_by, deadline 
			  FROM tasks WHERE assigned_to = $1 OR assigned_by = $1`

	rows, err := t.db.db.Query(query, userID)
	if err != nil {
		return []task.Task{}, fmt.Errorf("failed to get user tasks: %v", err)
	}
	defer rows.Close()

	var tasks []task.Task
	for rows.Next() {
		var t task.Task
		err := rows.Scan(&t.Id, &t.Name, &t.AssignedTo, &t.Description, &t.TaskStatus, &t.CreatedAt, &t.Priority, &t.AssignedBy, &t.Deadline)
		if err != nil {
			return []task.Task{}, fmt.Errorf("failed to scan task: %v", err)
		}
		tasks = append(tasks, t)
	}

	if err = rows.Err(); err != nil {
		return []task.Task{}, fmt.Errorf("error iterating over rows: %v", err)
	}

	return tasks, nil
}

// Get all tasks (without user filtering)
func (t *TaskRepo) GetAllTask() ([]task.Task, error) {
	query := `SELECT id, name, assigned_to, description, task_status, created_at, priority, assigned_by, deadline 
			  FROM tasks ORDER BY created_at DESC`

	rows, err := t.db.db.Query(query)
	if err != nil {
		return []task.Task{}, fmt.Errorf("failed to get all tasks: %v", err)
	}
	defer rows.Close()

	var tasks []task.Task
	for rows.Next() {
		var t task.Task
		err := rows.Scan(&t.Id, &t.Name, &t.AssignedTo, &t.Description, &t.TaskStatus, &t.CreatedAt, &t.Priority, &t.AssignedBy, &t.Deadline)
		if err != nil {
			return []task.Task{}, fmt.Errorf("failed to scan task: %v", err)
		}
		tasks = append(tasks, t)
	}

	if err = rows.Err(); err != nil {
		return []task.Task{}, fmt.Errorf("error iterating over rows: %v", err)
	}

	return tasks, nil
}

func (t *TaskRepo) GetUserTaskDb(taskStatus task.TaskStatus) (int, task.TaskStatus, error) {
	var count int
	query := `select count(*) from tasks where assigned_to=$1 and created_at < $2 and deadline > $3`
	err := t.db.db.QueryRow(query, taskStatus.Id, taskStatus.Timeline, taskStatus.Timeline).Scan(&count)
	if err != nil {
		return count, task.TaskStatus{}, fmt.Errorf("Failed to get Tasks count for user : %v", err)
	}
	return count, taskStatus, nil
}
