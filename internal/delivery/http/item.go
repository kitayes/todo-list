package delivery

import (
	"github.com/gin-gonic/gin"
	"todo/internal/models"

	"net/http"
	"strconv"
)

func (h *Handler) createItem(c *gin.Context) {
	userId, err := getUserId(c)
	if err != nil {
		h.newErrorResponse(c, err)
		return
	}

	listId, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		h.newErrorResponse(c, err)
		return
	}

	var input models.TodoItem
	if err := c.BindJSON(&input); err != nil {
		h.newErrorResponse(c, err)
		return
	}

	id, err := h.services.TodoItem.Create(userId, listId, input)
	if err != nil {
		h.newErrorResponse(c, err)
		return
	}

	c.JSON(http.StatusOK, map[string]interface{}{
		"id": id,
	})
}

func (h *Handler) getAllItems(c *gin.Context) {
	userId, err := getUserId(c)
	if err != nil {
		h.newErrorResponse(c, err)
		return
	}

	listId, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		h.newErrorResponse(c, err)
		return
	}

	items, err := h.services.TodoItem.GetAll(userId, listId)
	if err != nil {
		h.newErrorResponse(c, err)
		return
	}

	c.JSON(http.StatusOK, items)
}

func (h *Handler) getItemById(c *gin.Context) {
	userId, err := getUserId(c)
	if err != nil {
		h.newErrorResponse(c, err)
		return
	}

	itemId, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		h.newErrorResponse(c, err)
		return
	}

	item, err := h.services.TodoItem.GetById(userId, itemId)
	if err != nil {
		h.newErrorResponse(c, err)
		return
	}

	c.JSON(http.StatusOK, item)
}

func (h *Handler) updateItem(c *gin.Context) {
	userId, err := getUserId(c)
	if err != nil {
		h.newErrorResponse(c, err)
		return
	}

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		h.newErrorResponse(c, err)
		return
	}

	var input models.UpdateItemInput
	if err := c.BindJSON(&input); err != nil {
		h.newErrorResponse(c, err)
		return
	}

	if err := h.services.TodoItem.Update(userId, id, input); err != nil {
		h.newErrorResponse(c, err)
		return
	}

	c.JSON(http.StatusOK, statusResponse{"ok"})
}

func (h *Handler) deleteItem(c *gin.Context) {
	userId, err := getUserId(c)
	if err != nil {
		h.newErrorResponse(c, err)
		return
	}

	itemId, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		h.newErrorResponse(c, err)
		return
	}

	err = h.services.TodoItem.Delete(userId, itemId)
	if err != nil {
		h.newErrorResponse(c, err)
		return
	}

	c.JSON(http.StatusOK, statusResponse{"ok"})
}
