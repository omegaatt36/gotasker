package task

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/omegaatt36/gotasker/domain"
	"github.com/omegaatt36/gotasker/service/task"

	"github.com/gin-gonic/gin"
)

// Controller represents a task controller.
type Controller struct {
	service *task.Service
}

// NewController creates a new task controller.
func NewController(service *task.Service) *Controller {
	return &Controller{service: service}
}

// taskDetail defines DTO for domain.Task.
type taskDetail struct {
	ID     uint   `json:"id"`
	Name   string `json:"name"`
	Status int    `json:"status"`
}

func (task *taskDetail) fromDomain(domainTask *domain.Task) {
	task.ID = domainTask.ID
	task.Name = domainTask.Name
	task.Status = int(domainTask.Status)
}

// ListTasks lists all tasks.
func (x *Controller) ListTasks(c *gin.Context) {
	tasks, err := x.service.ListTasks(c.Request.Context())
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	taskDetails := make([]*taskDetail, len(tasks))
	for index := range tasks {
		taskDetails[index] = &taskDetail{}
		taskDetails[index].fromDomain(&tasks[index])
	}

	c.JSON(http.StatusOK, taskDetails)
}

// createTaskRequest defines the request for creating a task.
type createTaskRequest struct {
	Name string `json:"name" binding:"required"`
}

// CreateTask creates a new task.
func (x *Controller) CreateTask(c *gin.Context) {
	var req createTaskRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	if err := x.service.CreateTask(c.Request.Context(), task.CreateTaskRequest{
		Name: req.Name,
	}); err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	c.Status(http.StatusCreated)
}

// UpdateTaskRequest defines the request for updating a task.
type updateTaskRequest struct {
	Name   *string `json:"name"`
	Status *int    `json:"status"`
}

// UpdateTask updates a task.
func (x *Controller) UpdateTask(c *gin.Context) {
	taskID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	if taskID < 1 {
		c.AbortWithError(http.StatusBadRequest, errors.New("invalid task id"))
		return
	}

	var req updateTaskRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	var status *domain.TaskStatus
	if req.Status != nil {
		domainTaskStatus := domain.TaskStatus(*req.Status)
		if !domainTaskStatus.IsValid() {
			c.AbortWithError(http.StatusBadRequest, domain.ErrInvalidTaskStatus)
			return
		}

		status = &domainTaskStatus
	}

	if err := x.service.UpdateTask(c.Request.Context(), uint(taskID), task.UpdateTaskRequest{
		Name:   req.Name,
		Status: status,
	}); err != nil {
		if errors.Is(err, domain.ErrTaskNotFound) {
			c.AbortWithError(http.StatusNotFound, err)
			return
		}

		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	c.Status(http.StatusOK)
}

// DeleteTask deletes a task.
func (x *Controller) DeleteTask(c *gin.Context) {
	taskID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	if taskID < 1 {
		c.AbortWithError(http.StatusBadRequest, errors.New("invalid task id"))
		return
	}

	if err := x.service.DeleteTask(c.Request.Context(), uint(taskID)); err != nil {
		if errors.Is(err, domain.ErrTaskNotFound) {
			c.AbortWithError(http.StatusNotFound, err)
			return
		}

		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	c.Status(http.StatusOK)
}