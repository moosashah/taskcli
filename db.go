package main

import (
	"database/sql"
	"fmt"
	"log"
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

func initTodosDB(path string) (*taskDB, error) {
	db, err := sql.Open("sqlite3", path)
	if err != nil {
		return nil, err
	}
	t := taskDB{db, path}
	dbErr := t.createTable()
	if dbErr != nil {
		log.Fatalf("could not initialize db: %+v", dbErr)
	}
	fmt.Printf("db setup: %t\n", t.taskTableExists())
	return &t, nil
}
