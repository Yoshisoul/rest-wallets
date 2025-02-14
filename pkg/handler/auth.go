package handler

import (
	"database/sql"
	"errors"
	"net/http"

	"github.com/Yoshisoul/rest-wallets/pkg/models"
	"github.com/gin-gonic/gin"
)

func (h *Handler) signUp(c *gin.Context) {
	var input models.SignUpInput

	if err := c.BindJSON(&input); err != nil {
		newErrorResponse(c, http.StatusBadRequest, "invalid input body")
		return
	}

	id, err := h.services.Authorization.CreateUser(input)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, "service failure")
		return
	}

	c.JSON(http.StatusOK, map[string]interface{}{
		"id": id,
	})
}

func (h *Handler) signIn(c *gin.Context) {
	var input models.SignInInput

	if err := c.BindJSON(&input); err != nil {
		newErrorResponse(c, http.StatusBadRequest, "invalid input body")
		return
	}

	token, err := h.services.Authorization.GenerateToken(input.Username, input.Password)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			newErrorResponse(c, http.StatusNotFound, "invalid username or password")
			return
		}
		newErrorResponse(c, http.StatusInternalServerError, "service failure")
		return
	}

	c.JSON(http.StatusOK, map[string]interface{}{
		"token": token,
	})
}
