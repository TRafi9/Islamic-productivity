package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/go-redis/redis"
	"github.com/labstack/echo"
	"go.uber.org/zap"
)

func GetPrayerTimes(location string, client *redis.Client, logger *zap.SugaredLogger) (map[string]map[string]time.Time, error) {

	// Gets prayer times monthly

	apiLoc := "london"
	london, err := time.LoadLocation("Europe/London")
	if err != nil {
		fmt.Println("Error loading location:", err)
		return nil, err
	}
	// apiDate & apiDateString used to get all prayer times from 1st of current month
	apiDate := time.Now().In(london)
	// set as first day of month, use for loop later to add all dates of month to a map
	apiDate = time.Date(apiDate.Year(), apiDate.Month(), 01, 0, 0, 0, 0, apiDate.Location())
	apiDateString := apiDate.Format("02-01-2006")
	redisDateFormat := apiDate.Format("2006-01-02")
	// redisDateFormat := "01-02-2006"
	fmt.Printf("date string: %s \n", apiDateString)

	res, err := client.Get(redisDateFormat).Result()

	if err != nil {
		// NO REDIS DATA/ REDIS PULL FAILED, CALL API INSTEAD
		// error could also be that redis.Nil aka there was no value returned
		logger.Errorf("error with redis get call, continue to get data from api, err:  %w", err)

		// if no redis records exist for the date being queried, call the api code here and also upload to db!
		if err == redis.Nil {
			logger.Infof("No redis entry for this date, date %s, res: %s", redisDateFormat, res)
		}

		prayerTimesURL := fmt.Sprintf("https://muslimsalat.com/%s/monthly/%s/true.json?key==dd00aaed7ee591ead148b3af566d88f1", apiLoc, apiDateString)

		response, err := http.Get(prayerTimesURL)
		if err != nil {
			return nil, fmt.Errorf("failed to get response from API call to prayerTimesURL, err: %w", err)
		}
		//TODO check why we defer this specifically
		defer response.Body.Close()

		resBody, err := io.ReadAll(response.Body)
		if err != nil {
			return nil, fmt.Errorf("failed to read response body, err: %w", err)
		}
		jsonString := string(resBody)
		var ResponseStruct ResponseStruct

		err = json.Unmarshal([]byte(jsonString), &ResponseStruct)
		if err != nil {
			return nil, fmt.Errorf("error unmarshalling json into struct, err: %w", err)
		}

		prayerMonthMap := make(map[string]map[string]time.Time)

		// loops through json for all days of month
		// finds today and gets prayer times for today in string
		for i := range ResponseStruct.Items {
			prayerDate := ResponseStruct.Items[i].DateFor

			parsedDate, err := time.Parse("2006-1-2", prayerDate)
			if err != nil {
				return nil, fmt.Errorf("error parsing date: %w", err)
			}
			prayerDate = parsedDate.Format("2006-01-02")

			prayerDayMap := make(map[string]time.Time)

			prayerDayMap["Fajr"] = parseTime(prayerDate, ResponseStruct.Items[i].Fajr, location)
			prayerDayMap["Dhuhr"] = parseTime(prayerDate, ResponseStruct.Items[i].Dhuhr, location)
			prayerDayMap["Asr"] = parseTime(prayerDate, ResponseStruct.Items[i].Asr, location)
			prayerDayMap["Maghrib"] = parseTime(prayerDate, ResponseStruct.Items[i].Maghrib, location)
			prayerDayMap["Isha"] = parseTime(prayerDate, ResponseStruct.Items[i].Isha, location)

			prayerMonthMap[prayerDate] = prayerDayMap

		}
		newMap := make(map[string]map[string]string)
		for key, innerKeyVal := range prayerMonthMap {
			newMap[key] = make(map[string]string)
			for prayer, timeVal := range innerKeyVal {
				newMap[key][prayer] = timeVal.String()
			}

		}
		logger.Info("Uploading date values to redis...")
		for outerKey, innerMap := range newMap {
			innerMapJson, err := json.Marshal(innerMap)
			if err != nil {
				logger.Error("error marshalling inner map: %w", err)
				//TODO check if return is right or should add continue here
				return nil, err
			}
			// upload data to redis from api call, so it can be used in next cycle
			err = client.Set(outerKey, innerMapJson, 0).Err()
			if err != nil {
				logger.Errorf("error uploading outerKey %s and innerMapJson %s, error: %s", outerKey, innerMapJson, err)
				return nil, err
			}
		}
		logger.Info("Uploading successful")
		// logger.Info("Monthly prayer API call is")
		// logger.Info(prayerMonthMap)
		return prayerMonthMap, nil
	}

	monthlyDataRedis := make(map[string]map[string]time.Time)
	var prayerTimesRedis PrayerTimesRedis
	for day := apiDate; day.Month() == apiDate.Month(); day = day.AddDate(0, 0, 1) {
		redisDateKey := day.Format("2006-01-02")

		redisData, err := client.Get(redisDateKey).Result()

		if err != nil {
			logger.Errorf("redis data get request caused error, err: %w", err)
		}

		err = json.Unmarshal([]byte(redisData), &prayerTimesRedis)
		if err != nil {
			logger.Errorf("failed to unmarshal redis data into struct, err: %w", err)
		}
		// logger.Infof("PrayerTimesRedis Struct: %s", prayerTimesRedis.Fajr)

		dailyPrayerTimesMap := make(map[string]time.Time)
		dailyPrayerTimesMap["Fajr"] = parseRedisTimeString(prayerTimesRedis.Fajr, logger)
		dailyPrayerTimesMap["Dhuhr"] = parseRedisTimeString(prayerTimesRedis.Dhuhr, logger)
		dailyPrayerTimesMap["Asr"] = parseRedisTimeString(prayerTimesRedis.Asr, logger)
		dailyPrayerTimesMap["Maghrib"] = parseRedisTimeString(prayerTimesRedis.Maghrib, logger)
		dailyPrayerTimesMap["Isha"] = parseRedisTimeString(prayerTimesRedis.Isha, logger)

		monthlyDataRedis[redisDateKey] = dailyPrayerTimesMap

	}
	// logger.Info("MONTHLY DATA FROM REDIS CACHE IS: ")
	// logger.Info(monthlyDataRedis)
	return monthlyDataRedis, nil
}

