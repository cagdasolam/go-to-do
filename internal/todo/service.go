package todo

import (
	"errors"

	"gorm.io/gorm"
)

var (
	ErrTodoNotFound = errors.New("todo not found")
)

// Service handles business logic for todos
type Service interface {
	GetAll() ([]TodoResponse, error)
	GetByID(id uint) (*TodoResponse, error)
	Create(req CreateTodoRequest) (*TodoResponse, error)
	Update(id uint, req UpdateTodoRequest) (*TodoResponse, error)
	Delete(id uint) error
}

type service struct {
	repo Repository
}

// NewService creates a new todo service
func NewService(repo Repository) Service {
	return &service{repo: repo}
}

func (s *service) GetAll() ([]TodoResponse, error) {
	todos, err := s.repo.FindAll()
	if err != nil {
		return nil, err
	}

	response := make([]TodoResponse, len(todos))
	for i, t := range todos {
		response[i] = t.ToResponse()
	}
	return response, nil
}

func (s *service) GetByID(id uint) (*TodoResponse, error) {
	todo, err := s.repo.FindByID(id)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, ErrTodoNotFound
		}
		return nil, err
	}

	response := todo.ToResponse()
	return &response, nil
}

func (s *service) Create(req CreateTodoRequest) (*TodoResponse, error) {
	todo := &Todo{
		Title:       req.Title,
		Description: req.Description,
		Priority:    req.Priority,
		DueDate:     req.DueDate,
	}

	// Set default priority
	if todo.Priority == "" {
		todo.Priority = "medium"
	}

	if err := s.repo.Create(todo); err != nil {
		return nil, err
	}

	response := todo.ToResponse()
	return &response, nil
}

func (s *service) Update(id uint, req UpdateTodoRequest) (*TodoResponse, error) {
	todo, err := s.repo.FindByID(id)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, ErrTodoNotFound
		}
		return nil, err
	}

	// Update only provided fields
	if req.Title != nil {
		todo.Title = *req.Title
	}
	if req.Description != nil {
		todo.Description = *req.Description
	}
	if req.Completed != nil {
		todo.Completed = *req.Completed
	}
	if req.Priority != nil {
		todo.Priority = *req.Priority
	}
	if req.DueDate != nil {
		todo.DueDate = req.DueDate
	}

	if err := s.repo.Update(todo); err != nil {
		return nil, err
	}

	response := todo.ToResponse()
	return &response, nil
}

func (s *service) Delete(id uint) error {
	err := s.repo.Delete(id)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return ErrTodoNotFound
		}
		return err
	}
	return nil
}
