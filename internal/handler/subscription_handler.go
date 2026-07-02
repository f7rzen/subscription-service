package handler

import (
	"log/slog"
	"net/http"
	"strconv"
	"time"

	"github.com/f7rzen/subscription-service/internal/model"
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

type subscriptionResponseDTO struct {
	ID          int64     `json:"id"`
	ServiceName string    `json:"service_name"`
	Price       int       `json:"price"`
	UserID      string    `json:"user_id"`
	StartDate   string    `json:"start_date"`
	EndDate     *string   `json:"end_date"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
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

	c.JSON(http.StatusCreated, subscriptionResponse(subscription))
}

func (h *SubscriptionHandler) GetByID(c *gin.Context) {
	id, ok := parseSubscriptionID(c)
	if !ok {
		return
	}

	subscription, err := h.service.GetByID(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "subscription not found",
		})
		return
	}

	c.JSON(http.StatusOK, subscriptionResponse(subscription))
}

func (h *SubscriptionHandler) List(c *gin.Context) {
	subscriptions, err := h.service.List(c.Request.Context())
	if err != nil {
		h.logger.Error("failed to list subscriptions", slog.String("error", err.Error()))

		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "failed to list subscriptions",
		})
		return
	}

	response := make([]subscriptionResponseDTO, 0, len(subscriptions))

	for _, subscription := range subscriptions {
		response = append(response, subscriptionResponse(subscription))
	}

	c.JSON(http.StatusOK, response)
}

func (h *SubscriptionHandler) Update(c *gin.Context) {
	id, ok := parseSubscriptionID(c)
	if !ok {
		return
	}

	var input service.SubscriptionInput

	if err := c.ShouldBindJSON(&input); err != nil {
		h.logger.Warn("invalid request body", slog.String("error", err.Error()))

		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid request body",
		})
		return
	}

	subscription, err := h.service.Update(c.Request.Context(), id, input)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, subscriptionResponse(subscription))
}

func (h *SubscriptionHandler) Delete(c *gin.Context) {
	id, ok := parseSubscriptionID(c)
	if !ok {
		return
	}

	if err := h.service.Delete(c.Request.Context(), id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "failed to delete subscription",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "subscription deleted",
	})
}

func (h *SubscriptionHandler) Summary(c *gin.Context) {
	userID := c.Query("user_id")
	serviceName := c.Query("service_name")
	from := c.Query("from")
	to := c.Query("to")

	totalPrice, err := h.service.CalculateSummary(
		c.Request.Context(),
		userID,
		serviceName,
		from,
		to,
	)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"total_price": totalPrice,
	})
}

func parseSubscriptionID(c *gin.Context) (int64, bool) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil || id <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid subscription id",
		})
		return 0, false
	}

	return id, true
}

func subscriptionResponse(subscription model.Subscription) subscriptionResponseDTO {
	return subscriptionResponseDTO{
		ID:          subscription.ID,
		ServiceName: subscription.ServiceName,
		Price:       subscription.Price,
		UserID:      subscription.UserID,
		StartDate:   subscription.StartDate.Format("01-2006"),
		EndDate:     formatOptionalDate(subscription.EndDate),
		CreatedAt:   subscription.CreatedAt,
		UpdatedAt:   subscription.UpdatedAt,
	}
}

func formatOptionalDate(date *time.Time) *string {
	if date == nil {
		return nil
	}

	formatted := date.Format("01-2006")
	return &formatted
}
