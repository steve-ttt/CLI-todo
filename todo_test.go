package main

import (
	"os"
	"testing"
	"time"
	"strings"
)

func TestAdd(t *testing.T) {
	// Create a new instance of our TodoManager.
	manager := NewTodoManager()

	// Define the task to add, including a due date.
	taskDesc := "Learn Go TDD"
	now := time.Now()
	
	// Call the Add method. We'll need to update it to accept a time.Time value.
	manager.Add(taskDesc, now)

	// Now, check if the task was added successfully.
	tasks := manager.List()

	// This is our first assertion. We expect exactly one task.
	if len(tasks) != 1 {
		t.Errorf("Expected 1 task after adding, but got %d", len(tasks))
	}

	// This is our second assertion. We expect the task description to match.
	if tasks[0].Description != taskDesc {
		t.Errorf("Expected task description to be '%s', but got '%s'", taskDesc, tasks[0].Description)
	}

	// This is our third assertion. We expect the due date to match.
	if tasks[0].DueDate != now {
		t.Errorf("Expected task due date to be '%s', but got '%s'", now, tasks[0].DueDate)
	}
}

func TestMarkComplete(t *testing.T) {
	// Create a new task manager and add a new task.
	manager := NewTodoManager()
	taskDesc := "Mark a task as complete"
	now := time.Now()
	manager.Add(taskDesc, now)

	// Now, mark the task as complete. We expect this method to be created next.
	manager.Complete(1)

	// Get the tasks and check the status of the first one.
	tasks := manager.List()
	task := tasks[0]

	// Now we can use the assert logic.
	if !task.Completed {
		t.Errorf("Expected task to be completed, but it was not")
	}

}

func TestFormattedList(t *testing.T) {
	manager := NewTodoManager()
	now := time.Now()

	manager.Add("Buy milk", now)
	manager.Add("Go for a run", now)
	manager.Complete(1)

	// We format the time now to make the test assertions easier.
	formattedTime := now.Format("2006-01-02")
	expected := "1: [x] Buy milk (Due: " + formattedTime + ")\n2: [ ] Go for a run (Due: " + formattedTime + ")"
	
	// We call the method with the result of manager.List() now.
	actual := manager.FormattedList(manager.List())

	if strings.TrimSpace(actual) != expected {
		t.Errorf("Expected formatted list:\n---\n%s\n---\nbut got:\n---\n%s\n---", expected, actual)
	}
}

// TestIncompleteTasks tests that the manager can filter for incomplete tasks.
func TestIncompleteTasks(t *testing.T) {
	manager := NewTodoManager()
	now := time.Now()

	manager.Add("Write blog post", now)
	manager.Add("Record podcast", now)
	manager.Add("Reply to emails", now)

	// Mark the second task as complete.
	manager.Complete(2)

	// Get the incomplete tasks. We expect the first and third tasks.
	incompleteTasks := manager.Filter(func(task Task) bool {
		return !task.Completed
	})

	if len(incompleteTasks) != 2 {
		t.Fatalf("Expected 2 incomplete tasks, but got %d", len(incompleteTasks))
	}

	// We check for the presence of the expected IDs, regardless of order.
	expectedIDs := map[int]bool{1: true, 3: true}
	foundIDs := make(map[int]bool, 2)
	for _, task := range incompleteTasks {
		foundIDs[task.Id] = true
	}

	if len(foundIDs) != len(expectedIDs) {
		t.Fatalf("Expected 2 unique IDs, but got %d", len(foundIDs))
	}

	for id := range expectedIDs {
		if !foundIDs[id] {
			t.Errorf("Expected ID %d not found in incomplete tasks", id)
		}
	}
}

