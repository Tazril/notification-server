package notify

import (
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
)

type NotificationService struct {
	notifications map[string]*Notification
}

func NewNotificationService() *NotificationService {

	return &NotificationService{
		notifications: map[string]*Notification{},
	}

}

func (ns *NotificationService) CreateNotification(
	currentBTCPrice,
	marketTradeVolume,
	intraDayHighPrice,
	marketCap float64,
) *Notification {
	notfication := &Notification{
		ID:                uuid.NewString(),
		CurrentBTCPrice:   currentBTCPrice,
		MarketTradeVolume: marketTradeVolume,
		IntraDayHighPrice: intraDayHighPrice,
		MarketCap:         marketCap,
		State:             CREATED,
		Active:            true,
	}
	ns.notifications[notfication.ID] = notfication

	return notfication
}

func (ns *NotificationService) SendNotification(
	notificationId string,
	email string,
) error {

	notification, ok := ns.notifications[notificationId]

	if !ok {
		return errors.New("cannot send notification, notification not found in the db")
	}

	if !notification.Active {
		return errors.New("cannot send notification, notification has been deleted")
	}

	// send notification
	message := ns.getMessage(notification)
	err := ns.send(message, email)

	if err == nil {
		notification.State = SENT
	} else {
		notification.State = FAILED
	}
	notification.UpdatedAt = time.Now()

	return nil

}


func (ns *NotificationService) GetNotifications(state string) ([]*Notification, error) {

	if state != "" && state != string(CREATED) && state != string(SENT) && state  != string(FAILED) {
		return nil, errors.New("invalid state detected")
	}

	activeNotifications := []*Notification{}
	for _, notification := range ns.notifications {
		if notification.Active {
			if state == "" || notification.State == NotificationState(state) {
				activeNotifications = append(activeNotifications, notification)
			}
			
		}
	}

	return activeNotifications, nil
}

func (ns *NotificationService) DeleteNotification(notificationId string) error {

	notification, ok := ns.notifications[notificationId]

	if !ok {
		return errors.New("notification not found in the db")
	}

	if !notification.Active {
		return errors.New("notification has already been deleted")
	}

	notification.Active = false

	return nil
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
