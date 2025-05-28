package repository

import (
	"database/sql"
	_ "database/sql"
	"fmt"
	_ "github.com/lib/pq"
	"strings"
	"todo/internal/models"
)

type TodoListPostgres struct {
	db *sql.DB
}

func NewTodoListPostgres(db *sql.DB) *TodoListPostgres {
	return &TodoListPostgres{db: db}
}

// idk about tl(todolist)

func (r *TodoListPostgres) Create(userId int, list models.TodoList) (int, error) {
	tx, err := r.db.Begin()
	if err != nil {
		return 0, err
	}

	var id int
	createListQuery := fmt.Sprintf("INSERT INTO %s (title, description) VALUES ($1, $2) RETURNING id", todoListsTable)
	row := tx.QueryRow(createListQuery, list.Title, list.Description)
	if err := row.Scan(&id); err != nil {
		tx.Rollback()
		return 0, err
	}

	createUsersListQuery := fmt.Sprintf("INSERT INTO %s (user_id, list_id) VALUES ($1, $2)", usersListsTable)
	_, err = tx.Exec(createUsersListQuery, userId, id)
	if err != nil {
		tx.Rollback()
		return 0, err
	}

	return id, tx.Commit()
}

func (r *TodoListPostgres) GetAll(userId int) ([]models.TodoList, error) {
	query := fmt.Sprintf(`
		SELECT tl.id, tl.title, tl.description
		FROM %s tl
		INNER JOIN %s ul ON tl.id = ul.list_id
		WHERE ul.user_id = $1`,
		todoListsTable, usersListsTable)

	rows, err := r.db.Query(query, userId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var lists []models.TodoList
	for rows.Next() {
		var list models.TodoList
		err := rows.Scan(&list.Id, &list.Title, &list.Description)
		if err != nil {
			return nil, err
		}
		lists = append(lists, list)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return lists, nil
}

func (r *TodoListPostgres) GetById(userId, listId int) (models.TodoList, error) {
	query := fmt.Sprintf(`
		SELECT tl.id, tl.title, tl.description
		FROM %s tl
		INNER JOIN %s ul ON tl.id = ul.list_id
		WHERE ul.user_id = $1 AND ul.list_id = $2`,
		todoListsTable, usersListsTable)

	var list models.TodoList

	err := r.db.QueryRow(query, userId, listId).Scan(&list.Id, &list.Title, &list.Description)
	if err != nil {
		return models.TodoList{}, err
	}

	return list, nil
}

func (r *TodoListPostgres) Delete(userId, listId int) error {
	query := fmt.Sprintf("DELETE FROM %s tl USING %s ul WHERE tl.id = ul.list_id AND ul.user_id=$1 AND ul.list_id=$2",
		todoListsTable, usersListsTable)
	_, err := r.db.Exec(query, userId, listId)

	return err
}

func (r *TodoListPostgres) Update(userId, listId int, input models.UpdateListInput) error {
	setValues := make([]string, 0)
	args := make([]interface{}, 0)
	argId := 1

	if input.Title != nil {
		setValues = append(setValues, fmt.Sprintf("title=$%d", argId))
		args = append(args, *input.Title)
		argId++
	}

	if input.Description != nil {
		setValues = append(setValues, fmt.Sprintf("description=$%d", argId))
		args = append(args, *input.Description)
		argId++
	}

	// title=$1
	// description=$1
	// title=$1, description=$2
	setQuery := strings.Join(setValues, ", ")

	query := fmt.Sprintf("UPDATE %s tl SET %s FROM %s ul WHERE tl.id = ul.list_id AND ul.list_id=$%d AND ul.user_id=$%d",
		todoListsTable, setQuery, usersListsTable, argId, argId+1)
	args = append(args, listId, userId)

	_, err := r.db.Exec(query, args...)
	return err
}
