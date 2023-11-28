package main

import (
	"fmt"

	"github.com/labstack/echo"
	"go.uber.org/zap"
)

func main() {

	e := echo.New()
	z, _ := zap.NewProduction()
	logger := z.Sugar()

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

	api := e.Group("/api/v1")

	// v1.GET("/prayerTimes", getPrayerTimes(c*gin.Context))
	// check if this will be a post

	api.GET("/getPrayerTimes/:dateValue", func(c echo.Context) error {
		return todayPrayerHandler(c, Pt, logger)
	})

	e.Start(":8080")
}
