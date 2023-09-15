package main

import (
	"database/sql"
	"log"
	"os"
	"path/filepath"
	"reflect"
	"testing"

	_ "github.com/mattn/go-sqlite3"
)

func setup() *taskDB {
	path := filepath.Join(os.TempDir(), "test.db")
	db, err := sql.Open("sqlite3", path)
	if err != nil {
		log.Fatal(err)
	}

	t := taskDB{db, path}
	dbErr := t.createTable()
	if dbErr != nil {
		log.Fatal(dbErr)
	}
	return &t
}

func teardown(tDB *taskDB) {
	tDB.db.Close()
	os.Remove(tDB.dataDir)
}

func TestGetTask(t *testing.T) {
	tests := []struct {
		want task
	}{
		{
			want: task{
				ID:      1,
				Name:    "get milk",
				Project: "groceries",
				Status:  todo.String(),
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.want.Name, func(t *testing.T) {
			tDB := setup()
			defer teardown(tDB)

			if err := tDB.addTask(tc.want.Name, tc.want.Project); err != nil {
				t.Fatalf("could not add task, :%v", err)
			}
			task, err := tDB.getTask(tc.want.ID)
			if err != nil {
				t.Fatalf("could not get task of ID: %d. error: %v", tc.want.ID, err)
			}
			tc.want.Created = task.Created
			if !reflect.DeepEqual(task, tc.want) {
				t.Fatalf("got: %#v, want: %#v", task, tc.want)
			}
		})
	}
}

func TestDeleteTask(t *testing.T) {
	tests := []struct {
		want task
	}{
		{
			want: task{
				ID:      1,
				Name:    "get milk",
				Project: "groceries",
				Status:  todo.String(),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.want.Name, func(t *testing.T) {
			tDB := setup()
			defer teardown(tDB)
			if err := tDB.addTask(tt.want.Name, tt.want.Project); err != nil {
				t.Fatalf("Unable to insert task: %v", err)
			}
			task, err := tDB.getTask(tt.want.ID)
			if err != nil {
				t.Fatalf("unable to fetch task: %v", err)
			}
			tt.want.Created = task.Created
			if !reflect.DeepEqual(tt.want, task) {
				t.Fatalf("got %v, want %v", tt.want, task)
			}
			if err := tDB.deleteTask(tt.want.ID); err != nil {
				t.Fatalf("unable to delete task: %v", err)
			}
			if _, err := tDB.getTask(tt.want.ID); err == nil {
				t.Fatalf("task did not delete: %v", err)
			}
		})
	}
}

func TestUpdateTask(t *testing.T) {
	tests := []struct {
		old  *task
		new  *task
		want task
	}{
		{
			old: &task{
				ID:      1,
				Name:    "get milk",
				Project: "meep",
				Status:  todo.String(),
			},
			new: &task{
				ID:      1,
				Name:    "get pasta",
				Project: "groceries",
				Status:  todo.String(),
			},
			want: task{
				ID:      1,
				Name:    "get pasta",
				Project: "groceries",
				Status:  todo.String(),
			},
		},
		{
			old: &task{
				ID:      1,
				Name:    "get milk",
				Project: "groceries",
				Status:  todo.String(),
			},
			new: &task{
				ID:      1,
				Name:    "get milk",
				Project: "groceries",
				Status:  todo.String(), //update status
			},
			want: task{
				ID:      1,
				Name:    "get milk",
				Project: "groceries",
				Status:  todo.String(), //update status
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.want.Name, func(t *testing.T) {
			tDB := setup()
			defer teardown(tDB)
			if err := tDB.addTask(tt.old.Name, tt.old.Project); err != nil {
				t.Fatalf("could not create task: %v", err)
			}
			//update task
			if err := tDB.updateTask(*tt.new); err != nil {
				t.Fatalf("could not update task: %v", err)
			}
			task, err := tDB.getTask(tt.want.ID)
			if err != nil {
				t.Fatalf("could not get task: %v", err)
			}
			tt.want.Created = task.Created
			if !reflect.DeepEqual(task, tt.want) {
				t.Fatalf("got: %#v, want: %#v", task, tt.want)
			}
		})
	}
}

func TestGetTasks(t *testing.T) {
	tests := []struct {
		want []task
	}{
		{
			want: []task{
				{
					ID:      1,
					Name:    "get milk",
					Project: "groceries",
					Status:  todo.String(),
				},
				{
					ID:      2,
					Name:    "get milk",
					Project: "groceries",
					Status:  todo.String(),
				},
				{
					ID:      3,
					Name:    "get strawberries",
					Project: "groceries",
					Status:  todo.String(),
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run("get all tasks", func(t *testing.T) {
			tDB := setup()
			defer teardown(tDB)
			for _, task := range tt.want {
				if err := tDB.addTask(task.Name, task.Project); err != nil {
					t.Fatalf("error adding task: %v", err)
				}
			}
			got, err := tDB.getTasks()
			if err != nil {
				t.Fatalf("error getting tasks: %v", err)
			}
			for k, tc := range got {
				tt.want[k].Created = tc.Created
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Logf("got length: %v, want length: %v", len(got), len(tt.want))
				t.Errorf("taskDB.getTasks() = %v, want %v", got, tt.want)
			}
		})
	}
}
