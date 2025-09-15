package repository

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	"strings"
	"todo/internal/models"
)

type TodoItemPostgres struct {
	db *sql.DB
}

func NewTodoItemPostgres(db *sql.DB) *TodoItemPostgres {
	return &TodoItemPostgres{db: db}
}

func (r *TodoItemPostgres) Create(listId int, item models.TodoItem) (int, error) {
	tx, err := r.db.Begin()
	if err != nil {
		return 0, err
	}

	var itemId int
	createItemQuery := fmt.Sprintf("INSERT INTO %s (title, description) values ($1, $2) RETURNING id", todoItemsTable)

	row := tx.QueryRow(createItemQuery, item.Title, item.Description)
	err = row.Scan(&itemId)
	if err != nil {
		tx.Rollback()
		return 0, err
	}

	createListItemsQuery := fmt.Sprintf("INSERT INTO %s (list_id, item_id) values ($1, $2)", listsItemsTable)
	_, err = tx.Exec(createListItemsQuery, listId, itemId)
	if err != nil {
		tx.Rollback()
		return 0, err
	}

	return itemId, tx.Commit()
}

func (r *TodoItemPostgres) GetAll(userId, listId int) ([]models.TodoItem, error) {
	query := fmt.Sprintf(`SELECT ti.id, ti.title, ti.description, ti.done FROM %s ti 
                          INNER JOIN %s li ON li.item_id = ti.id
                          INNER JOIN %s ul ON ul.list_id = li.list_id 
                          WHERE li.list_id = $1 AND ul.user_id = $2`,
		todoItemsTable, listsItemsTable, usersListsTable)

	rows, err := r.db.Query(query, listId, userId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var items []models.TodoItem
	for rows.Next() {
		var item models.TodoItem
		if err := rows.Scan(&item.Id, &item.Title, &item.Description, &item.Done); err != nil {
			return nil, err
		}
		items = append(items, item)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return items, nil
}

func (r *TodoItemPostgres) GetById(userId, itemId int) (models.TodoItem, error) {
	var item models.TodoItem

	query := fmt.Sprintf(`
		SELECT ti.id, ti.title, ti.description, ti.done 
		FROM %s ti 
		INNER JOIN %s li ON li.item_id = ti.id
		INNER JOIN %s ul ON ul.list_id = li.list_id 
		WHERE ti.id = $1 AND ul.user_id = $2`,
		todoItemsTable, listsItemsTable, usersListsTable)

	row := r.db.QueryRow(query, itemId, userId)

	err := row.Scan(&item.Id, &item.Title, &item.Description, &item.Done)
	if err != nil {
		if err == sql.ErrNoRows {
			// Если запись не найдена, можно вернуть кастомную ошибку
			return item, fmt.Errorf("item not found")
		}
		return item, err
	}

	return item, nil
}

func (r *TodoItemPostgres) Delete(userId, itemId int) error {
	query := `
		DELETE FROM todo_items ti 
		USING list_items li, user_lists ul 
		WHERE ti.id = li.item_id 
		  AND li.list_id = ul.list_id 
		  AND ul.user_id = $1 
		  AND ti.id = $2`

	_, err := r.db.Exec(query, userId, itemId)
	if err != nil {
		return fmt.Errorf("failed to delete todo item: %w", err)
	}

	return nil
}

func (r *TodoItemPostgres) Update(userId, itemId int, input models.UpdateItemInput) error {
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

	if input.Done != nil {
		setValues = append(setValues, fmt.Sprintf("done=$%d", argId))
		args = append(args, *input.Done)
		argId++
	}

	setQuery := strings.Join(setValues, ", ")

	query := fmt.Sprintf(`UPDATE %s ti SET %s FROM %s li, %s ul
									WHERE ti.id = li.item_id AND li.list_id = ul.list_id AND ul.user_id = $%d AND ti.id = $%d`,
		todoItemsTable, setQuery, listsItemsTable, usersListsTable, argId, argId+1)
	args = append(args, userId, itemId)

	_, err := r.db.Exec(query, args...)
	return err
}
