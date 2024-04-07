//go:generate go-enum -f=$GOFILE

package domain

import (
	"context"
	"errors"
)

var (
	ErrTaskNotFound = errors.New("task not found")
)

// Task represents a task.
type Task struct {
	ID     uint
	Name   string
	Status TaskStatus
}

// TaskStatus represents a task status.
// ENUM(incomplete, completed)
type TaskStatus int

// TaskRepository represents a task repository.
type TaskRepository interface {
	CreateTask(ctx context.Context, req CreateTaskRequest) error
	ListTasks(ctx context.Context) ([]Task, error)
	UpdateTask(ctx context.Context, id uint, req UpdateTaskRequest) error
	DeleteTask(ctx context.Context, id uint) error
}

// CreateTaskRequest defines the request for creating a task.
type CreateTaskRequest struct {
	Name string
}

// UpdateTaskRequest defines the request for updating a task.
type UpdateTaskRequest struct {
	Name   *string
	Status *TaskStatus
}
