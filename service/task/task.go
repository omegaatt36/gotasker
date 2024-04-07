package task

import (
	"context"
	"errors"

	"github.com/omegaatt36/gotasker/domain"
)

// Service represents a task service.
type Service struct {
	repo domain.TaskRepository
}

// NewService creates a new task service.
func NewService(repo domain.TaskRepository) *Service {
	return &Service{repo: repo}
}

// ListTasks lists all tasks.
func (s *Service) ListTasks(ctx context.Context) ([]domain.Task, error) {
	return s.repo.ListTasks(ctx)
}

// CreateTaskRequest defines the request for creating a task.
type CreateTaskRequest struct {
	Name string
}

// CreateTask creates a new task.
func (s *Service) CreateTask(ctx context.Context, req CreateTaskRequest) error {
	if req.Name == "" {
		return errors.New("task name is required")
	}

	return s.repo.CreateTask(ctx, domain.CreateTaskRequest{
		Name: req.Name,
	})
}

// UpdateTaskRequest defines the request for updating a task.
type UpdateTaskRequest struct {
	Name   *string
	Status *domain.TaskStatus
}

// UpdateTask updates a task.
func (s *Service) UpdateTask(ctx context.Context, id uint, req UpdateTaskRequest) error {
	if req.Status != nil && !req.Status.IsValid() {
		return errors.New("invalid status")
	}

	return s.repo.UpdateTask(ctx, id, domain.UpdateTaskRequest{
		Name:   req.Name,
		Status: req.Status,
	})
}

// DeleteTask deletes a task.
func (s *Service) DeleteTask(ctx context.Context, id uint) error {
	return s.repo.DeleteTask(ctx, id)
}
