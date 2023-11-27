package main

import (
	"fmt"

	"github.com/gin-gonic/gin"
)

func main() {

	router := gin.Default()

	location := "Europe/London"
	pt, err := GetPrayerTimes(location)
	if err != nil {
		fmt.Errorf("error executing GetPrayerTimes, err %w", err)
	}
	// store it in the global variable to be accessed outside
	fmt.Println(pt.Isha)
	fmt.Println(pt.PrayerDate)

	v1 := router.Group("/v1")
	{
		// v1.GET("/prayerTimes", getPrayerTimes(c*gin.Context))
		// check if this will be a post
		v1.POST("/")
	}
	router.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})
	router.Run(":8080")
}
