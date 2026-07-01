package service

import (
	"context"
	"errors"
	"log/slog"
	"strings"
	"time"

	"github.com/f7rzen/subscription-service/internal/model"
	"github.com/f7rzen/subscription-service/internal/repository"
	"github.com/google/uuid"
)

type SubscriptionService struct {
	repo   *repository.SubscriptionRepository
	logger *slog.Logger
}

func NewSubscriptionService(repo *repository.SubscriptionRepository, logger *slog.Logger) *SubscriptionService {
	return &SubscriptionService{
		repo:   repo,
		logger: logger,
	}
}

type SubscriptionInput struct {
	ServiceName string  `json:"service_name"`
	Price       int     `json:"price"`
	UserID      string  `json:"user_id"`
	StartDate   string  `json:"start_date"`
	EndDate     *string `json:"end_date,omitempty"`
}

func (s *SubscriptionService) Create(ctx context.Context, input SubscriptionInput) (model.Subscription, error) {
	serviceName := strings.TrimSpace(input.ServiceName)
	if serviceName == "" {
		return model.Subscription{}, errors.New("service_name is required")
	}

	if input.Price <= 0 {
		return model.Subscription{}, errors.New("price must be greater than zero")
	}

	if err := uuid.Validate(input.UserID); err != nil {
		return model.Subscription{}, errors.New("invalid user_id")
	}

	startDate, err := time.Parse("01-2006", input.StartDate)
	if err != nil {
		return model.Subscription{}, errors.New("invalid start_date format, expected MM-YYYY")
	}

	var endDate *time.Time
	if input.EndDate != nil {
		parsedEndDate, err := time.Parse("01-2006", *input.EndDate)
		if err != nil {
			return model.Subscription{}, errors.New("invalid end_date format, expected MM-YYYY")
		}

		if parsedEndDate.Before(startDate) {
			return model.Subscription{}, errors.New("end_date must be greater than or equal to start_date")
		}

		endDate = &parsedEndDate
	}

	subscription := model.Subscription{
		ID:          uuid.NewString(),
		ServiceName: serviceName,
		Price:       input.Price,
		UserID:      input.UserID,
		StartDate:   startDate,
		EndDate:     endDate,
	}

	createdSubscription, err := s.repo.Create(ctx, subscription)
	if err != nil {
		s.logger.Error("failed to create subscription", slog.String("error", err.Error()))
		return model.Subscription{}, err
	}

	s.logger.Info(
		"subscription created",
		slog.String("id", createdSubscription.ID),
		slog.String("user_id", createdSubscription.UserID),
		slog.String("service_name", createdSubscription.ServiceName),
	)

	return createdSubscription, nil
}

func (s *SubscriptionService) GetByID(ctx context.Context, id string) (model.Subscription, error) {
	if err := uuid.Validate(id); err != nil {
		return model.Subscription{}, errors.New("invalid subscription id")
	}

	subscription, err := s.repo.GetByID(ctx, id)
	if err != nil {
		s.logger.Error("failed to get subscription", slog.String("id", id), slog.String("error", err.Error()))
		return model.Subscription{}, err
	}

	return subscription, nil
}

func (s *SubscriptionService) List(ctx context.Context) ([]model.Subscription, error) {
	subscriptions, err := s.repo.List(ctx)
	if err != nil {
		s.logger.Error("failed to list subscriptions", slog.String("error", err.Error()))
		return nil, err
	}

	return subscriptions, nil
}

