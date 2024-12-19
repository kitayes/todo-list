package models

type TodoList struct {
	Title       string `json:"title" db:"title" binding:"required"`
	Description string `json:"description" db:"description"`
	Id          int    `json:"id" db:"id"`
}
