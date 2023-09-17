package main

import (
	"errors"
	"fmt"
	"strconv"
	"time"

	"github.com/charmbracelet/bubbles/list"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "tasks",
	Short: "CLI task management tool for ~slaying~ your to do list",
	Args:  cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		return cmd.Help()
	},
}

var addCmd = &cobra.Command{
	Use:   "add",
	Short: "Add a new task with an optional project name",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		t, err := initTasksDB(setupPath())
		if err != nil {
			return err
		}
		defer t.db.Close()

		project, err := cmd.Flags().GetString("project")
		if err != nil {
			return err
		}
		if err := t.addTask(args[0], project); err != nil {
			return err
		}
		return nil
	},
}

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List all your tasks",
	Args:  cobra.NoArgs,
	PreRunE: func(cmd *cobra.Command, args []string) error {
		if cmd.Flags().Changed("status") && cmd.Flags().Changed("project") {
			return errors.New("Can only filter by project or status, not both.")
		}
		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		t, err := initTasksDB(setupPath())
		if err != nil {
			return err
		}
		defer t.db.Close()

		prog, err := cmd.Flags().GetInt("status")
		if err != nil {
			return err
		}

		var status string
		switch prog {
		case int(inProgress):
			status = inProgress.String()
		case int(done):
			status = done.String()
		case int(todo):
			status = todo.String()
		default:
			status = ""
		}

		project, err := cmd.Flags().GetString("project")
		if err != nil {
			return err
		}

		var tasks []task

		if status != "" {
			tasks, err = t.getTasksByStatus(status)
			if err != nil {
				return err
			}
		} else if project != "" {
			tasks, err = t.getTasksByProject(project)
			if err != nil {
				return err
			}
		} else {
			tasks, err = t.getTasks()
			if err != nil {
				return err
			}
		}
		if len(tasks) == 0 {
			fmt.Println("No tasks found")
		} else {
			table := setupTable(tasks)
			fmt.Println(table.View())
		}

		return nil
	},
}

var deleteCmd = &cobra.Command{
	Use:   "delete",
	Short: "Delete a task",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		t, err := initTasksDB(setupPath())
		if err != nil {
			return err
		}
		defer t.db.Close()

		id, err := strconv.Atoi(args[0])
		if err != nil {
			return err
		}
		return t.deleteTask(uint(id))
	},
}

var updateCmd = &cobra.Command{
	Use:   "update",
	Short: "Update a task",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		t, err := initTasksDB(setupPath())
		if err != nil {
			return err
		}
		defer t.db.Close()
		name, err := cmd.Flags().GetString("name")
		if err != nil {
			return err
		}
		project, err := cmd.Flags().GetString("project")
		if err != nil {
			return err
		}
		prog, err := cmd.Flags().GetInt("status")
		if err != nil {
			return err
		}
		id, err := strconv.Atoi(args[0])
		if err != nil {
			return err
		}
		var status string
		switch prog {
		case int(inProgress):
			status = inProgress.String()
		case int(done):
			status = done.String()
		default:
			status = todo.String()
		}

		updatedTask := task{uint(id), name, project, status, time.Time{}}
		return t.updateTask(updatedTask)

	},
}

func init() {
	addCmd.Flags().StringP("project", "p", "", "specify a project for your task")
	rootCmd.AddCommand(addCmd)
	listCmd.Flags().IntP("status", "s", -1, "specify a status for your task")
	listCmd.Flags().StringP("project", "p", "", "specify a project for your task")
	rootCmd.AddCommand(listCmd)
	rootCmd.AddCommand(deleteCmd)
	updateCmd.Flags().StringP("project", "p", "", "specify a project for your task")
	updateCmd.Flags().StringP("name", "n", "", "specify a name for your task")
	updateCmd.Flags().IntP("status", "s", int(todo), "specify a status for your task")
	rootCmd.AddCommand(updateCmd)
}
