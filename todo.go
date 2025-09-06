package main

import (
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"
)

// Task is a struct that represents a single task in our todo list.
type Task struct {
	Id          int
	Description string
	Completed   bool
	DueDate     time.Time
}

// TodoManager is a struct that manages our tasks.
type TodoManager struct {
	Count int
	Tasks map[int]*Task
}

// NewTodoManager returns a new instance of TodoManager.
func NewTodoManager() *TodoManager {
	return &TodoManager{
		Tasks: make(map[int]*Task),
	}
}

// Add adds a new task to the manager.
func (tm *TodoManager) Add(description string, due time.Time) {
	tm.Count += 1
	tm.Tasks[tm.Count] = &Task{
		Id:          tm.Count,
		Description: description,
		Completed:   false,
		DueDate:     due,
	}
}

// mark a task complete given its ID.
func (tm *TodoManager) Complete(id int) bool {
	task, ok := tm.Tasks[id] // check task ID exists in the map
	if ok {
		task.Completed = true
		return true
	}
	return false
}

// List returns a slice of all tasks.
func (tm *TodoManager) List() []Task {
	// Create a new slice to hold the tasks.
	var tasks []Task
	for _, task := range tm.Tasks {
		// Dereference the pointer to get the Task value.
		tasks = append(tasks, *task)
	}
	return tasks
}

// returns a formated list of all tasks
func (tm *TodoManager) FormattedList(taskList []Task) string {
	builder := strings.Builder{}
	for _, task := range taskList {

		builder.WriteString(strconv.Itoa(task.Id))
		builder.WriteString(": ")
		if task.Completed == true {
			builder.WriteString("[x]")
		} else {
			builder.WriteString("[ ]")
		}
		builder.WriteString(" ")
		builder.WriteString(task.Description)
		builder.WriteString(" (Due: ")
		builder.WriteString(task.DueDate.Format("2006-01-02"))
		builder.WriteString(")\n")
	}

	return builder.String()
}

func (tm *TodoManager) Filter(fn func(task Task) bool) []Task {
	var filteredTasks []Task
	tasks := tm.List()
	for _, task := range tasks {
		if fn(task) {
			filteredTasks = append(filteredTasks, task)
		}
	}
	return filteredTasks

}

func (tm *TodoManager) Remove(index int) {
	delete(tm.Tasks, index)
}

// SaveToFile saves the tasks to a JSON file on disk.
func (tm *TodoManager) SaveToFile(filename string) error {
	jsonData, err := json.MarshalIndent(tm.List(), "", "  ")
	if err != nil {
		return fmt.Errorf("error marshalling to JSON: %w", err)
	}
	err = os.WriteFile(filename, jsonData, 0644)
	if err != nil {
		return fmt.Errorf("error writing to file: %w", err)
	}
	return nil
}

// LoadFromFile loads tasks from a JSON file on disk.
func (tm *TodoManager) LoadFromFile(filename string) error {
	data, err := os.ReadFile(filename)
	if err != nil {
		return fmt.Errorf("error reading file: %w", err)
	}

	var tasks []Task
	err = json.Unmarshal(data, &tasks)
	if err != nil {
		return fmt.Errorf("error unmarshalling JSON: %w", err)
	}

	for i := range tasks {
		tm.Tasks[tasks[i].Id] = &tasks[i]
		if tasks[i].Id > tm.Count {
			tm.Count = tasks[i].Id
		}
	}

	return nil
}
