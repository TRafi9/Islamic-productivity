package main

import (
	"fmt"

	"github.com/gin-gonic/gin"
)

func main() {

	router := gin.Default()

	// prayerTimes, err := GetPrayerTimes("Europe/London")
	// if err != nil {
	// 	// TODO find out why you cant return an fmt.errorf here?
	// 	return

	// }
	var prayerTimes PrayerTimes

	location := "Europe/London"
	pt, err := GetPrayerTimes(location)
	if err != nil {
		fmt.Errorf("error executing GetPrayerTimes, err %w", err)
		return
	}
	// store it in the global variable to be accessed outside
	// anonymous function

	prayerTimes = pt
	fmt.Println(prayerTimes.Fajr)

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
