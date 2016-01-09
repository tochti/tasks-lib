package tasks

import (
	"database/sql"
	"fmt"

	"github.com/tochti/gin-gum/gumauth"

	"gopkg.in/gorp.v1"
)

var (
	UserTable       = "User"
	TasksTable      = "Task"
	TasksUsersTable = "TasksUsers"
)

func Q(q string, s ...interface{}) string {
	return fmt.Sprintf(q, s...)
}

func GorpInit(db *sql.DB) *gorp.DbMap {
	dbMap := &gorp.DbMap{
		Db:      db,
		Dialect: gorp.MySQLDialect{"InnoDB", "UTF8"},
	}

	dbMap.AddTable(gumauth.User{}).SetKeys(true, "ID")
	dbMap.AddTable(Task{}).SetKeys(true, "id")
	dbMap.AddTable(TasksUsers{}).SetKeys(false, "user_id", "task_id")

	return dbMap
}

// Give back all tasks for a given user. Included task which are done.
// Orderd by name
func ReadAllTasksUser(db *gorp.DbMap, userID int64) ([]Task, error) {
	tasks := []Task{}
	q := Q(`SELECT task.id, task.name, task.desc, task.expires, task.done
		FROM %v tasksUsers, %v task
		WHERE tasksUsers.user_id=%v
		AND task.id=tasksUsers.task_id
		ORDER BY task.name`, TasksUsersTable, TasksTable, userID)
	_, err := db.Select(&tasks, q)
	if err != nil {
		return []Task{}, err
	}

	return tasks, nil
}

// Give back a task for a given id and user id
func ReadOneTaskUser(db *gorp.DbMap, taskID, userID int64) (Task, error) {
	task := Task{}
	q := Q(`SELECT task.id, task.name, task.desc, task.expires, task.done
		FROM %v tasksUsers, %v task
		WHERE tasksUsers.user_id=%v
		AND tasksUsers.task_id=?
		AND task.id=tasksUsers.task_id`,
		TasksUsersTable, TasksTable, userID)

	err := db.SelectOne(&task, q, taskID)
	if err != nil {
		return Task{}, err
	}

	return task, nil
}

// Create a new task an bind it to an user
func CreateTaskUser(db *gorp.DbMap, task *Task, userID int64) error {
	trans, err := db.Begin()
	if err != nil {
		return err
	}

	trans.Insert(task)
	trans.Insert(TasksUsers{userID, task.ID})

	return trans.Commit()
}

// Delete a task by id and dissolve user connection
func DeleteTaskUser(db *gorp.DbMap, task *Task, userID int64) error {
	trans, err := db.Begin()
	if err != nil {
		return err
	}

	trans.Delete(TasksUsers{userID, task.ID})
	trans.Delete(task)

	return trans.Commit()
}
