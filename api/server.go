package api

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"github.com/omegaatt36/gotasker/api/task"
	"github.com/omegaatt36/gotasker/domain/stub"
	"github.com/omegaatt36/gotasker/logging"
	taskService "github.com/omegaatt36/gotasker/service/task"

	"github.com/gin-gonic/gin"
)

// Server is an api server
type Server struct {
	router *gin.Engine

	taskController *task.Controller
}

// NewServer creates a new server
func NewServer() *Server {
	apiEngine := gin.New()
	apiEngine.RedirectTrailingSlash = true

	// FIXME: use real data store
	taskController := task.NewController(taskService.NewService(
		stub.NewInMemoryTaskRepository(),
	))

	return &Server{
		router: apiEngine,

		taskController: taskController,
	}
}

// Start starts the server
func (s *Server) Start(ctx context.Context, appPort string) <-chan struct{} {
	s.router.Use(corsMiddleware())
	s.registerRoutes()

	srv := &http.Server{
		Addr:    fmt.Sprintf(":%s", appPort),
		Handler: s.router,
	}

	closeChain := make(chan struct{})
	go func() {
		defer func() {
			logging.Info("api stopped")
			closeChain <- struct{}{}
			close(closeChain)
		}()

		<-ctx.Done()
		if err := srv.Shutdown(ctx); err != nil {
			logging.Fatal("Server Shutdown: ", err)
		}
	}()

	logging.Info("starts serving...")

	go func() {
		if err := srv.ListenAndServe(); err != nil &&
			!errors.Is(err, http.ErrServerClosed) {
			logging.Fatalf("listen: %s\n", err)
		}
	}()

	return closeChain
}

func (s *Server) registerRoutes() {
	groupedRouter := s.router.Group("")

	groupedRouter.Use(injectLogging([]string{}), recovery())

	groupFilmLog := groupedRouter.Group("/tasks")
	groupFilmLog.GET("", s.taskController.ListTasks)
	groupFilmLog.POST("", s.taskController.CreateTask)
	groupFilmLog.PUT("/:id", s.taskController.UpdateTask)
	groupFilmLog.DELETE("/:id", s.taskController.DeleteTask)
}
