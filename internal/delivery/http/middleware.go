package delivery

import (
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	"strings"
	"todo/internal/models"
)

const (
	authorizationHeader = "Authorization"
	userCtx             = "userId"
)

func (h *Handler) userIdentity(c *gin.Context) {
	header := c.GetHeader(authorizationHeader)
	if header == "" {
		h.newErrorResponse(c, models.ErrUnauthorized)
		return
	}

	headerParts := strings.Split(header, " ")
	if len(headerParts) != 2 || headerParts[0] != "Bearer" {
		h.newErrorResponse(c, models.ErrUnauthorized)
		return
	}

	if len(headerParts[1]) == 0 {
		h.newErrorResponse(c, models.ErrUnauthorized)
		return
	}

	userId, err := h.services.Authorization.ParseToken(headerParts[1])
	if err != nil {
		h.newErrorResponse(c, err)
		return
	}

	c.Set(userCtx, userId)
}

func getUserId(c *gin.Context) (int, error) {
	id, ok := c.Get(userCtx)
	if !ok {
		return 0, errors.New("user id not found")
	}

	idInt, ok := id.(int)
	if !ok {
		return 0, errors.Wrap(models.ErrInvalidInput, "user id is of invalid type")
	}

	return idInt, nil
}

// osh valid invalidinp, jwt unauth, notfound sql
