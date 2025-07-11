package persistance

import (
	"task_service/src/internal/core/task"
)

type TaskRepo struct {
	db *Database
}

func NewTaskRepo(d *Database) TaskRepo {
	return TaskRepo{db: d}
}

var emptyTask task.Task

func (t *TaskRepo) CreateNewTask(task task.Task) (task.Task, error) {
	var id int
	query := `insert into tasks(name,assigned_to,description,task_status,priority,assigned_by,deadline) values($1,$2,$3,$4,$5,$6,$7) returning id`
	//todo: Empty task_status part can be handled
	err := t.db.db.QueryRow(query, task.Name, task.AssignedTo, task.Description, task.TaskStatus, task.Priority, task.AssignedBy, task.Deadline).Scan(&id)
	if err != nil {
		return emptyTask, err
	}
	task.Id = id
	return task, nil
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
