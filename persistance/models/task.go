package models

import "fmt"

// task related constants
const (
	KeyTaskAutoIncrementID = "tasks_auto_increment_id"
	KeyTaskHMap            = "tasks_map"
)

// Task represents a task.
type Task struct {
	ID     uint   `json:"id"`
	Name   string `json:"name"`
	Status int    `json:"status"`
}

// Key returns key.
func (t *Task) Key() string {
	return fmt.Sprintf("%d", t.ID)
}