func parseRedisTimeString(redisTimeString string, logger *zap.SugaredLogger) time.Time {

	layout := "2006-01-02 15:04:05 -0700 MST"
	parsedTime, err := time.Parse(layout, redisTimeString)
	if err != nil {
		logger.Errorf("error parsing time string, %w", err)
	}
	return parsedTime
}

func parseTime(dateVal string, timeVal string, location string) time.Time {

	loc, err := time.LoadLocation(location)
	if err != nil {

		return time.Time{}
	}
	timeVal = strings.TrimSpace(timeVal)
	timeParsed, err := time.Parse("3:04 pm", timeVal)
	if err != nil {
		return time.Time{}
	}

	dateParsed, err := time.Parse("2006-01-02", dateVal)
	if err != nil {
		return time.Time{}
	}

	finalTime := time.Date(dateParsed.Year(), dateParsed.Month(), dateParsed.Day(), timeParsed.Hour(), timeParsed.Minute(), 0, 0, loc)
	return finalTime
}

// returned struct from GetPrayerTimes function
type PrayerTimes struct {
	PrayerDate string
	Fajr       time.Time
	Dhuhr      time.Time
	Asr        time.Time
	Maghrib    time.Time
	Isha       time.Time
}
type PrayerTimesRedis struct {
	Fajr    string
	Dhuhr   string
	Asr     string
	Maghrib string
	Isha    string
}