func (s *SubscriptionService) Update(ctx context.Context, id string, input SubscriptionInput) (model.Subscription, error) {
	if err := uuid.Validate(id); err != nil {
		return model.Subscription{}, errors.New("invalid subscription id")
	}

	serviceName := strings.TrimSpace(input.ServiceName)
	if serviceName == "" {
		return model.Subscription{}, errors.New("service_name is required")
	}

	if input.Price <= 0 {
		return model.Subscription{}, errors.New("price must be greater than zero")
	}

	if err := uuid.Validate(input.UserID); err != nil {
		return model.Subscription{}, errors.New("invalid user_id")
	}

	startDate, err := time.Parse("01-2006", input.StartDate)
	if err != nil {
		return model.Subscription{}, errors.New("invalid start_date format, expected MM-YYYY")
	}

	var endDate *time.Time
	if input.EndDate != nil {
		parsedEndDate, err := time.Parse("01-2006", *input.EndDate)
		if err != nil {
			return model.Subscription{}, errors.New("invalid end_date format, expected MM-YYYY")
		}

		if parsedEndDate.Before(startDate) {
			return model.Subscription{}, errors.New("end_date must be greater than or equal to start_date")
		}

		endDate = &parsedEndDate
	}

	subscription := model.Subscription{
		ID:          id,
		ServiceName: serviceName,
		Price:       input.Price,
		UserID:      input.UserID,
		StartDate:   startDate,
		EndDate:     endDate,
	}

	updatedSubscription, err := s.repo.Update(ctx, subscription)
	if err != nil {
		s.logger.Error("failed to update subscription", slog.String("error", err.Error()))
		return model.Subscription{}, err
	}

	s.logger.Info(
		"subscription updated",
		slog.String("id", updatedSubscription.ID),
		slog.String("user_id", updatedSubscription.UserID),
		slog.String("service_name", updatedSubscription.ServiceName),
	)

	return updatedSubscription, nil
}

func (s *SubscriptionService) Delete(ctx context.Context, id string) error {
	if err := uuid.Validate(id); err != nil {
		return errors.New("invalid subscription id")
	}

	if err := s.repo.Delete(ctx, id); err != nil {
		s.logger.Error("failed to delete subscription", slog.String("error", err.Error()))
		return err
	}

	s.logger.Info("subscription deleted", slog.String("id", id))

	return nil
}

func (s *SubscriptionService) CalculateSummary(ctx context.Context, userID string, serviceName string, from string, to string) (int, error) {
	if err := uuid.Validate(userID); err != nil {
		return 0, errors.New("invalid user_id")
	}

	serviceName = strings.TrimSpace(serviceName)
	if serviceName == "" {
		return 0, errors.New("service_name is required")
	}

	fromDate, err := time.Parse("01-2006", from)
	if err != nil {
		return 0, errors.New("invalid from date format, expected MM-YYYY")
	}

	toDate, err := time.Parse("01-2006", to)
	if err != nil {
		return 0, errors.New("invalid to date format, expected MM-YYYY")
	}

	if toDate.Before(fromDate) {
		return 0, errors.New("to date must be greater than or equal to from date")
	}

	subscriptions, err := s.repo.FindByUserIDAndServiceName(ctx, userID, serviceName)
	if err != nil {
		s.logger.Error("failed to find subscriptions for summary", slog.String("error", err.Error()))
		return 0, err
	}

	totalPrice := 0

	for _, subscription := range subscriptions {
		subscriptionStart := subscription.StartDate
		subscriptionEnd := toDate

		if subscription.EndDate != nil {
			subscriptionEnd = *subscription.EndDate
		}

		overlapStart := maxDate(subscriptionStart, fromDate)
		overlapEnd := minDate(subscriptionEnd, toDate)

		if overlapEnd.Before(overlapStart) {
			continue
		}

		monthsCount := countMonthsInclusive(overlapStart, overlapEnd)
		totalPrice += subscription.Price * monthsCount
	}

	s.logger.Info(
		"summary calculated",
		slog.String("user_id", userID),
		slog.String("service_name", serviceName),
		slog.Int("total_price", totalPrice),
	)

	return totalPrice, nil
}

func maxDate(a time.Time, b time.Time) time.Time {
	if a.After(b) {
		return a
	}

	return b
}

func minDate(a time.Time, b time.Time) time.Time {
	if a.Before(b) {
		return a
	}

	return b
}

func countMonthsInclusive(start time.Time, end time.Time) int {
	startYear, startMonth, _ := start.Date()
	endYear, endMonth, _ := end.Date()

	return (endYear-startYear)*12 + int(endMonth-startMonth) + 1
}
