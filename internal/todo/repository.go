package todo

import (
	"gorm.io/gorm"
)

// Repository handles database operations for todos
type Repository interface {
	FindAll() ([]Todo, error)
	FindByID(id uint) (*Todo, error)
	Create(todo *Todo) error
	Update(todo *Todo) error
	Delete(id uint) error
}

type repository struct {
	db *gorm.DB
}

// NewRepository creates a new todo repository
func NewRepository(db *gorm.DB) Repository {
	return &repository{db: db}
}

func (r *repository) FindAll() ([]Todo, error) {
	var todos []Todo
	err := r.db.Find(&todos).Error
	return todos, err
}

func (r *repository) FindByID(id uint) (*Todo, error) {
	var todo Todo
	err := r.db.First(&todo, id).Error
	if err != nil {
		return nil, err
	}
	return &todo, nil
}

func (r *repository) Create(todo *Todo) error {
	return r.db.Create(todo).Error
}

func (r *repository) Update(todo *Todo) error {
	return r.db.Save(todo).Error
}

func (r *repository) Delete(id uint) error {
	result := r.db.Delete(&Todo{}, id)
	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return result.Error
}
