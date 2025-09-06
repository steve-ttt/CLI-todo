# Simple Go To-Do List Application

This is a simple command-line to-do list application built with Go. The project was developed using a Test-Driven Development (TDD) approach, focusing on writing tests for each piece of functionality before implementing the code.
Features

   * Add Tasks: Add a new task with a description and optional due date.

   * List Tasks: View a list of all tasks.

   * Complete Tasks: Mark a task as completed using its unique ID.

   * Remove Tasks: Delete a task by its unique ID.

   * Save & Load: Save tasks to a file on disk and load them on startup.

   * Filtering: Filter tasks to show only incomplete ones.

How to Run

  *  Ensure you have a Go environment set up.

   * Navigate to the project directory in your terminal.

   * Comple the application with 
   ```
   go build .
   ```

## Usage

Here are some examples of how to use the available commands:
Add a new task
```
$ ./todo "Learn more about TDD"
Task added.
```
## Add a task with a specific due date
```
$ ./todo add "Go grocery shopping" 2025-09-07
Task added.
```

List all tasks
```
$ ./todo list
1: [ ] Learn more about TDD (Due: 2025-09-07)
2: [ ] Go grocery shopping (Due: 2025-09-07)
```

Mark a task as complete
```
$ ./todo complete 1
Task 1 marked as complete.
```

View incomplete tasks
```
$ ./todo incomplete
2: [ ] Go grocery shopping (Due: 2025-09-07)
```

Remove a task
```
$ ./todo remove 2
Task 2 removed.
```

## Project Design

This application's design emerged organically through the TDD process. The core logic for managing tasks is encapsulated within the TodoManager struct in todo.go, while the command-line interface is handled separately in main.go.

