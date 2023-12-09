package main

import (
	"encoding/json"
	"fmt"
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
	Pt, err := GetPrayerTimes(location, client, logger)
	if err != nil {
		fmt.Errorf("error executing GetPrayerTimes, err %w", err)
	}
	//TODO add Pt to redis memory
	// convert to strings as redis cannot store time vars

	newMap := make(map[string]map[string]string)
	for key, innerKeyVal := range Pt {
		newMap[key] = make(map[string]string)
		for prayer, timeVal := range innerKeyVal {
			newMap[key][prayer] = timeVal.String()
		}

	}

	for outerKey, innerMap := range newMap {
		innerMapJson, err := json.Marshal(innerMap)
		if err != nil {
			logger.Error("error marshalling inner map: %w", err)
			//TODO check if return is right or should add continue here
			return
		}
		err = client.Set(outerKey, innerMapJson, 0).Err()
		if err != nil {
			logger.Errorf("error uploading outerKey %s and innerMapJson %s, error: %s", outerKey, innerMapJson, err)
			return
		}
		// logger.Info("Uploaded data %s : %s", outerKey, innerMapJson)
	}
	//TODO add panic and recover if it fails to upload to memory

	api := e.Group("/api/v1")

	api.GET("/getPrayerTimes/:dateValue", func(c echo.Context) error {
		return todayPrayerHandler(c, Pt, logger)
	})

	e.Start(":8080")
}
