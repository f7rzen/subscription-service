package handler

import (
	"log/slog"
	"net/http"
	"time"

	"github.com/f7rzen/subscription-service/internal/service"
	"github.com/gin-gonic/gin"
)

type SubscriptionHandler struct {
	service *service.SubscriptionService
	logger  *slog.Logger
}

func NewSubscriptionHandler(service *service.SubscriptionService, logger *slog.Logger) *SubscriptionHandler {
	return &SubscriptionHandler{
		service: service,
		logger:  logger,
	}
}

func (h *SubscriptionHandler) Create(c *gin.Context) {
	var input service.SubscriptionInput

	if err := c.ShouldBindJSON(&input); err != nil {
		h.logger.Warn("invalid request body", slog.String("error", err.Error()))

		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid request body",
		})
		return
	}

	subscription, err := h.service.Create(c.Request.Context(), input)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"id":           subscription.ID,
		"service_name": subscription.ServiceName,
		"price":        subscription.Price,
		"user_id":      subscription.UserID,
		"start_date":   subscription.StartDate.Format("01-2006"),
		"end_date":     formatOptionalDate(subscription.EndDate),
		"created_at":   subscription.CreatedAt,
		"updated_at":   subscription.UpdatedAt,
	})
}

func formatOptionalDate(date *time.Time) *string {
	if date == nil {
		return nil
	}

	formatted := date.Format("01-2006")
	return &formatted
}
