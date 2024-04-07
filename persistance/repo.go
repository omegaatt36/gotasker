package persistance

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"slices"

	"github.com/omegaatt36/gotasker/domain"
	"github.com/omegaatt36/gotasker/persistance/models"

	"github.com/redis/go-redis/v9"
)

// RedisRepo represents a redis repository.
type RedisRepo struct {
	client *redis.Client
}

// NewRedisRepo creates a new redis repository.
func NewRedisRepo(client *redis.Client) *RedisRepo {
	return &RedisRepo{client: client}
}

// CreateTask creates a new task.
func (r *RedisRepo) CreateTask(ctx context.Context, req domain.CreateTaskRequest) error {
	id, err := r.client.Incr(ctx, models.KeyTaskAutoIncrementID).Result()
	if err != nil {
		return fmt.Errorf("failed to create task: %w", err)
	}

	modelTask := models.Task{
		ID:     uint(id),
		Name:   req.Name,
		Status: int(domain.TaskStatusIncomplete),
	}

	bs, err := json.Marshal(modelTask)
	if err != nil {
		return fmt.Errorf("failed to marshal task: %w", err)
	}

	if err := r.client.HSet(ctx, models.KeyTaskHMap, modelTask.Key(), string(bs)).Err(); err != nil {
		return fmt.Errorf("failed to create task: %w", err)
	}

	return nil
}

// ListTasks lists all tasks.
func (r *RedisRepo) ListTasks(ctx context.Context) ([]domain.Task, error) {
	tasks, err := r.client.HGetAll(ctx, models.KeyTaskHMap).Result()
	if err != nil {
		return nil, fmt.Errorf("failed to list tasks: %w", err)
	}

	modelTasks := make([]models.Task, 0, len(tasks))
	for _, task := range tasks {
		var modelTask models.Task
		if err := json.Unmarshal([]byte(task), &modelTask); err != nil {
			return nil, fmt.Errorf("failed to unmarshal task: %w", err)
		}
		modelTasks = append(modelTasks, modelTask)
	}

	slices.SortStableFunc(modelTasks, func(left, right models.Task) int {
		if left.ID < right.ID {
			return -1
		}

		return 1
	})

	result := make([]domain.Task, len(modelTasks))
	for index, t := range modelTasks {
		result[index] = domain.Task{
			ID:     t.ID,
			Name:   t.Name,
			Status: domain.TaskStatus(t.Status),
		}
	}

	return result, nil
}

// UpdateTask updates a task.
func (r *RedisRepo) UpdateTask(ctx context.Context, id uint, req domain.UpdateTaskRequest) error {

	modelTask := models.Task{
		ID: id,
	}

	bs, err := r.client.HGet(ctx, models.KeyTaskHMap, modelTask.Key()).Bytes()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return domain.ErrTaskNotFound
		}

		return fmt.Errorf("failed to get task: %w", err)
	}

	if err := json.Unmarshal(bs, &modelTask); err != nil {
		return fmt.Errorf("failed to unmarshal task: %w", err)
	}

	if req.Name != nil {
		modelTask.Name = *req.Name
	}

	if req.Status != nil {
		modelTask.Status = int(*req.Status)
	}

	bs, err = json.Marshal(modelTask)
	if err != nil {
		return fmt.Errorf("failed to marshal task: %w", err)
	}

	if err := r.client.HSet(ctx, models.KeyTaskHMap, modelTask.Key(), string(bs)).Err(); err != nil {
		return fmt.Errorf("failed to update task: %w", err)
	}

	return nil
}

// DeleteTask deletes a task.
func (r *RedisRepo) DeleteTask(ctx context.Context, id uint) error {

	modelTask := models.Task{
		ID: id,
	}

	_, err := r.client.HGet(ctx, models.KeyTaskHMap, modelTask.Key()).Bytes()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return domain.ErrTaskNotFound
		}

		return fmt.Errorf("failed to get task: %w", err)
	}

	if err := r.client.Del(ctx, modelTask.Key()).Err(); err != nil {
		// TODO: check if task not found
		if errors.Is(err, redis.Nil) {
			return domain.ErrTaskNotFound
		}
		return fmt.Errorf("failed to delete task: %w", err)
	}

	return nil
}
