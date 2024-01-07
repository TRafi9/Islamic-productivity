package main

import (
	"time"

	"github.com/go-redis/redis"
	"go.uber.org/zap"
)

// have to use customError as scheduler doesnt accept
// echo context as a param so need to eliminate echo err responses
// and use custom errors instead
// the other choice was to wrap in a handlerfunc
type CustomError struct {
	Message string
}

func (e *CustomError) Error() string {
	return e.Message
}

func prayerTimesCronJob(client *redis.Client, logger *zap.SugaredLogger, location string, Pt map[string]map[string]time.Time) (latestPt map[string]map[string]time.Time, err error) {
	// get current date in Pt stored date format to check if exists in Pt, if it doesnt exist in Pt, run GetPrayerTimes to update Pt with new values
	// and store all the data in redis.
	// the frontend will refresh daily on a cron job (need to add) which will get the prayer times of the new day from Pt
	logger.Info("Cron task running to update Pt var/ redis cache if needed...")
	london, err := time.LoadLocation("Europe/London")
	if err != nil {
		logger.Errorf("Error loading location: %w", err)
		return nil, err
	}
	apiDate := time.Now().In(london)
	apiDate = time.Date(apiDate.Year(), apiDate.Month(), 01, 0, 0, 0, 0, apiDate.Location())
	redisDateFormat := apiDate.Format("2006-01-02")

	_, ok := Pt[redisDateFormat]
	if !ok {
		// Could not find data in Pt var for this date, and therefore needs to update Pt var, can just run GetPrayerTimes again for this
		logger.Info("No prayer time existed for first day of month, updating prayer times in redis and in Pt variable...")
		latestPt, err := GetPrayerTimes(location, client, logger)
		if err != nil {
			logger.Errorf("Error running getPrayerTimes from cronjob run err: %w", err)

			return nil, &CustomError{Message: "Error updating Pt"}
		}

		return latestPt, nil
	}
	logger.Info("Cron job ran today, prayer times do not need updating")
	return nil, &CustomError{Message: "No need to update Pt, prayer times do not need updating, as Pt[redisDateFormat] exists"}
}
