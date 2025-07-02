package notify

import "time"

// Create a notification. Line items may include current price of BTC, market trade volume, intra day high price, market cap
// Send a notification to an email
// List sent notifications (sent, outstanding, failed etc.)
// Delete a notification

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
