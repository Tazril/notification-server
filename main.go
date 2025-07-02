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
    c.String(http.StatusOK, "This is my home page")
}

func main() {

	// localTesting()

	// Create Gin router
	router := gin.Default()

	// Register Routes
	router.GET("/", homePage)

	// Start the server
	router.Run()



	
}