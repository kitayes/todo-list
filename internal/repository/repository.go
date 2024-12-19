package repository

import (
	"context"
	"github.com/jmoiron/sqlx"
	"todo/internal/models"
)

type Authorization interface {
	CreateUser(user models.User) (int, error)
	GetUser(username, password string) (models.User, error)
}

type TodoList interface {
	Create(userId int, list models.TodoList) (int, error)
	GetAll(userId int) ([]models.TodoList, error)
	GetById(userId, listId int) (models.TodoList, error)
	Delete(userId, listId int) error
	Update(userId, listId int, input models.UpdateListInput) error
}

type TodoItem interface {
	Create(listId int, item models.TodoItem) (int, error)
	GetAll(userId, listId int) ([]models.TodoItem, error)
	GetById(userId, itemId int) (models.TodoItem, error)
	Delete(userId, itemId int) error
	Update(userId, itemId int, input models.UpdateItemInput) error
}

type Repository struct {
	cfg *Config
	db  *sqlx.DB
	Authorization
	TodoList
	TodoItem
}

func NewRepository(cfg *Config) *Repository {
	return &Repository{
		cfg: cfg,
	}
}

func (r *Repository) Run(_ context.Context) error {
	return nil
}

func (r *Repository) Stop(_ context.Context) error {
	return r.db.Close()
}

func (r *Repository) Init() error {
	var err error
	r.db, err = newPostgresDB(r.cfg)
	if err != nil {
		return err
	}
	r.Authorization = NewAuthPostgres(r.db)
	r.TodoList = NewTodoListPostgres(r.db)
	r.TodoItem = NewTodoItemPostgres(r.db)

	return nil
}
