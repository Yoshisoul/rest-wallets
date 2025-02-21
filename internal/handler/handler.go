package handler

import (
	"github.com/Yoshisoul/rest-wallets/internal/service"
	"github.com/gin-gonic/gin"
)

type Handler struct {
	services *service.Service
}

func NewHandler(services *service.Service) *Handler {
	return &Handler{services: services}
}

func (h *Handler) InitRoutes() *gin.Engine {
	router := gin.New()

	auth := router.Group("/auth")
	{
		auth.POST("/sign-up", h.signUp)
		auth.POST("/sign-in", h.signIn)
	}

	api := router.Group("/api/v1")
	{
		wallets := api.Group("/wallets", h.userIdentity)
		{
			wallets.POST("/", h.createWallet)
			wallets.GET("/", h.getAllWalletsFromUser)
			wallets.GET("/:id", h.getWalletById)
			wallets.DELETE("/:id", h.deleteWallet)
			// updates using transactions
		}

		transcactions := api.Group("/transactions")
		{
			transcactions.POST("/", h.createTransaction)
			transcactions.GET("/", h.getAllTransactions)
			transcactions.GET("/:id", h.getTransactionById)
			// can't update and delete transactions
		}
	}

	return router
}
