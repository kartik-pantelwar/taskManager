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

func (t *TaskRepo) CreateNewTask(task task.Task) (task.Task, error) {
	var id int
	query := `insert into tasks(name,assigned_to,description,task_status,priority) values($1,$2,$3,$4,$5) returning id`
	//todo: Empty task_status part can be handled
	err := t.db.db.QueryRow(query, task.Name, task.AssignedTo, task.Description, task.TaskStatus, task.Priority).Scan(&id)
	if err != nil {
		return emptyTask, fmt.Errorf("Failed to Add a New Task in Database")
	}
	task.Id = id
	return task, nil
}

func (t *TaskRepo) UpdateOldTask(task task.Task) (task.Task, error) {
	query := `update tasks set name=$1, assigned_to=$2, description=$3, task_status=$4, priority=$5 where id=$6`
	//todo: Empty task_status part can be handled
	//for this, before writing query, check in task struct, if any field is empty, then do not add it into the query.
	_, err := t.db.db.Exec(query, task.Name, task.AssignedTo, task.Description, task.TaskStatus, task.Priority, task.Id)
	if err != nil {
		return emptyTask, fmt.Errorf("Failed to Add a New Task in Database")
	}
	// task.Id = id
	return task, nil
}

func (t *TaskRepo) GetAllTaskDb() ([]task.Task, error) {
	var tasks []task.Task
	query := "select * from tasks"
	rows, err := t.db.db.Query(query)
	if err != nil {
		return []task.Task{}, fmt.Errorf("Failed to get Task from Database")
	}
	defer rows.Close()
	for rows.Next() {
		var currentTask task.Task
		err := rows.Scan(&currentTask.Id, &currentTask.Name, &currentTask.AssignedTo, &currentTask.Description, &currentTask.TaskStatus, &currentTask.CreatedAt, &currentTask.Priority)
		if err != nil {
			return []task.Task{}, fmt.Errorf("Failed to display tasks")
		}
		tasks = append(tasks, currentTask)
	}
	return tasks, nil
}
