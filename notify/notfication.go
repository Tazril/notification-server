package notify

import "time"

type NotificationState string

const (
	CREATED NotificationState = "CREATED"
	SENT NotificationState = "SENT"
	FAILED NotificationState = "FAILED"
)

type Notification struct {
	ID string
	CurrentBTCPrice float64
	MarketTradeVolume float64
	IntraDayHighPrice float64
	MarketCap float64
	State NotificationState // CREATED, SENT, FAILED
	Active bool
	CreatedAt time.Time
	UpdatedAt time.Time
}
