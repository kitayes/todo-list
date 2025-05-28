package delivery

import (
	"net/http"
	"strconv"
	"todo/internal/models"

	"github.com/gin-gonic/gin"
)

func (h *Handler) createList(c *gin.Context) {
	userId, err := getUserId(c)
	if err != nil {
		h.newErrorResponse(c, err)
		return
	}

	var input models.TodoList
	if err := c.BindJSON(&input); err != nil {
		h.newErrorResponse(c, err)
		return
	}

	id, err := h.services.TodoList.Create(userId, input)
	if err != nil {
		h.newErrorResponse(c, err)
		return
	}

	c.JSON(http.StatusOK, map[string]interface{}{
		"id": id,
	})
}

type getAllListsResponse struct {
	Data []models.TodoList `json:"data"`
}

func (h *Handler) getAllLists(c *gin.Context) {
	userId, err := getUserId(c)
	if err != nil {
		h.newErrorResponse(c, err)
		return
	}

	lists, err := h.services.TodoList.GetAll(userId)
	if err != nil {
		h.newErrorResponse(c, err)
		return
	}

	c.JSON(http.StatusOK, getAllListsResponse{
		Data: lists,
	})
}

func (h *Handler) getListById(c *gin.Context) {
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

	list, err := h.services.TodoList.GetById(userId, id)
	if err != nil {
		h.newErrorResponse(c, err)
		return
	}

	c.JSON(http.StatusOK, list)
}

func (h *Handler) updateList(c *gin.Context) {
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

	var input models.UpdateListInput
	if err := c.BindJSON(&input); err != nil {
		h.newErrorResponse(c, err)
		return
	}

	if err := h.services.TodoList.Update(userId, id, input); err != nil {
		h.newErrorResponse(c, err)
		return
	}

	c.JSON(http.StatusOK, statusResponse{"ok"})
}

func (h *Handler) deleteList(c *gin.Context) {
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

	err = h.services.TodoList.Delete(userId, id)
	if err != nil {
		h.newErrorResponse(c, err)
		return
	}

	c.JSON(http.StatusOK, statusResponse{
		Status: "ok",
	})
}
