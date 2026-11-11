package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"text/tabwriter"
)

func storagePath() (string, error) {
	configDir, err := os.UserConfigDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(configDir, "cli-task-manager", "tasks.json"), nil
}

type Task struct {
	ID        int    `json:"id"`
	Title     string `json:"title"`
	Completed bool   `json:"completed"`
}

func list(tasks []Task) {
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 3, ' ', 0)
	fmt.Fprintln(w, "")
	fmt.Fprintln(w, "ID\tTITLE\tSTATUS")
	fmt.Fprintln(w, "--\t-----\t------")
	for _, task := range tasks {
		status := "[ ]"
		if task.Completed {
			status = "[*]"
		}
		fmt.Fprintf(w, "%d\t%s\t%s\n", task.ID, task.Title, status)
	}
	w.Flush()
}

func help() {
	fmt.Println(`
Commands list:
    list:               prints task list
    add [task]:         add a task
    complete [ID->int]: complete a task by its ID as an integer number
    delete [ID->int]:   delete a task
    help:               shows Commands list.`)
}

func add(tasks *[]Task, task Task) {
	*tasks = append(*tasks, task)
	fmt.Printf("The task is succesfully added to the list. Maybe u should go and do some of %d of them?\n", len(*tasks))
}
func complete(tasks *[]Task, id int) bool {
	for i := range *tasks {
		if (*tasks)[i].ID == id {
			(*tasks)[i].Completed = true
			fmt.Printf("%d\t%s\t%s\n", (*tasks)[i].ID, (*tasks)[i].Title, "[*]")
			return true
		}
	}
	return false
}

func delete(tasks []Task, id int) ([]Task, bool) {
	for i, task := range tasks {
		if task.ID == id {
			return append(tasks[:i], tasks[i+1:]...), true
		}
	}
	return tasks, false
}
func nextID(tasks []Task) int {
	maxID := 0
	for _, task := range tasks {
		if task.ID > maxID {
			maxID = task.ID
		}
	}
	return maxID + 1
}

func load() ([]Task, error) {
	path, err := storagePath()
	if err != nil {
		return nil, err
	}
	data, err := os.ReadFile(path)
	if errors.Is(err, os.ErrNotExist) {
		return []Task{}, nil
	}
	if err != nil {
		return nil, err
	}

	var tasks []Task
	if err := json.Unmarshal(data, &tasks); err != nil {
		return nil, err
	}
	return tasks, nil
}

func save(tasks []Task) error {
	data, err := json.MarshalIndent(tasks, "", "  ")
	if err != nil {
		return err
	}
	// TODO: write optional logs
	path, err := storagePath()
	if err != nil {
		return err
	}
	if err := os.MkdirAll(filepath.Dir(path), 0700); err != nil {
		return err
	}
	return os.WriteFile(path, data, 0600)
}

func main() {
	if len(os.Args) < 2 {
		help()
		return
	}

	tasks, err := load()
	if err != nil {
		fmt.Println("Failed to load tasks:", err)
		return
	}
	switch os.Args[1] {
	case "list":
		list(tasks)
	case "add":
		if len(os.Args) < 3 {
			help()
			return
		}
		if strings.TrimSpace(os.Args[2]) == "" {
			fmt.Println("Task title cannot be empty")
			return
		}
		add(&tasks, Task{
			ID:        nextID(tasks),
			Title:     os.Args[2],
			Completed: false,
		})
		if err := save(tasks); err != nil {
			fmt.Println("Failed to save tasks:", err)
		}
	case "complete":
		if len(os.Args) < 3 {
			help()
			return
		}
		//another error handling and I gonna kms
		id, err := strconv.Atoi(os.Args[2])
		if err != nil {
			fmt.Println("Invalid ID:", err)
			return
		}
		if !complete(&tasks, id) {
			fmt.Printf("Task %d not found\n", id)
			return
		}
		if err := save(tasks); err != nil {
			fmt.Println("Failed to save tasks:", err)
		}
	case "delete":
		if len(os.Args) < 3 {
			help()
			return
		}
		id, err := strconv.Atoi(os.Args[2])
		if err != nil {
			fmt.Println("Invalid ID:", err)
			return
		}
		var found bool
		tasks, found = delete(tasks, id)
		if !found {
			fmt.Printf("Task %d not found\n", id)
			return
		}
		if err := save(tasks); err != nil {
			fmt.Println("Failed to save tasks:", err)
		}
	default:
		help()
	}
}
