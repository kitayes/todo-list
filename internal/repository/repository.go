package repository

import (
	"context"
	"database/sql"
	_ "github.com/lib/pq"
	"todo/internal/models"
)

type Logger interface {
	Error(format string, v ...interface{})
	Warn(format string, v ...interface{})
	Info(format string, v ...interface{})
	Debug(format string, v ...interface{})
}

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
	db  *sql.DB
	Authorization
	TodoList
	TodoItem
	logger Logger
}

func NewRepository(cfg *Config, logger Logger) *Repository {
	return &Repository{
		cfg:    cfg,
		logger: logger,
	}
}

func (r *Repository) Run(_ context.Context) {
}

func (r *Repository) Stop() {
	err := r.db.Close()
	if err != nil {
		r.logger.Error(err.Error())
	}
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
