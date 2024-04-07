package stub

import (
	"context"
	"errors"
	"slices"
	"sync"
	"time"

	"github.com/omegaatt36/gotasker/domain"
)

type task struct {
	ID        uint
	CreatedAt int64
	Name      string
	Status    domain.TaskStatus
}

// InMemoryTaskRepository is an stub implementation of in-memory task repository.
type InMemoryTaskRepository struct {
	sync.RWMutex

	taskAutoIncrementIDSequence uint

	tasks []task
}

// NewInMemoryTaskRepository creates a new in-memory task repository.
func NewInMemoryTaskRepository() *InMemoryTaskRepository {
	return &InMemoryTaskRepository{}
}

var _ domain.TaskRepository = (*InMemoryTaskRepository)(nil)

// CreateTask creates a new task.
func (repo *InMemoryTaskRepository) CreateTask(ctx context.Context, req domain.CreateTaskRequest) error {
	repo.Lock()
	defer repo.Unlock()

	repo.taskAutoIncrementIDSequence++
	repo.tasks = append(repo.tasks, task{
		ID:        repo.taskAutoIncrementIDSequence,
		CreatedAt: time.Now().Unix(),
		Name:      req.Name,
		Status:    domain.TaskStatusIncomplete,
	})

	return nil
}

// ListTasks lists all tasks.
func (repo *InMemoryTaskRepository) ListTasks(ctx context.Context) ([]domain.Task, error) {
	repo.RLock()
	defer repo.RUnlock()

	tasks := make([]task, len(repo.tasks))
	copy(tasks, repo.tasks)

	slices.SortStableFunc(tasks, func(left, right task) int {
		if left.CreatedAt > right.CreatedAt {
			return 0
		}

		return 1
	})

	result := make([]domain.Task, len(tasks))
	for index, t := range tasks {
		result[index] = domain.Task{
			ID:     t.ID,
			Name:   t.Name,
			Status: t.Status,
		}
	}

	return result, nil
}

// UpdateTask updates a task.
func (repo *InMemoryTaskRepository) UpdateTask(ctx context.Context, id uint, req domain.UpdateTaskRequest) error {
	repo.Lock()
	defer repo.Unlock()

	var indexOf *int
	for index, t := range repo.tasks {
		if t.ID == id {
			indexOf = &index
			break
		}
	}

	if indexOf == nil {
		return errors.New("task not found")
	}

	if req.Name != nil {
		repo.tasks[*indexOf].Name = *req.Name
	}
	if req.Status != nil {
		repo.tasks[*indexOf].Status = *req.Status
	}

	return nil
}

// DeleteTask deletes a task.
func (r *InMemoryTaskRepository) DeleteTask(ctx context.Context, id uint) error {
	r.Lock()
	defer r.Unlock()

	var indexOf *int
	for index, t := range r.tasks {
		if t.ID == id {
			indexOf = &index
			break
		}
	}

	if indexOf == nil {
		return errors.New("task not found")
	}

	r.tasks = append(r.tasks[:*indexOf], r.tasks[*indexOf+1:]...)

	return nil
}
