package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"
)

// main is the entry point for our command-line application.
func main() {
	manager := NewTodoManager()
	filename := "tasks.json"

	// Try to load tasks from file on startup
	if err := manager.LoadFromFile(filename); err != nil {
		fmt.Println("No existing tasks found. Starting with an empty list.")
	}

	args := os.Args[1:]
	if len(args) == 0 {
		fmt.Println("No command provided. Use 'add', 'list', 'complete', 'remove', 'save', or 'load'.")
		return
	}

	switch args[0] {
	case "add":
		if len(args) < 2 {
			fmt.Println("Usage: todo add <description> [due_date_YYYY-MM-DD]")
			return
		}
		description := strings.Join(args[1:], " ")
		dueDate := time.Now().Add(24 * time.Hour) // Default due date is tomorrow
		
		// Simple due date parsing
		parts := strings.Split(description, " ")
		if len(parts) > 1 {
			lastPart := parts[len(parts)-1]
			if t, err := time.Parse("2006-01-02", lastPart); err == nil {
				dueDate = t
				description = strings.Join(parts[:len(parts)-1], " ")
			}
		}

		manager.Add(description, dueDate)
		fmt.Println("Task added.")

	case "list":
		tasks := manager.List()
		fmt.Println(manager.FormattedList(tasks))

	case "incomplete":
		tasks := manager.Filter(func(task Task) bool {
			return !task.Completed
		})
		fmt.Println(manager.FormattedList(tasks))

	case "complete":
		if len(args) < 2 {
			fmt.Println("Usage: todo complete <id>")
			return
		}
		id, err := strconv.Atoi(args[1])
		if err != nil {
			fmt.Println("Invalid ID. Please provide a number.")
			return
		}
		if manager.Complete(id) {
			fmt.Printf("Task %d marked as complete.\n", id)
		} else {
			fmt.Printf("Task %d not found.\n", id)
		}

	case "remove":
		if len(args) < 2 {
			fmt.Println("Usage: todo remove <id>")
			return
		}
		id, err := strconv.Atoi(args[1])
		if err != nil {
			fmt.Println("Invalid ID. Please provide a number.")
			return
		}
		manager.Remove(id)
		fmt.Printf("Task %d removed.\n", id)
	case "save":
		if err := manager.SaveToFile(filename); err != nil {
			fmt.Printf("Error saving tasks: %v\n", err)
		} else {
			fmt.Println("Tasks saved successfully.")
		}

	case "load":
		if err := manager.LoadFromFile(filename); err != nil {
			fmt.Printf("Error loading tasks: %v\n", err)
		} else {
			fmt.Println("Tasks loaded successfully.")
		}
	
	default:
		fmt.Println("Unknown command. Use 'add', 'list', 'complete', 'remove', 'save', or 'load'.")
	}

	// Save tasks automatically on exit
	if err := manager.SaveToFile(filename); err != nil {
		fmt.Printf("Error autosaving tasks: %v\n", err)
	}
}
