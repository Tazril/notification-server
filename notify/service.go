package notify

import (
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	_ "github.com/mattn/go-sqlite3"
)

type NotificationService struct {
	db *sql.DB
}

func NewNotificationService() *NotificationService {
	db, err := sql.Open("sqlite3", "./notifications.db")
	if err != nil {
		panic(fmt.Sprintf("Failed to connect to database: %v", err))
	}

	// Create table if not exists
	createTableSQL := `
	CREATE TABLE IF NOT EXISTS notifications (
		id TEXT PRIMARY KEY,
		current_btc_price REAL,
		market_trade_volume REAL,
		intra_day_high_price REAL,
		market_cap REAL,
		state TEXT,
		active BOOLEAN,
		created_at DATETIME,
		updated_at DATETIME
	);`
	
	if _, err := db.Exec(createTableSQL); err != nil {
		panic(fmt.Sprintf("Failed to create table: %v", err))
	}

	return &NotificationService{
		db: db,
	}
}

func (ns *NotificationService) CreateNotification(
	currentBTCPrice,
	marketTradeVolume,
	intraDayHighPrice,
	marketCap float64,
) *Notification {
	notification := &Notification{
		ID:                uuid.NewString(),
		CurrentBTCPrice:   currentBTCPrice,
		MarketTradeVolume: marketTradeVolume,
		IntraDayHighPrice: intraDayHighPrice,
		MarketCap:         marketCap,
		State:             CREATED,
		Active:            true,
		CreatedAt:         time.Now(),
		UpdatedAt:         time.Now(),
	}

	insertSQL := `INSERT INTO notifications (id, current_btc_price, market_trade_volume, intra_day_high_price, market_cap, state, active, created_at, updated_at)
	VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)`
	
	_, err := ns.db.Exec(insertSQL, notification.ID, notification.CurrentBTCPrice, notification.MarketTradeVolume, 
		notification.IntraDayHighPrice, notification.MarketCap, string(notification.State), notification.Active, 
		notification.CreatedAt, notification.UpdatedAt)
	
	if err != nil {
		panic(fmt.Sprintf("Failed to insert notification: %v", err))
	}

	return notification
}

func (ns *NotificationService) SendNotification(
	notificationId string,
	email string,
) error {
	// Get notification
	notification, err := ns.getNotificationByID(notificationId)
	if err != nil {
		return errors.New("cannot send notification, notification not found in the db")
	}

	if !notification.Active {
		return errors.New("cannot send notification, notification has been deleted")
	}

	// Send notification
	message := ns.getMessage(notification)
	err = ns.send(message, email)

	// Update state
	var newState NotificationState
	if err == nil {
		newState = SENT
	} else {
		newState = FAILED
	}

	updateSQL := `UPDATE notifications SET state = ?, updated_at = ? WHERE id = ?`
	_, dbErr := ns.db.Exec(updateSQL, string(newState), time.Now(), notificationId)
	if dbErr != nil {
		return fmt.Errorf("failed to update notification state: %v", dbErr)
	}

	return err
}

func (ns *NotificationService) GetNotifications(state string) ([]*Notification, error) {
	if state != "" && state != string(CREATED) && state != string(SENT) && state != string(FAILED) {
		return nil, errors.New("invalid state detected")
	}

	var query string
	var args []interface{}

	if state == "" {
		query = `SELECT id, current_btc_price, market_trade_volume, intra_day_high_price, market_cap, state, active, created_at, updated_at 
		FROM notifications WHERE active = 1`
	} else {
		query = `SELECT id, current_btc_price, market_trade_volume, intra_day_high_price, market_cap, state, active, created_at, updated_at 
		FROM notifications WHERE active = 1 AND state = ?`
		args = append(args, state)
	}

	rows, err := ns.db.Query(query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to query notifications: %v", err)
	}
	defer rows.Close()

	var notifications []*Notification
	for rows.Next() {
		notification := &Notification{}
		var stateStr string
		err := rows.Scan(&notification.ID, &notification.CurrentBTCPrice, &notification.MarketTradeVolume,
			&notification.IntraDayHighPrice, &notification.MarketCap, &stateStr, &notification.Active,
			&notification.CreatedAt, &notification.UpdatedAt)
		if err != nil {
			return nil, fmt.Errorf("failed to scan notification: %v", err)
		}
		notification.State = NotificationState(stateStr)
		notifications = append(notifications, notification)
	}

	return notifications, nil
}

func (ns *NotificationService) DeleteNotification(notificationId string) error {
	// Check if notification exists and is active
	notification, err := ns.getNotificationByID(notificationId)
	if err != nil {
		return errors.New("notification not found in the db")
	}

	if !notification.Active {
		return errors.New("notification has already been deleted")
	}

	// Soft delete by setting active to false
	updateSQL := `UPDATE notifications SET active = 0, updated_at = ? WHERE id = ?`
	_, err = ns.db.Exec(updateSQL, time.Now(), notificationId)
	if err != nil {
		return fmt.Errorf("failed to delete notification: %v", err)
	}

	return nil
}

func (ns *NotificationService) getNotificationByID(id string) (*Notification, error) {
	query := `SELECT id, current_btc_price, market_trade_volume, intra_day_high_price, market_cap, state, active, created_at, updated_at 
	FROM notifications WHERE id = ?`
	
	row := ns.db.QueryRow(query, id)
	notification := &Notification{}
	var stateStr string
	
	err := row.Scan(&notification.ID, &notification.CurrentBTCPrice, &notification.MarketTradeVolume,
		&notification.IntraDayHighPrice, &notification.MarketCap, &stateStr, &notification.Active,
		&notification.CreatedAt, &notification.UpdatedAt)
	
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("notification not found")
		}
		return nil, fmt.Errorf("database error: %v", err)
	}
	
	notification.State = NotificationState(stateStr)
	return notification, nil
}

func (ns *NotificationService) getMessage(notification *Notification) string {
	return fmt.Sprintf(`Hi There 
	current price of BTC is %.2f
	market trade volume is  %.2f
	intra day high price is  %.2f
	market  is  %.2f`,
		notification.CurrentBTCPrice,
		notification.MarketTradeVolume,
		notification.IntraDayHighPrice,
		notification.MarketCap)
}

func (ns *NotificationService) send(message string, email string) error {
	fmt.Printf("Notification message: %s \nsent to %v\n", message, email)
	return nil
}

func (ns *NotificationService) Close() error {
	return ns.db.Close()
}

