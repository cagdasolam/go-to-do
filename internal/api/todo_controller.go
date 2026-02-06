package api

import (
	"net/http"
	"strconv"

	"example.com/mod/internal/todo"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// TodoController handles HTTP requests for todos
// Controller only handles: request parsing, validation, calling service, response formatting
type TodoController struct {
	service todo.Service
	logger  *zap.SugaredLogger
}

// NewTodoController creates a new todo controller
func NewTodoController(service todo.Service, logger *zap.SugaredLogger) *TodoController {
	return &TodoController{
		service: service,
		logger:  logger,
	}
}

// GetAllTodos godoc
// @Summary      List all todos
// @Description  Get all todos from the database
// @Tags         todos
// @Accept       json
// @Produce      json
// @Success      200  {array}   todo.TodoResponse
// @Failure      500  {object}  map[string]string
// @Router       /todos [get]
func (tc *TodoController) GetAllTodos(c *gin.Context) {
	todos, err := tc.service.GetAll()
	if err != nil {
		tc.logger.Errorw("Failed to fetch todos", "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch todos"})
		return
	}

	c.JSON(http.StatusOK, todos)
}

// GetTodoByID godoc
// @Summary      Get a todo by ID
// @Description  Get a single todo by its ID
// @Tags         todos
// @Accept       json
// @Produce      json
// @Param        id   path      int  true  "Todo ID"
// @Success      200  {object}  todo.TodoResponse
// @Failure      400  {object}  map[string]string
// @Failure      404  {object}  map[string]string
// @Router       /todos/{id} [get]
func (tc *TodoController) GetTodoByID(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	result, err := tc.service.GetByID(uint(id))
	if err != nil {
		if err == todo.ErrTodoNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "Todo not found"})
			return
		}
		tc.logger.Errorw("Failed to fetch todo", "error", err, "id", id)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch todo"})
		return
	}

	c.JSON(http.StatusOK, result)
}

// CreateTodo godoc
// @Summary      Create a new todo
// @Description  Create a new todo item
// @Tags         todos
// @Accept       json
// @Produce      json
// @Param        todo  body      todo.CreateTodoRequest  true  "Todo to create"
// @Success      201   {object}  todo.TodoResponse
// @Failure      400   {object}  map[string]string
// @Failure      500   {object}  map[string]string
// @Router       /todos [post]
func (tc *TodoController) CreateTodo(c *gin.Context) {
	var req todo.CreateTodoRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	result, err := tc.service.Create(req)
	if err != nil {
		tc.logger.Errorw("Failed to create todo", "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create todo"})
		return
	}

	tc.logger.Infow("Todo created", "id", result.ID)
	c.JSON(http.StatusCreated, result)
}

// UpdateTodo godoc
// @Summary      Update a todo
// @Description  Update an existing todo by ID
// @Tags         todos
// @Accept       json
// @Produce      json
// @Param        id    path      int                     true  "Todo ID"
// @Param        todo  body      todo.UpdateTodoRequest  true  "Todo updates"
// @Success      200   {object}  todo.TodoResponse
// @Failure      400   {object}  map[string]string
// @Failure      404   {object}  map[string]string
// @Failure      500   {object}  map[string]string
// @Router       /todos/{id} [put]
func (tc *TodoController) UpdateTodo(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	var req todo.UpdateTodoRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	result, err := tc.service.Update(uint(id), req)
	if err != nil {
		if err == todo.ErrTodoNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "Todo not found"})
			return
		}
		tc.logger.Errorw("Failed to update todo", "error", err, "id", id)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update todo"})
		return
	}

	tc.logger.Infow("Todo updated", "id", result.ID)
	c.JSON(http.StatusOK, result)
}

// DeleteTodo godoc
// @Summary      Delete a todo
// @Description  Delete a todo by ID (soft delete)
// @Tags         todos
// @Accept       json
// @Produce      json
// @Param        id   path      int  true  "Todo ID"
// @Success      200  {object}  map[string]string
// @Failure      400  {object}  map[string]string
// @Failure      404  {object}  map[string]string
// @Failure      500  {object}  map[string]string
// @Router       /todos/{id} [delete]
func (tc *TodoController) DeleteTodo(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	err = tc.service.Delete(uint(id))
	if err != nil {
		if err == todo.ErrTodoNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "Todo not found"})
			return
		}
		tc.logger.Errorw("Failed to delete todo", "error", err, "id", id)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete todo"})
		return
	}

	tc.logger.Infow("Todo deleted", "id", id)
	c.JSON(http.StatusOK, gin.H{"message": "Todo deleted successfully"})
}
