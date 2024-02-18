package main

import (
	"fmt"
	"math/big"
	"os"
	"strconv"
	"time"

	"crypto/rand"

	"net/http"
	"net/smtp"

	"github.com/go-redis/redis"
	"github.com/labstack/echo/v4"
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

func generateRandomCode(codeLength int) (int, error) {
	if codeLength <= 0 {
		return 0, fmt.Errorf("codeLength must be a positive integer")
	}

	// Determine the maximum value for the given codeLength
	maxValue := big.NewInt(10)
	maxValue.Exp(maxValue, big.NewInt(int64(codeLength)), nil)

	// Generate a random number within the specified range
	randomNum, err := rand.Int(rand.Reader, maxValue)
	if err != nil {
		return 0, err
	}

	return int(randomNum.Int64()), nil
}

func sendEmailVerification(c echo.Context, verificationCode int, logger *zap.SugaredLogger) error {
	email_sender := os.Getenv("VERIFICATION_EMAIL")
	verification_email_password := os.Getenv("VERIFICATION_EMAIL_PASSWORD")

	email_recipient := []string{"talhar9@gmail.com"}
	auth := smtp.PlainAuth("", email_sender, verification_email_password, "smtp.gmail.com")

	subject := "The Productive Muslim Verification Code"
	body := fmt.Sprintf(`Hey!
Here is your verification code: %s
Please enter it at the following link to register your account!
%s
Note that the code will expire within an hour!
Kind regards,
The Productive Muslim team`, strconv.Itoa(verificationCode), "http://localhost:3000/verify_email_view")

	msg := fmt.Sprintf("From: %s\nTo: %s \nSubject: %s\n\n%s",
		email_sender, "talhar9@gmail.com", subject, body)

	err := smtp.SendMail("smtp.gmail.com:587",
		auth,
		email_sender,
		email_recipient,
		[]byte(msg),
	)
	if err != nil {
		logger.Errorf("could not send email to recipient err: %s", err.Error())
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "failed to send verification email to recipient"})
	}

	//https://www.youtube.com/watch?v=H0HZc4FgX7E&t=249s&ab_channel=CodingwithRobby
	//https://pkg.go.dev/net/smtp#example-PlainAuth

	// generate random code

	// add email client
	return nil
}