// Response represents the entire JSON response structure.
type ResponseStruct struct {
	Title             string       `json:"title"`
	Query             string       `json:"query"`
	For               string       `json:"for"`
	Method            int          `json:"method"`
	PrayerMethodName  string       `json:"prayer_method_name"`
	Daylight          bool         `json:"daylight"`
	Timezone          string       `json:"timezone"`
	MapImage          string       `json:"map_image"`
	SeaLevel          string       `json:"sealevel"`
	TodayWeather      TodayWeather `json:"today_weather"`
	Link              string       `json:"link"`
	QiblaDirection    string       `json:"qibla_direction"`
	Latitude          string       `json:"latitude"`
	Longitude         string       `json:"longitude"`
	Address           string       `json:"address"`
	City              string       `json:"city"`
	State             string       `json:"state"`
	PostalCode        string       `json:"postal_code"`
	Country           string       `json:"country"`
	CountryCode       string       `json:"country_code"`
	Items             []PrayerItem `json:"items"`
	StatusValid       int          `json:"status_valid"`
	StatusCode        int          `json:"status_code"`
	StatusDescription string       `json:"status_description"`
}

// TodayWeather represents the weather information for today.
type TodayWeather struct {
	Pressure    string `json:"pressure"`
	Temperature string `json:"temperature"`
}

// PrayerItem represents the prayer times for a specific date.
type PrayerItem struct {
	DateFor string `json:"date_for"`
	Fajr    string `json:"fajr"`
	Shurooq string `json:"shurooq"`
	Dhuhr   string `json:"dhuhr"`
	Asr     string `json:"asr"`
	Maghrib string `json:"maghrib"`
	Isha    string `json:"isha"`
}

func todayPrayerHandler(c echo.Context, pt map[string]map[string]time.Time, logger *zap.SugaredLogger) error {

	incomingDate := c.Param("dateValue")
	logger.Info(incomingDate)
	// regex check to see if string is valid date format
	pattern := `^\d{4}-\d{2}-\d{2}$`
	re := regexp.MustCompile(pattern)
	if !re.MatchString(incomingDate) {
		return c.JSON(http.StatusBadRequest, fmt.Sprintf("date value from api incorrect: %d", http.StatusBadRequest))
	}

	prayers := pt[incomingDate]

	c.JSON(http.StatusOK, prayers)
	return nil

}

// DO NOT DELETE THIS
// TODO NEED TO ADD BACK WHEN UPLOADING USER DATA TO POSTGRESQL

type UserDataRequestBody struct {
	CurrentPrayerName string `json:"currentPrayerName"`
	CurrentPrayerTime string `json:"currentPrayerTime"`
	LastPrayerName    string `json:"lastPrayerName"`
	LastPrayerTime    string `json:"lastPrayerTime"`
	ProductiveValue   bool   `json:"productiveValue"`
}

func handlePostUserData(c echo.Context, logger *zap.SugaredLogger, db *sql.DB) error {
	// this function parses the incoming data and uploads it to the postgres database

	var incomingData UserDataRequestBody
	if err := c.Bind(&incomingData); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request body for incoming user submission"})
	}
	logger.Infof("User posted Data %s", incomingData)

	err := uploadUserInput(c, logger, db, incomingData)
	if err != nil {
		logger.Errorf("Failed to upload users input! err: %s", err.Error())
	}

	return c.JSON(http.StatusOK, incomingData)
}

