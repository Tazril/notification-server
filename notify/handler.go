package notify

import (
	"net/http"
	"github.com/gin-gonic/gin"
)

type NotificationServiceHandler struct {
	notificationService *NotificationService
}

func NewNotificationServiceHandler(notificationService *NotificationService) *NotificationServiceHandler {
    return &NotificationServiceHandler{
		notificationService,
    }
}

type CreateNotificationRequest struct {
	CurrentBTCPrice   float64 `json:"current_btc_price" binding:"required"`
	MarketTradeVolume float64 `json:"market_trade_volume" binding:"required"`
	IntraDayHighPrice float64 `json:"intra_day_high_price" binding:"required"`
	MarketCap         float64 `json:"market_cap" binding:"required"`
}

type SendNotificationRequest struct {
	Email string `json:"email" binding:"required,email"`
}

func (h *NotificationServiceHandler) CreateNotification(c *gin.Context) {
	var req CreateNotificationRequest
	
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	
	notification := h.notificationService.CreateNotification(
		req.CurrentBTCPrice,
		req.MarketTradeVolume,
		req.IntraDayHighPrice,
		req.MarketCap,
	)
	
	c.JSON(http.StatusCreated, gin.H{
		"message": "Notification created successfully",
		"notification": notification,
	})
}

func (h *NotificationServiceHandler) SendNotification(c *gin.Context) {
	notificationID := c.Param("id")
	
	var req SendNotificationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	
	if err := h.notificationService.SendNotification(notificationID, req.Email); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	
	c.JSON(http.StatusOK, gin.H{
		"message": "Notification sent successfully",
	})
}

func (h *NotificationServiceHandler) GetNotifications(c *gin.Context) {
	state := c.Query("state") // Optional query parameter
	
	notifications, err := h.notificationService.GetNotifications(state)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	
	c.JSON(http.StatusOK, gin.H{
		"notifications": notifications,
		"count": len(notifications),
	})
}

func (h *NotificationServiceHandler) DeleteNotification(c *gin.Context) {
	notificationID := c.Param("id")
	
	if err := h.notificationService.DeleteNotification(notificationID); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	
	c.JSON(http.StatusOK, gin.H{
		"message": "Notification deleted successfully",
	})
}

func (h *NotificationServiceHandler) HealthCheck(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status": "healthy",
		"service": "notification-server",
	})
}
