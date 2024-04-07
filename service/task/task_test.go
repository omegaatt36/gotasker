package task_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/omegaatt36/gotasker/domain"
	"github.com/omegaatt36/gotasker/domain/stub"
	"github.com/omegaatt36/gotasker/service/task"
	"github.com/omegaatt36/gotasker/util"

	"github.com/stretchr/testify/suite"
)

type TaskServiceTaskSuite struct {
	suite.Suite
}

func (s *TaskServiceTaskSuite) SetupSuite() {
}

func (s *TaskServiceTaskSuite) TestCreateTask() {
	repo := stub.NewInMemoryTaskRepository()
	service := task.NewService(repo)

	tasksInRepo, err := repo.ListTasks(context.Background())
	s.NoError(err)
	s.Empty(tasksInRepo)

	s.T().Run("without name", func(t *testing.T) {
		s.Error(service.CreateTask(context.Background(), task.CreateTaskRequest{
			Name: "",
		}), "task name is required")
	})
	s.T().Run("success", func(t *testing.T) {
		s.NoError(service.CreateTask(context.Background(), task.CreateTaskRequest{
			Name: "task 1",
		}))
	})

	tasksInRepo, err = repo.ListTasks(context.Background())
	s.NoError(err)
	s.Len(tasksInRepo, 1)

	s.Equal("task 1", tasksInRepo[0].Name)
	s.Equal(domain.TaskStatusIncomplete, tasksInRepo[0].Status)

	s.T().Run("another task", func(t *testing.T) {
		s.NoError(service.CreateTask(context.Background(), task.CreateTaskRequest{
			Name: "task 2",
		}))
	})

	tasksInRepo, err = repo.ListTasks(context.Background())
	s.NoError(err)
	s.Len(tasksInRepo, 2)

	s.Equal("task 1", tasksInRepo[0].Name)
	s.Equal(domain.TaskStatusIncomplete, tasksInRepo[0].Status)
	s.Equal("task 2", tasksInRepo[1].Name)
	s.Equal(domain.TaskStatusIncomplete, tasksInRepo[1].Status)

	s.T().Run("duplicated task name is allowed", func(t *testing.T) {
		s.NoError(service.CreateTask(context.Background(), task.CreateTaskRequest{
			Name: "task 1",
		}))
	})

	tasksInRepo, err = repo.ListTasks(context.Background())
	s.NoError(err)
	s.Len(tasksInRepo, 3)

	s.Equal("task 1", tasksInRepo[0].Name)
	s.Equal(domain.TaskStatusIncomplete, tasksInRepo[0].Status)
	s.Equal("task 2", tasksInRepo[1].Name)
	s.Equal(domain.TaskStatusIncomplete, tasksInRepo[1].Status)
	s.Equal("task 1", tasksInRepo[2].Name)
	s.Equal(domain.TaskStatusIncomplete, tasksInRepo[2].Status)
}

func (s *TaskServiceTaskSuite) TestListTasks() {
	repo := stub.NewInMemoryTaskRepository()
	service := task.NewService(repo)

	// 1.22 new feature: range a number like other language :D
	for index := range 10 {
		repo.CreateTask(context.Background(), domain.CreateTaskRequest{
			Name: fmt.Sprintf("task %d", index+1),
		})
	}

	tasksInRepo, err := service.ListTasks(context.Background())
	s.NoError(err)
	s.Len(tasksInRepo, 10)

	// assert order by created_at(id)
	for index := range 10 {
		s.Equal(fmt.Sprintf("task %d", index+1), tasksInRepo[index].Name)
		s.Equal(domain.TaskStatusIncomplete, tasksInRepo[index].Status)
	}
}