func uploadUserInput(c echo.Context, logger *zap.SugaredLogger, db *sql.DB, userData UserDataRequestBody) error {
	// insertSQL := `
	// CREATE TABLE IF NOT EXISTS user_submissions (
	// 	random_primary_key UUID DEFAULT gen_random_uuid() PRIMARY KEY,
	// 	user_id VARCHAR(255),
	// 	productive_val BOOLEAN,
	// 	first_prayer_name VARCHAR(255),
	// 	second_prayer_name VARCHAR(255),
	// 	first_prayer_time TIMESTAMP,
	// 	second_prayer_time TIMESTAMP,
	// 	ingestion_timestamp TIMESTAMP
	// );
	// `

	insertSQL := fmt.Sprintf(`
	INSERT INTO user_submissions (
		user_id, productive_val, first_prayer_name,
		second_prayer_name, first_prayer_time,
		second_prayer_time, ingestion_timestamp
	) VALUES (
		'talha_1', %s, '%s', '%s',
		'%s', '%s',
		'2023-12-18 12:34:56'
	);
	`, strconv.FormatBool(userData.ProductiveValue), userData.LastPrayerName, userData.CurrentPrayerName, userData.LastPrayerTime, userData.CurrentPrayerTime)

	_, err := db.Exec(insertSQL)
	if err != nil {
		logger.Fatalf("Failed to execute database sql statement, err: %w", err)
		return err
	} else {
		// this is how you query data that has been sent to the database
		///
		sql_select := "SELECT * FROM user_submissions LIMIT 1;"

		rows, err := db.Query(sql_select)
		if err != nil {
			logger.Error(err)
			return err
		}
		defer rows.Close() // Don't forget to close the rows when done
		var (
			random_primary_key  string
			user_id             string
			productive_val      string
			first_prayer_name   string
			second_prayer_name  string
			first_prayer_time   string
			second_prayer_time  string
			ingestion_timestamp string
		)

		for rows.Next() {
			err := rows.Scan(&random_primary_key, &user_id, &productive_val, &first_prayer_name, &second_prayer_name, &first_prayer_time, &second_prayer_time, &ingestion_timestamp)
			if err != nil {
				logger.Error(err)
				return err
			}
		}

		// logger.Info(random_primary_key, user_id, productive_val, first_prayer_name, second_prayer_name, first_prayer_time, second_prayer_time, ingestion_timestamp)
		///
		logger.Info("SUCCESSFULLY UPLOADED TO POSTGRES DB!")
		return nil
	}
}

type UserCredentials struct {
	UserEmail    string `json:"userEmail"`
	UserPassword string `json:"userPassword"`
}

