package main

import (
	"os"

	"github.com/go-redis/redis"
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
	client := redis.NewClient(&redis.Options{
		Addr:     "redis-13336.c304.europe-west1-2.gce.cloud.redislabs.com:13336",
		Password: pass,
		DB:       0,
	})

	e := echo.New()

	//TODO make getPrayerTimes return the whole month
	// use a cron job to run get prayertimes
	// use the same cron job to add the prayer times to a redis in memory database for that month
	// that way the data will persist if the application goes down
	// add v1 GET api to make a call, given a date, to recieve all prayer times for that day, from the redis server

	location := "Europe/London"
	//TODO make it globally readable,
	// use concurrency to
	// pt object stores everything, not have data for this month, 100 ppl call app at same time for new month, 100 calls, so go handles them concurrently,
	// problem , when handler cant find data

	Pt, err := GetPrayerTimes(location, client, logger)
	if err != nil {
		logger.Errorf("error executing GetPrayerTimes, err %w", err)
	}

	//TODO add panic and recover if it fails to upload to memory

	api := e.Group("/api/v1")

	api.GET("/getPrayerTimes/:dateValue", func(c echo.Context) error {
		return todayPrayerHandler(c, Pt, logger)
	})

	//TODO CONTINUE FROM HERE
	// make Pt a global variable, then use a cloud function to trigger the handler function underneath
	// handler function should just rerun getPrayerTimes, which will run for the new day and push the month data to redis
	// therefore you want the cloud function to trigger just after midnight on the first of a new month to populate redis data
	// use mutex to lock Pt while you update it, and then open it up once mutex is done

	// you can also use an infinite for-loop instead that will sleep daily and trigger Pt to run if it is the first of the month
	// but this is a copout version (although very viable and efficient heuheh)
	api.GET("/getPrayerTimes/:dateValue", func(c echo.Context) error {
		Pt, err := GetPrayerTimes(location, client, logger)
		return err
	})

	e.Start(":8080")
}
