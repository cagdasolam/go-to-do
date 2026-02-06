package api

import (
	"example.com/mod/internal/todo"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

func SetupRoutes(router *gin.Engine, logger *zap.SugaredLogger, db *gorm.DB) {
	// Auto migrate the todo model
	if err := db.AutoMigrate(&todo.Todo{}); err != nil {
		logger.Fatalf("Failed to migrate database: %v", err)
	}

	// Initialize layers (Dependency Injection)
	// Repository -> Service -> Controller
	todoRepo := todo.NewRepository(db)
	todoService := todo.NewService(todoRepo)
	todoController := NewTodoController(todoService, logger)

	// API v1 routes
	v1 := router.Group("/api/v1")
	{
		todos := v1.Group("/todos")
		{
			todos.GET("", todoController.GetAllTodos)
			todos.GET("/:id", todoController.GetTodoByID)
			todos.POST("", todoController.CreateTodo)
			todos.PUT("/:id", todoController.UpdateTodo)
			todos.DELETE("/:id", todoController.DeleteTodo)
		}
	}
}