// TestCompletedTasks tests that the manager can filter for completed tasks.
func TestCompletedTasks(t *testing.T) {
	manager := NewTodoManager()
	now := time.Now()

	manager.Add("Write blog post", now)
	manager.Add("Record podcast", now)
	manager.Add("Reply to emails", now)

	// Mark the second and third tasks as complete.
	manager.Complete(2)
	manager.Complete(3)

	// Get the completed tasks. We expect the second and third tasks.
	completedTasks := manager.Filter(func(task Task) bool {
		return task.Completed
	})

	if len(completedTasks) != 2 {
		t.Fatalf("Expected 2 completed tasks, but got %d", len(completedTasks))
	}

	// We check for the presence of the expected IDs, regardless of order.
	expectedIDs := map[int]bool{2: true, 3: true}
	foundIDs := make(map[int]bool, 2)
	for _, task := range completedTasks {
		foundIDs[task.Id] = true
	}

	for id := range expectedIDs {
		if !foundIDs[id] {
			t.Errorf("Expected ID %d not found in completed tasks", id)
		}
	}
}

// test we can delete a task
func TestRemove(t *testing.T) {
	// Create a new instance of our TodoManager.
	manager := NewTodoManager()

	// Define the task we want to add, including a due date.
	taskDesc := "Learn Go TDD"
	now := time.Now()
	
	// Call the Add method. We'll need to update it to accept a time.Time value.
	manager.Add(taskDesc, now)

	// Check if the task was added successfully.
	tasks := manager.List()
	
	// This is our first assertion. We expect exactly one task.
	if len(tasks) != 1 {
		t.Errorf("Expected 1 task after adding, but got %d", len(tasks))
	}

	// remove task
	manager.Remove(tasks[0].Id)
	
	// get the updated list
	updatedTasks := manager.List()

	// Assertion to test task has been removed
	if len(updatedTasks) != 0 {
		t.Errorf("Expected no tasks in list , but got %d", len(updatedTasks))
	}
	
	// Assertion to check if the removed task can be found
	// We'll use our Filter method to check if the task with the removed ID exists.
	foundTasks := manager.Filter(func(task Task) bool {
		return task.Id == tasks[0].Id
	})
	
	if len(foundTasks) != 0 {
		t.Errorf("Expected removed task to not be found in the list, but it was.")
	}
}

// TestSaveAndLoad tests that tasks can be saved to disk and loaded back.
func TestSaveAndLoad(t *testing.T) {
	manager := NewTodoManager()
	now := time.Now()

	// Add some tasks to our original manager.
	manager.Add("Buy vegan groceries", now)
	manager.Add("Schedule a meeting", now)
	manager.Complete(1)

	// Define a filename and defer its removal for a clean test environment.
	filename := "test_tasks.json"
	defer os.Remove(filename)
	
	// Call the new SaveToFile method.
	err := manager.SaveToFile(filename)
	if err != nil {
		t.Fatalf("Expected no error on save, but got: %v", err)
	}

	// Now, create a new manager and load the tasks from the file.
	loadedManager := NewTodoManager()
	err = loadedManager.LoadFromFile(filename)
	if err != nil {
		t.Fatalf("Expected no error on load, but got: %v", err)
	}

	// We compare the contents of both managers.
	originalTasks := manager.List()
	loadedTasks := loadedManager.List()
	
	if len(originalTasks) != len(loadedTasks) {
		t.Fatalf("Expected %d tasks after loading, but got %d", len(originalTasks), len(loadedTasks))
	}

	// We'll convert the lists to maps for easy comparison, ignoring order.
	originalMap := make(map[int]*Task, len(originalTasks))
	for _, task := range originalTasks {
		originalMap[task.Id] = &task
	}
	
	for _, loadedTask := range loadedTasks {
		originalTask, ok := originalMap[loadedTask.Id]
		if !ok {
			t.Errorf("Loaded task with ID %d not found in original manager", loadedTask.Id)
		}
		
		// The comparison of the tasks needs to be deep.
		// For a simple case, we can compare fields.
		if originalTask.Description != loadedTask.Description || originalTask.Completed != loadedTask.Completed {
			t.Errorf("Mismatch for task with ID %d: Expected %+v, Got %+v", originalTask.Id, originalTask, loadedTask)
		}
	}
}

