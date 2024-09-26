package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"strconv"
	"time"
)

type Task struct {
	ID        int       `json:"id"`
	Title     string    `json:"title"`
	CreatedAt time.Time `json:"CreatedAt"`
	Status    bool      `json:"status"`
}

func ReadTasks(filename string) ([]Task, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, fmt.Errorf("error opening file: %w", err)
	}
	defer file.Close()

	data, err := os.ReadFile(filename)
	if err != nil {
		return nil, fmt.Errorf("error reading file: %w", err)
	}
	var tasks []Task
	err = json.Unmarshal(data, &tasks)
	if err != nil {
		return nil, fmt.Errorf("error deserializing JSON: %w", err)
	}

	return tasks, nil

}

func writeTasks(filename string, tasks []Task) error {
	file, err := os.Create(filename)
	if err != nil {
		return fmt.Errorf("error creating file: %w", err)
	}
	defer file.Close()

	data, err := json.MarshalIndent(tasks, "", "  ")
	if err != nil {
		return fmt.Errorf("error encoding JSON: %w", err)
	}

	_, err = file.Write(data)
	if err != nil {
		return fmt.Errorf("error writing to file: %w", err)
	}

	return nil
}

func AddTasks(name string, filename string) error {
	tasks, err := ReadTasks(filename)
	if err != nil {
		return err
	}

	newTask := Task{
		ID:        len(tasks) + 5381,
		Title:     name,
		CreatedAt: time.Now(),
		Status:    false,
	}

	tasks = append(tasks, newTask)
	return writeTasks(filename, tasks)
}

func ListTasks(filename string) error {
	tasks, err := ReadTasks(filename)
	if err != nil {
		return err
	}

	if len(tasks) == 0 {
		fmt.Println("No tasks found.")
		return nil
	}

	fmt.Println("Tasks:")
	for _, task := range tasks {
		status := "Not completed"
		if task.Status {
			status = "Completed"
		}

		fmt.Printf("%d. %s [%s]\n", task.ID, task.Title, status)
	}

	return nil
}

func CompleteTask(id int, filename string) error {
	tasks, err := ReadTasks(filename)
	if err != nil {
		return err
	}

	for i, task := range tasks {
		if task.ID == id {
			tasks[i].Status = true
			fmt.Printf("Task %d marked as completed.\n", id)
			return writeTasks(filename, tasks)
		}
	}

	return fmt.Errorf("task with ID %d not found", id)
}

func main() {

	const filename = "tasks.json"

	addCmd := flag.NewFlagSet("add", flag.ExitOnError)
	Completecmd := flag.NewFlagSet("complete", flag.ExitOnError)

	AddTaskName := addCmd.String("name", "", "Taskname")

	if len(os.Args) < 2 {
		fmt.Println("expected 'add', 'list', or 'complete' subcommands")
		os.Exit(1)
	}

	switch os.Args[1] {
	case "add":
		addCmd.Parse(os.Args[2:])
		if *AddTaskName == "" {
			fmt.Println("Task name cannot be empty")
			os.Exit(1)
		}
		err := AddTasks(*AddTaskName, filename)
		if err != nil {
			fmt.Println("Error adding task:", err)
		} else {
			fmt.Println("Task added successfully!")
		}
	case "list":
		err := ListTasks(filename)
		if err != nil {
			fmt.Println("Error listing tasks:", err)
		}
	case "complete":
		Completecmd.Parse(os.Args[2:])
		if Completecmd.Arg(0) == "" {
			fmt.Println("Please provide the task ID to complete")
			os.Exit(1)
		}
		id, err := strconv.Atoi(Completecmd.Arg(0))
		if err != nil {
			fmt.Println("Invalid task ID")
			os.Exit(1)
		}
		err = CompleteTask(id, filename)
		if err != nil {
			fmt.Println("Error completing task:", err)
		}
	default:
		fmt.Println("expected 'add', 'list', or 'complete' subcommands")
		os.Exit(1)
	}

}
