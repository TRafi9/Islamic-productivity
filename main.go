package main

import (
	"fmt"
	"os"

	"github.com/labstack/echo"
	"go.uber.org/zap"
)

func readFile(filepath string) (string, error) {
	content, err := os.ReadFile(filepath)
	if err != nil {
		return "", err
	}
	contentString := string(content)
	return contentString, nil
}

func main() {
	z, _ := zap.NewProduction()
	logger := z.Sugar()

	pass, err := readFile("./pass.txt")
	if err != nil {
		// do something here
		return
	}
	logger.Infof("pass %s", pass)
	// client := redis.NewClient(&redis.Options{
	// 	Addr:     "redis-13336.c304.europe-west1-2.gce.cloud.redislabs.com",
	// 	Password: pass,     // no password set
	// 	DB:       12040201, // use default DB
	// })

	e := echo.New()

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
