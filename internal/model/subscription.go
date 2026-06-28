package model

import "time"

type Subscription struct {
	ID          string     `db:"id"`
	ServiceName string     `db:"service_name"`
	Price       int        `db:"price"`
	UserID      string     `db:"user_id"`
	StartDate   time.Time  `db:"start_date"`
	EndDate     *time.Time `db:"end_date"`
	CreatedAt   time.Time  `db:"created_at"`
	UpdatedAt   time.Time  `db:"updated_at"`
}

type SummaryResponse struct {
	TotalPrice int `json:"total_price"`
}
