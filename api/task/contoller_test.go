package task_test

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/omegaatt36/gotasker/api/task"
	"github.com/omegaatt36/gotasker/domain"
	"github.com/omegaatt36/gotasker/persistance"
	"github.com/omegaatt36/gotasker/persistance/database"
	taskService "github.com/omegaatt36/gotasker/service/task"
	"github.com/omegaatt36/gotasker/util"
	"github.com/stretchr/testify/suite"
)

type TaskControllerSuite struct {
	suite.Suite
}

func (s *TaskControllerSuite) SetupSuite() {
}

func (s *TaskControllerSuite) TestListTasks() {
	miniredis := database.InitializeTestingRedis()
	defer miniredis.Close()

	database.Initialize(context.Background(), miniredis.Addr(), "")

	repo := persistance.NewRedisRepo(database.Redis())
	service := taskService.NewService(repo)
	controller := task.NewController(service)

	type taskDetail struct {
		ID     uint   `json:"id"`
		Name   string `json:"name"`
		Status int    `json:"status"`
	}

	req := util.HTTPTestRequest{
		ServedURL:            "/tasks",
		RequestURLWithParams: "/tasks",
		Method:               http.MethodGet,
		HandleFuncs: []gin.HandlerFunc{
			controller.ListTasks,
		},
	}

	{
		resp, err := util.HTTPTest(req)

		s.NoError(err)
		s.Equal(http.StatusOK, resp.StatusCode)

		var tasks []taskDetail

		s.NoError(json.Unmarshal(resp.Body, &tasks))
		s.Len(tasks, 0)
	}

	// insert 10 tasks
	for index := range 10 {
		s.NoError(repo.CreateTask(context.Background(), domain.CreateTaskRequest{
			Name: fmt.Sprintf("task %d", index+1),
		}))
	}

	{
		resp, err := util.HTTPTest(req)

		s.NoError(err)
		s.Equal(http.StatusOK, resp.StatusCode)

		var tasks []taskDetail

		s.NoError(json.Unmarshal(resp.Body, &tasks))
		s.Len(tasks, 10)

		for index := range tasks {
			s.Equal(fmt.Sprintf("task %d", index+1), tasks[index].Name)
		}
	}
}

func (s *TaskControllerSuite) TestCreateTask() {
	miniredis := database.InitializeTestingRedis()
	defer miniredis.Close()

	database.Initialize(context.Background(), miniredis.Addr(), "")

	repo := persistance.NewRedisRepo(database.Redis())
	service := taskService.NewService(repo)
	controller := task.NewController(service)

	s.T().Run("without name", func(t *testing.T) {
		resp, err := util.HTTPTest(util.HTTPTestRequest{
			ServedURL:            "/tasks",
			RequestURLWithParams: "/tasks",
			Method:               http.MethodPost,
			HandleFuncs: []gin.HandlerFunc{
				controller.CreateTask,
			},
			Payload: map[string]any{},
		})
		s.NoError(err)
		s.Equal(http.StatusBadRequest, resp.StatusCode)
	})

	s.T().Run("success", func(t *testing.T) {
		resp, err := util.HTTPTest(util.HTTPTestRequest{
			ServedURL:            "/tasks",
			RequestURLWithParams: "/tasks",
			Method:               http.MethodPost,
			HandleFuncs: []gin.HandlerFunc{
				controller.CreateTask,
			},
			Payload: map[string]any{
				"name": "task 1",
			},
		})
		s.NoError(err)
		s.Equal(http.StatusCreated, resp.StatusCode)

		tasksInRepo, err := repo.ListTasks(context.Background())
		s.NoError(err)
		s.Len(tasksInRepo, 1)
		s.Equal("task 1", tasksInRepo[0].Name)
		s.Equal(domain.TaskStatusIncomplete, tasksInRepo[0].Status)
	})

	s.T().Run("success - duplicated name", func(t *testing.T) {
		resp, err := util.HTTPTest(util.HTTPTestRequest{
			ServedURL:            "/tasks",
			RequestURLWithParams: "/tasks",
			Method:               http.MethodPost,
			HandleFuncs: []gin.HandlerFunc{
				controller.CreateTask,
			},
			Payload: map[string]any{
				"name": "task 1",
			},
		})
		s.NoError(err)
		s.Equal(http.StatusCreated, resp.StatusCode)

		tasksInRepo, err := repo.ListTasks(context.Background())
		s.NoError(err)
		s.Len(tasksInRepo, 2)
		s.Equal("task 1", tasksInRepo[0].Name)
		s.Equal(domain.TaskStatusIncomplete, tasksInRepo[0].Status)
		s.Equal("task 1", tasksInRepo[1].Name)
		s.Equal(domain.TaskStatusIncomplete, tasksInRepo[1].Status)
	})
}