func (s *TaskServiceTaskSuite) TestUpdateTask() {
	repo := stub.NewInMemoryTaskRepository()
	service := task.NewService(repo)

	s.NoError(repo.CreateTask(context.Background(), domain.CreateTaskRequest{
		Name: "task 1",
	}))

	tasksInRepo, err := repo.ListTasks(context.Background())
	s.NoError(err)
	s.Len(tasksInRepo, 1)
	s.Equal("task 1", tasksInRepo[0].Name)
	s.Equal(domain.TaskStatusIncomplete, tasksInRepo[0].Status)

	s.T().Run("invalid status", func(t *testing.T) {
		s.Error(service.UpdateTask(context.Background(), 1, task.UpdateTaskRequest{
			Status: util.Pointer(domain.TaskStatus(99999)),
		}))
	})

	{ // update only status
		s.T().Run("update status success", func(t *testing.T) {
			s.NoError(service.UpdateTask(context.Background(), 1, task.UpdateTaskRequest{
				Status: util.Pointer(domain.TaskStatusCompleted),
			}))
		})

		tasksInRepo, err = repo.ListTasks(context.Background())
		s.NoError(err)
		s.Len(tasksInRepo, 1)
		s.Equal("task 1", tasksInRepo[0].Name)
		s.Equal(domain.TaskStatusCompleted, tasksInRepo[0].Status)
	}
	{ // update only name
		s.T().Run("update status success", func(t *testing.T) {
			s.NoError(service.UpdateTask(context.Background(), 1, task.UpdateTaskRequest{
				Name: util.Pointer("task 1 updated"),
			}))
		})

		tasksInRepo, err = repo.ListTasks(context.Background())
		s.NoError(err)
		s.Len(tasksInRepo, 1)
		s.Equal("task 1 updated", tasksInRepo[0].Name)
		s.Equal(domain.TaskStatusCompleted, tasksInRepo[0].Status)
	}
	{ // update both name and status
		s.T().Run("update status success", func(t *testing.T) {
			s.NoError(service.UpdateTask(context.Background(), 1, task.UpdateTaskRequest{
				Status: util.Pointer(domain.TaskStatusIncomplete),
				Name:   util.Pointer("task 1 updated 2"),
			}))
		})

		tasksInRepo, err = repo.ListTasks(context.Background())
		s.NoError(err)
		s.Len(tasksInRepo, 1)
		s.Equal("task 1 updated 2", tasksInRepo[0].Name)
		s.Equal(domain.TaskStatusIncomplete, tasksInRepo[0].Status)
	}

	s.T().Run("task not found", func(t *testing.T) {
		s.Error(service.UpdateTask(context.Background(), 2, task.UpdateTaskRequest{
			Status: util.Pointer(domain.TaskStatusCompleted),
		}))
	})
}

func (s *TaskServiceTaskSuite) TestDeleteTask() {
	repo := stub.NewInMemoryTaskRepository()
	service := task.NewService(repo)

	for index := range 10 {
		repo.CreateTask(context.Background(), domain.CreateTaskRequest{
			Name: fmt.Sprintf("task %d", index+1),
		})
	}

	tasksInRepo, err := repo.ListTasks(context.Background())
	s.NoError(err)
	s.Len(tasksInRepo, 10)

	s.NoError(service.DeleteTask(context.Background(), 1))
	s.NoError(service.DeleteTask(context.Background(), 2))
	s.NoError(service.DeleteTask(context.Background(), 3))

	tasksInRepo, err = repo.ListTasks(context.Background())
	s.NoError(err)
	s.Len(tasksInRepo, 7)

	s.Equal("task 4", tasksInRepo[0].Name)
	s.Equal("task 5", tasksInRepo[1].Name)
	s.Equal("task 6", tasksInRepo[2].Name)
	s.Equal("task 7", tasksInRepo[3].Name)
	s.Equal("task 8", tasksInRepo[4].Name)
	s.Equal("task 9", tasksInRepo[5].Name)
	s.Equal("task 10", tasksInRepo[6].Name)
}

func TestTaskService(t *testing.T) {
	suite.Run(t, new(TaskServiceTaskSuite))
}
