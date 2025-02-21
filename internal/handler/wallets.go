package handler

import (
	"database/sql"
	"errors"
	"net/http"

	"github.com/Yoshisoul/rest-wallets/internal/models"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func (h *Handler) createWallet(c *gin.Context) {
	id, err := getUserId(c)
	if err != nil {
		return
	}

	uuid, err := h.services.Wallet.Create(id)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, "service failure")
		return
	}

	c.JSON(http.StatusOK, map[string]interface{}{
		"uuid": uuid,
	})
}

type getAllWalletsResponse struct {
	Wallets []models.Wallet `json:"data"`
}

func (h *Handler) getAllWalletsFromUser(c *gin.Context) {
	id, err := getUserId(c)
	if err != nil {
		return
	}

	wallets, err := h.services.Wallet.GetAllFromUser(id)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, getAllWalletsResponse{
		Wallets: wallets,
	})
}

func (h *Handler) getWalletById(c *gin.Context) {
	userId, err := getUserId(c)
	if err != nil {
		return
	}

	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, "invalid id param")
		return
	}

	wallet, err := h.services.Wallet.GetByIdFromUser(userId, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			newErrorResponse(c, http.StatusNotFound, "wallet not found")
			return
		}
		newErrorResponse(c, http.StatusInternalServerError, "service failure")
		return
	}

	c.JSON(http.StatusOK, wallet)
}

func (h *Handler) deleteWallet(c *gin.Context) {
	userId, err := getUserId(c)
	if err != nil {
		return
	}

	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, "invalid id param")
		return
	}

	err = h.services.Wallet.Delete(userId, id)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, statusResponse{
		Status: "ok",
	})
}
