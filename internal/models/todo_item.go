package models

type TodoItem struct {
	Title       string `json:"title" db:"title" binding:"required"`
	Description string `json:"description" db:"description"`
	Id          int    `json:"id" db:"id"`
	Done        bool   `json:"done" db:"done"`
}
