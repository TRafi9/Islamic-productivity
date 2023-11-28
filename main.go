package main

import (
	"fmt"

	"github.com/gin-gonic/gin"
)

func main() {

	router := gin.Default()

	//TODO make getPrayerTimes return the whole month
	// use a cron job to run get prayertimes
	// use the same cron job to add the prayer times to a redis in memory database for that month
	// that way the data will persist if the application goes down
	// add v1 GET api to make a call, given a date, to recieve all prayer times for that day, from the redis server

	location := "Europe/London"
	Pt, err := GetPrayerTimes(location)
	if err != nil {
		fmt.Errorf("error executing GetPrayerTimes, err %w", err)
	}
	// store it in the global variable to be accessed outside
	fmt.Println(Pt)

	v1 := router.Group("/api/v1")
	{
		// v1.GET("/prayerTimes", getPrayerTimes(c*gin.Context))
		// check if this will be a post
		v1.GET("/getPrayerTimes/:dateValue", todayPrayerHandler)
	}
	router.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})
	router.Run(":8080")
}
