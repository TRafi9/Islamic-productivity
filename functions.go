package main

import (
	"os"
	"time"

	"crypto/rand"
	"io"

	"net/smtp"

	"github.com/go-redis/redis"
	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
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

func hashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

func checkPasswordHash(password string, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

func currentTimeStampPostgres() string {
	currentTime := time.Now()
	formattedTime := currentTime.Format("2006-01-02 15:04:05")
	return formattedTime
}

func currentTimePlusHourPostgres() string {
	currentTime := time.Now()
	oneHourLater := currentTime.Add(time.Hour)
	formattedTime := oneHourLater.Format("2006-01-02 15:04:05")
	return formattedTime
}

func generateRandomCode(max int) string {
	// [...] means size of array is determined by the number of values provided in the array itself
	table := [...]byte{1, 2, 3, 4, 5, 6, 7, 8, 9}
	b := make([]byte, max)
	n, err := io.ReadAtLeast(rand.Reader, b, max)
	if n != max {
		panic(err)
	}
	for i := 0; i < len(b); i++ {
		b[i] = table[int(b[i])%len(table)]
	}
	return string(b)
}

func sendEmailVerification(verificationCode string) error {
	email := os.Getenv("VERIFICATION_EMAIL")
	verification_email_password := os.Getenv("VERIFICATION_EMAIL_PASSWORD")
	hostName = 
	auth := smtp.PlainAuth("", email, verification_email_password, "smtp.gmail.com")
	//https://www.youtube.com/watch?v=H0HZc4FgX7E&t=249s&ab_channel=CodingwithRobby
	//https://pkg.go.dev/net/smtp#example-PlainAuth

	// generate random code

	// add email client
	return nil
}
