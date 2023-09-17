package main

import (
	"database/sql"
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

type taskDB struct {
	db      *sql.DB
	dataDir string
}

type task struct {
	ID      uint
	Name    string
	Project string
	Status  string
	Created time.Time
}

type status int

const (
	todo status = iota
	inProgress
	done
)

func (s status) String() string {
	return [...]string{"todo", "in progress", "done"}[s]
}

func (orig *task) merge(t task) {
	uValues := reflect.ValueOf(&t).Elem()
	oValues := reflect.ValueOf(orig).Elem()
	for i := 0; i < uValues.NumField(); i++ {
		uField := uValues.Field(i).Interface()
		if oValues.CanSet() {
			if v, ok := uField.(int64); ok && uField != 0 {
				oValues.Field(i).SetInt(v)
			}
			if v, ok := uField.(string); ok && uField != "" {
				oValues.Field(i).SetString(v)
			}
		}
	}
}

func (t *taskDB) createTable() error {
	_, err := t.db.Exec(`CREATE TABLE IF NOT EXISTS tasks (
		id INTEGER,
		name TEXT NOT NULL,
		project TEXT,
		status TEXT,
		created DATETIME,
		PRIMARY KEY (id AUTOINCREMENT)
	)`)
	return err
}

func (t *taskDB) taskTableExists() bool {
	if _, err := t.db.Query("SELECT * FROM tasks"); err == nil {
		return true
	}
	return false
}

func (t *taskDB) addTask(name, project string) error {
	_, err := t.db.Exec(
		"INSERT INTO tasks(name,project,status,created) VALUES(?,?,?,?)",
		name,
		project,
		todo.String(),
		time.Now(),
	)
	return err
}

func (t *taskDB) getTask(id uint) (task, error) {
	var task task
	err := t.db.QueryRow("SELECT * FROM tasks WHERE id = ?", id).Scan(
		&task.ID, &task.Name, &task.Project, &task.Status, &task.Created,
	)
	return task, err
}

func (t *taskDB) deleteTask(id uint) error {
	_, err := t.db.Exec("DELETE FROM tasks WHERE id = ?", id)
	return err
}

func (t *taskDB) updateTask(task task) error {
	orig, err := t.getTask(task.ID)
	if err != nil {
		return err
	}
	orig.merge(task)
	_, err = t.db.Exec(
		"UPDATE tasks SET name = ?, project = ?, status = ? WHERE id = ?",
		orig.Name,
		orig.Project,
		orig.Status,
		orig.ID,
	)
	return err
}

func (t *taskDB) getTasks() ([]task, error) {
	var tasks []task
	rows, err := t.db.Query("SELECT * FROM tasks")
	if err != nil {
		return tasks, fmt.Errorf("unable to get values; %w", err)
	}
	for rows.Next() {
		var task task
		err = rows.Scan(
			&task.ID,
			&task.Name,
			&task.Project,
			&task.Status,
			&task.Created,
		)
		if err != nil {
			return tasks, err
		}
		tasks = append(tasks, task)
	}
	return tasks, nil
}

func (t *taskDB) getTasksByStatus(status string) ([]task, error) {
	var tasks []task
	rows, err := t.db.Query(`SELECT * FROM tasks WHERE status = $1`, status)
	if err != nil {
		return tasks, fmt.Errorf("unable to get values: %w", err)
	}
	for rows.Next() {
		var task task
		err = rows.Scan(
			&task.ID,
			&task.Name,
			&task.Project,
			&task.Status,
			&task.Created,
		)
		if err != nil {
			return tasks, err
		}
		tasks = append(tasks, task)
	}
	return tasks, err
}

func (t *taskDB) getTasksByProject(project string) ([]task, error) {
	var tasks []task
	rows, err := t.db.Query(`SELECT * FROM tasks WHERE project LIKE ?`, "%"+project+"%")
	for rows.Next() {
		var task task
		err = rows.Scan(
			&task.ID,
			&task.Name,
			&task.Project,
			&task.Status,
			&task.Created,
		)
		if err != nil {
			return tasks, err
		}
		tasks = append(tasks, task)
	}
	return tasks, err
}

func initTaskDir(path string) error {
	if _, err := os.Stat(path); err != nil {
		if os.IsNotExist(err) {
			return os.Mkdir(path, 0o770)
		}
		return err
	}
	return nil
}

func initTasksDB(path string) (*taskDB, error) {
	db, err := sql.Open("sqlite3", filepath.Join(path, "tasks.db"))
	if err != nil {
		return nil, err
	}
	t := taskDB{db, path}
	if !t.taskTableExists() {
		if err := t.createTable(); err != nil {
			return nil, err
		}
	}
	return &t, nil
}
