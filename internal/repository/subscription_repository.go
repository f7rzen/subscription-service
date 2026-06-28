package repository

import (
	"context"

	"github.com/f7rzen/subscription-service/internal/model"
	"github.com/jmoiron/sqlx"
)

type SubscriptionRepository struct {
	db *sqlx.DB
}

func NewSubscriptionRepository(db *sqlx.DB) *SubscriptionRepository {
	return &SubscriptionRepository{
		db: db,
	}
}

func (r *SubscriptionRepository) Create(ctx context.Context, subscription model.Subscription) (model.Subscription, error) {
	query := `
		INSERT INTO subscriptions (
			id,
			service_name,
			price,
			user_id,
			start_date,
			end_date
		)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id, service_name, price, user_id, start_date, end_date, created_at, updated_at
	`

	var createdSubscription model.Subscription
	
	err := r.db.GetContext(
		ctx,
		&createdSubscription,
		query,
		subscription.ID,
		subscription.ServiceName,
		subscription.Price,
		subscription.UserID,
		subscription.StartDate,
		subscription.EndDate,
	)
	if err != nil {
		return model.Subscription{}, err
	}

	return createdSubscription, nil
}

func (r *SubscriptionRepository) GetByID(ctx context.Context, id string) (model.Subscription, error) {
	query := `
		SELECT id, service_name, price, user_id, start_date, end_date, created_at, updated_at
		FROM subscriptions
		WHERE id = $1
	`

	var subscription model.Subscription

	err := r.db.GetContext(ctx, &subscription, query, id)
	if err != nil {
		return model.Subscription{}, err
	}

	return subscription, nil
}

func (r *SubscriptionRepository) List(ctx context.Context) ([]model.Subscription, error) {
	query := `
		SELECT id, service_name, price, user_id, start_date, end_date, created_at, updated_at
		FROM subscriptions
		ORDER BY created_at DESC
	`

	var subscriptions []model.Subscription

	err := r.db.SelectContext(ctx, &subscriptions, query)
	if err != nil {
		return nil, err
	}

	return subscriptions, nil
}

func (r *SubscriptionRepository) Update(ctx context.Context, subscription model.Subscription) (model.Subscription, error) {
	query := `
		UPDATE subscriptions
		SET
			service_name = $1,
			price = $2,
			user_id = $3,
			start_date = $4,
			end_date = $5,
			updated_at = NOW()
		WHERE id = $6
		RETURNING id, service_name, price, user_id, start_date, end_date, created_at, updated_at
	`

	var updatedSubscription model.Subscription

	err := r.db.GetContext(
		ctx,
		&updatedSubscription,
		query,
		subscription.ServiceName,
		subscription.Price,
		subscription.UserID,
		subscription.StartDate,
		subscription.EndDate,
		subscription.ID,
	)
	if err != nil {
		return model.Subscription{}, err
	}

	return updatedSubscription, nil
}

func (r *SubscriptionRepository) Delete(ctx context.Context, id string) error {
	query := `
		DELETE FROM subscriptions
		WHERE id = $1
	`

	_, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return err
	}

	return nil
}

func (r *SubscriptionRepository) FindByUserIDAndServiceName(ctx context.Context, userID string, serviceName string) ([]model.Subscription, error) {
	query := `
		SELECT id, service_name, price, user_id, start_date, end_date, created_at, updated_at
		FROM subscriptions
		WHERE user_id = $1
		  AND service_name = $2
	`

	var subscriptions []model.Subscription

	err := r.db.SelectContext(ctx, &subscriptions, query, userID, serviceName)
	if err != nil {
		return nil, err
	}

	return subscriptions, nil
}