func handleCreateUser(c echo.Context, logger *zap.SugaredLogger, db *sql.DB) error {
	// this function parses the incoming user data, calls an encryption on the password then uploads it to the db
	var incomingUserRegistration UserCredentials
	if err := c.Bind(&incomingUserRegistration); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request body for new user registration"})
	}
	hashed_password, err := hashPassword(incomingUserRegistration.UserPassword)
	if err != nil {
		logger.Errorf(fmt.Sprintf("Failed to hash password, err: %s", err.Error()))
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to hash password"})
	}

	incomingUserRegistration.UserPassword = hashed_password

	// get current timestamp in postgres format to register when new user was created, can use to delete values
	currentTime := currentTimeStampPostgres()
	// set email verification for users initially to false when user is created, they need to confirm later
	verifiedEmail := false

	// primary key is users email here to stop multiple registrations of same user email going through
	// createTableSQL := `
	// CREATE TABLE IF NOT EXISTS users (
	// 	user_id VARCHAR(255) PRIMARY KEY,
	// 	password_hash VARCHAR(255),
	// 	creation_timestamp TIMESTAMP
	//  verified_email BOOLEAN
	// );
	// `

	logger.Info(incomingUserRegistration.UserPassword)

	insertSQL := `
    INSERT INTO users (
        user_id, password_hash, creation_timestamp, verified_email
    ) VALUES (
        $1, $2, $3, $4
    )
`

	_, err = db.Exec(insertSQL, incomingUserRegistration.UserEmail, incomingUserRegistration.UserPassword, currentTime, verifiedEmail)
	if err != nil {
		logger.Errorf("Failed to execute database sql statement, err: %w", err)
		return c.JSON(http.StatusAlreadyReported, map[string]string{"error": "Failed to upload user data to server, is the email already in use?"})
	}

	// add send email verification function here before returning registered user?
	// generate random passphrase for email verification confirmation
	verificationCode, err := generateRandomCode(6)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Could not generate random code"})
	}
	logger.Infof("verification code is %s", verificationCode)
	expiryTime := currentTimePlusHourPostgres()
	logger.Infof("expiry time is %s", expiryTime)

	// setup verification table in db to hold verification codes for users to register
	createTableSQL := `
	CREATE TABLE IF NOT EXISTS email_verification_check (
		user_id VARCHAR(255) PRIMARY KEY,
		email_verification_code INT,
		expiry_time TIMESTAMP
	);
	`
	_, err = db.Exec(createTableSQL)
	if err != nil {
		logger.Errorf("Failed to execute database sql statement, err: %w", err)
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to upload user data to server, is the email already in use?"})
	}
	logger.Info("db for user email verification created")

	// insertVerificationCodeSQL := `
	// INSERT INTO email_verification_check (
	// 	user_id, email_verification_code, expiry_time
	// ) VALUES (
	// 	$1, $2, $3
	// );
	// `
	// _, err = db.Exec(insertVerificationCodeSQL, incomingUserRegistration.UserEmail, verificationCode, expiryTime)
	// if err != nil {
	// 	logger.Errorf("Failed to insert email verification code in DB, err %w", err)
	// 	return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to upload email verification code to db"})
	// }

	err = sendEmailVerification(c, verificationCode, logger)
	if err != nil {
		logger.Errorf("Error sending email, err: %s", err.Error())
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Send email verification failed"})
	}
	logger.Info("email should have sent!")

	// email user a random key to their email, (which is also set in the insert statement of a new DB table called verified_email_check)
	// this should have created date and expiration date on it
	// once user registers, if successful, pop up a verify email button, taking user to a new page (or have it change from greyed out)
	// this page takes user to a new form page, asking for email + verification code
	// this page makes a query to verified_email_check db, checking if there is a record where email + verification code exists + login expiry time<current time
	// if record exists,return 200 and display registration verified, else 404 verification not valid/ expired and show message on frontend
	// pop up resend email verification button -> this reruns email verification somehow?
	return nil
}

// TODO continue from here
func handleLogin(c echo.Context, logger *zap.SugaredLogger, db *sql.DB) error {
	// gets hashed pass from db, compares it to users logged in password, then allows auth to continue or stops it
	var loginCredentials UserCredentials

	if err := c.Bind(&loginCredentials); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to bind user login credentials"})
	}

	queryPassword := `
	SELECT password_hash, verified_email 
	FROM users
	WHERE user_id = $1
	LIMIT 1`

	rows, err := db.Query(queryPassword, loginCredentials.UserEmail)
	if err != nil {
		logger.Info("query error")
		logger.Error(err)
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to run db.Query"})
	}

	defer rows.Close() // Don't forget to close the rows when done
	var hashed_password_from_db string
	var verified_email bool

	for rows.Next() {
		err := rows.Scan(&hashed_password_from_db, &verified_email)
		logger.Infof("hashed password from db is: %s and verified email status is : %s", hashed_password_from_db, strconv.FormatBool(verified_email))
		if err != nil {
			logger.Info("rows.Next() error")
			logger.Error(err)
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed in rows.Next() loop"})
		}
	}

	isPassHashed := checkPasswordHash(loginCredentials.UserPassword, hashed_password_from_db)

	if isPassHashed {
		logger.Infof("password is hashed correctly and login details match!")
		if verified_email {
			logger.Info("email is verified as well")
			// returning nil will return status OK by default
			return nil
		} else {
			logger.Info("email is not verified")
			return c.JSON(http.StatusNotAcceptable, map[string]string{"error": "Email is not verified"})
		}
		//TODO return OK if email verified, otherwise return error asking to verify email
	} else {
		logger.Info("password hash is incorrect!")
		return c.JSON(http.StatusForbidden, map[string]string{"error": "User credentials invalid"})
	}

	return nil
}
