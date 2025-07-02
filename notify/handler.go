package notify

import "github.com/gin-gonic/gin"

type NotificationServiceHandler struct {
	notificationService *NotificationService
}

func NewNotificationServiceHandler(notificationService *NotificationService) *NotificationServiceHandler {
    return &NotificationServiceHandler{
		notificationService,
    }
}

func (h* NotificationServiceHandler) CreateNotification(c *gin.Context) {
	
	
	
}
