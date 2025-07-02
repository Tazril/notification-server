package main

import (
	"bitgo-go/notify"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

func localTesting() {
	ns := notify.NewNotificationService()

	notfication := ns.CreateNotification(120.0, 7000.0, 10000.0, 500.0)
	ns.CreateNotification(110.0, 7000.0, 10000.0, 500.0)

	if err := ns.SendNotification(notfication.ID, "hello@gmail.com"); err != nil {
		fmt.Println(err.Error())
	}

	notifications, err := ns.GetNotifications("")

	if err != nil {
		fmt.Println(err.Error())
	}

	fmt.Println("printing notifications")
	for _, notification := range notifications {
		fmt.Printf("notification: %#v\n", *notification)
	}
	fmt.Println()

	if err := ns.DeleteNotification(notfication.ID); err != nil {
		fmt.Println(err.Error())
	}

	if err := ns.SendNotification(notfication.ID, "hello@gmail.com"); err != nil {
		fmt.Println(err.Error())
	}
}

func homePage(c *gin.Context) {
    c.String(http.StatusOK, "Welcome to the Bitcoin Notification Server")
}

func main() {
	// localTesting()

	// Create notification service
	notificationService := notify.NewNotificationService()
	
	// Create notification handler
	notificationHandler := notify.NewNotificationServiceHandler(notificationService)

	// Create Gin router
	router := gin.Default()

	// Register Routes
	router.GET("/", homePage)
	
	// Health check endpoint
	router.GET("/health", notificationHandler.HealthCheck)
	
	// Notification API routes
	notificationRoutes := router.Group("/notifications")
	{
		notificationRoutes.POST("", notificationHandler.CreateNotification)
		notificationRoutes.GET("", notificationHandler.GetNotifications)
		notificationRoutes.POST("/:id/send", notificationHandler.SendNotification)
		notificationRoutes.DELETE("/:id", notificationHandler.DeleteNotification)
	}

	// Start the server
	fmt.Println("Starting Bitcoin Notification Server on :8080")
	fmt.Println("Available endpoints:")
	fmt.Println("  GET  /                     - Home page")
	fmt.Println("  GET  /health               - Health check")
	fmt.Println("  POST /notifications        - Create notification")
	fmt.Println("  GET  /notifications        - Get notifications")
	fmt.Println("  POST /notifications/:id/send - Send notification")
	fmt.Println("  DELETE /notifications/:id  - Delete notification")
	
	router.Run()
}