func (s *TaskControllerSuite) TestUpdateTask() {
	miniredis := database.InitializeTestingRedis()
	defer miniredis.Close()

	database.Initialize(context.Background(), miniredis.Addr(), "")

	repo := persistance.NewRedisRepo(database.Redis())
	service := taskService.NewService(repo)
	controller := task.NewController(service)

	s.T().Run("invalid id - not found", func(t *testing.T) {
		resp, err := util.HTTPTest(util.HTTPTestRequest{
			ServedURL:            "/tasks/:id",
			RequestURLWithParams: "/tasks/1",
			Method:               http.MethodPut,
			HandleFuncs: []gin.HandlerFunc{
				controller.UpdateTask,
			},
		})
		s.NoError(err)
		s.Equal(http.StatusNotFound, resp.StatusCode)
	})

	for index := range 10 {
		s.NoError(repo.CreateTask(context.Background(), domain.CreateTaskRequest{
			Name: fmt.Sprintf("task %d", index+1),
		}))
	}

	s.T().Run("invalid id - not numeric", func(t *testing.T) {
		resp, err := util.HTTPTest(util.HTTPTestRequest{
			ServedURL:            "/tasks/:id",
			RequestURLWithParams: "/tasks/a",
			Method:               http.MethodPut,
			HandleFuncs: []gin.HandlerFunc{
				controller.UpdateTask,
			},
		})
		s.NoError(err)
		s.Equal(http.StatusBadRequest, resp.StatusCode)
	})

	s.T().Run("invalid id - less than 1", func(t *testing.T) {
		resp, err := util.HTTPTest(util.HTTPTestRequest{
			ServedURL:            "/tasks/:id",
			RequestURLWithParams: "/tasks/0",
			Method:               http.MethodPut,
			HandleFuncs: []gin.HandlerFunc{
				controller.UpdateTask,
			},
		})
		s.NoError(err)
		s.Equal(http.StatusBadRequest, resp.StatusCode)
	})

	s.T().Run("invalid status", func(t *testing.T) {
		resp, err := util.HTTPTest(util.HTTPTestRequest{
			ServedURL:            "/tasks/:id",
			RequestURLWithParams: "/tasks/1",
			Method:               http.MethodPut,
			HandleFuncs: []gin.HandlerFunc{
				controller.UpdateTask,
			},
			Payload: map[string]any{
				"status": "invalid",
			},
		})
		s.NoError(err)
		s.Equal(http.StatusBadRequest, resp.StatusCode)
	})

	s.T().Run("success", func(t *testing.T) {
		resp, err := util.HTTPTest(util.HTTPTestRequest{
			ServedURL:            "/tasks/:id",
			RequestURLWithParams: "/tasks/1",
			Method:               http.MethodPut,
			HandleFuncs: []gin.HandlerFunc{
				controller.UpdateTask,
			},
			Payload: map[string]any{
				"name":   "task 1 - updated",
				"status": domain.TaskStatusCompleted,
			},
		})
		s.NoError(err)
		s.Equal(http.StatusOK, resp.StatusCode)

		tasksInRepo, err := repo.ListTasks(context.Background())
		s.NoError(err)
		s.Len(tasksInRepo, 10)
		s.Equal("task 1 - updated", tasksInRepo[0].Name)
		s.Equal(domain.TaskStatusCompleted, tasksInRepo[0].Status)
	})
}

func (s *TaskControllerSuite) TestDeleteTask() {
	miniredis := database.InitializeTestingRedis()
	defer miniredis.Close()

	database.Initialize(context.Background(), miniredis.Addr(), "")

	repo := persistance.NewRedisRepo(database.Redis())
	service := taskService.NewService(repo)
	controller := task.NewController(service)

	s.T().Run("invalid id - not found", func(t *testing.T) {
		resp, err := util.HTTPTest(util.HTTPTestRequest{
			ServedURL:            "/tasks/:id",
			RequestURLWithParams: "/tasks/1",
			Method:               http.MethodDelete,
			HandleFuncs: []gin.HandlerFunc{
				controller.DeleteTask,
			},
		})
		s.NoError(err)
		s.Equal(http.StatusNotFound, resp.StatusCode)
	})

	for index := range 10 {
		s.NoError(repo.CreateTask(context.Background(), domain.CreateTaskRequest{
			Name: fmt.Sprintf("task %d", index+1),
		}))
	}

	s.T().Run("invalid id - not numeric", func(t *testing.T) {
		resp, err := util.HTTPTest(util.HTTPTestRequest{
			ServedURL:            "/tasks/:id",
			RequestURLWithParams: "/tasks/a",
			Method:               http.MethodDelete,
			HandleFuncs: []gin.HandlerFunc{
				controller.DeleteTask,
			},
		})
		s.NoError(err)
		s.Equal(http.StatusBadRequest, resp.StatusCode)
	})

	s.T().Run("invalid id - less than 1", func(t *testing.T) {
		resp, err := util.HTTPTest(util.HTTPTestRequest{
			ServedURL:            "/tasks/:id",
			RequestURLWithParams: "/tasks/-999",
			Method:               http.MethodDelete,
			HandleFuncs: []gin.HandlerFunc{
				controller.DeleteTask,
			},
		})
		s.NoError(err)
		s.Equal(http.StatusBadRequest, resp.StatusCode)
	})

	s.T().Run("success", func(t *testing.T) {
		resp, err := util.HTTPTest(util.HTTPTestRequest{
			ServedURL:            "/tasks/:id",
			RequestURLWithParams: "/tasks/1",
			Method:               http.MethodDelete,
			HandleFuncs: []gin.HandlerFunc{
				controller.DeleteTask,
			},
		})
		s.NoError(err)
		s.Equal(http.StatusOK, resp.StatusCode)

		tasksInRepo, err := repo.ListTasks(context.Background())
		s.NoError(err)
		s.Len(tasksInRepo, 9)

		for index := 1; index <= 9; index++ {
			s.Equal(fmt.Sprintf("task %d", index+1), tasksInRepo[index-1].Name)
		}
	})
}

func TestTaskController(t *testing.T) {
	suite.Run(t, new(TaskControllerSuite))
}
