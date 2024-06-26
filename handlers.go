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
	// "github.com/golang-jwt/jwt"
	jwt "github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
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

// We add jwt.RegisteredClaims as an embedded type, to provide fields like expiry time
type Claims struct {
	UserType string
	User     string
	jwt.RegisteredClaims
}

func todayPrayerHandler(c echo.Context, pt map[string]map[string]time.Time, logger *zap.SugaredLogger, hmacSecret []byte) error {
	logger.Info("accessing todaysPrayerHandler")

	// below is the scode to parse jwt token & get claims from it + error handling
	tokenString := c.Request().Header.Get("Authorization")

	if tokenString == "" {
		logger.Error("token string empty")
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "No valid JWT token"})

	}
	// // initialise new instance of claims
	claims := &Claims{}

	// // Parse the JWT string and store the result in `claims`.
	// // Note that we are passing the key in this method as well. This method will return an error
	// // if the token is invalid (if it has expired according to the expiry time we set on sign in),
	// // or if the signature does not match
	_, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (any, error) {
		return hmacSecret, nil
	})
	if err != nil {
		if err == jwt.ErrSignatureInvalid {
			logger.Error("Invalid jwt signature")
			return c.JSON(http.StatusUnauthorized, map[string]string{"error": "Unauthorized access, invalid jwt signature"})
		}
		logger.Error("Bad request returned")
		logger.Error(err.Error())
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Bad request parsing jwt claims"})
	}
	// note - this is not a security check and is redundant because an altered jwt is blocked at middleware level
	// this is just to show what we can do with claims, this code will never be reached
	if claims.UserType != "user" {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid claims"})
	}

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

func handlePostUserData(c echo.Context, logger *zap.SugaredLogger, db *sql.DB, hmacSecret []byte) error {
	// this function parses the incoming data and uploads it to the postgres database
	logger.Info("hit postUserData")

	// TODO the following code is repeated in getAllStats, functionalise this and test this properly to make sure all is okay:
	// gets user email from authorization header
	tokenString := c.Request().Header.Get("Authorization")
	logger.Infof("token string is ... %s", tokenString)

	if tokenString == "" {
		logger.Error("token string empty")
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "No valid JWT token"})

	}
	// // initialise new instance of claims
	claims := &Claims{}

	// parse jwt claims into claims var using secret
	_, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (any, error) {
		return hmacSecret, nil
	})

	if err != nil {
		if err == jwt.ErrSignatureInvalid {
			logger.Error("Invalid jwt signature")
			return c.JSON(http.StatusUnauthorized, map[string]string{"error": "Unauthorized access, invalid jwt signature"})
		}
		logger.Error("Bad request returned")
		logger.Error(err.Error())
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Bad request parsing jwt claims"})
	}

	if claims.UserType != "user" {
		logger.Error(" user type is not - user")
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid claims"})
	}

	userEmail := claims.User
	logger.Infof("user email is ... %s", userEmail)
	if userEmail == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "User not logged in"})
	}
	// upload incoming user data to db
	var incomingData UserDataRequestBody
	if err := c.Bind(&incomingData); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request body for incoming user submission"})
	}
	logger.Infof("User posted Data %s", incomingData)

	err = uploadUserInput(c, logger, db, incomingData, userEmail)
	if err != nil {
		logger.Errorf("Failed to upload users input! err: %s", err.Error())
	}

	return c.JSON(http.StatusOK, incomingData)
}

func uploadUserInput(c echo.Context, logger *zap.SugaredLogger, db *sql.DB, userData UserDataRequestBody, userEmail string) error {
	currentTime := time.Now()
	currentTimeFormatted := currentTime.Format("2006-01-02 15:04:05")

	insertSQL := fmt.Sprintf(`
	INSERT INTO user_submissions (
		user_id, productive_val, first_prayer_name,
		second_prayer_name, first_prayer_time,
		second_prayer_time, ingestion_timestamp
	) VALUES (
		'%s', %t, '%s', '%s',
		'%s', '%s',
		'%s'
	);
	`, userEmail, userData.ProductiveValue, userData.LastPrayerName, userData.CurrentPrayerName,
		userData.LastPrayerTime, userData.CurrentPrayerTime, currentTimeFormatted)

	_, err := db.Exec(insertSQL)
	if err != nil {
		logger.Errorf("Failed to execute database sql statement, err: %w", err)
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

	expiryTime := currentTimePlusHourPostgres()

	insertVerificationCodeSQL := `
	INSERT INTO email_verification_check (
		user_id, email_verification_code, expiry_time
	) VALUES (
		$1, $2, $3
	);
	`
	_, err = db.Exec(insertVerificationCodeSQL, incomingUserRegistration.UserEmail, verificationCode, expiryTime)
	if err != nil {
		logger.Errorf("Failed to insert email verification code in DB, err %w", err)
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to upload email verification code to db"})
	}

	err = sendEmailVerification(c, verificationCode, incomingUserRegistration.UserEmail, logger)
	if err != nil {
		logger.Errorf("Error sending email, err: %s", err.Error())
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Send email verification failed"})
	}
	logger.Info("email should have sent to user")

	return c.JSON(http.StatusOK, map[string]string{"error": ""})
}

type EmailVerificationDetailsFromFrontend struct {
	UserEmail        string `json:"userEmail"`
	VerificationCode string `json:"verificationCode"`
}

type EmailVerificationDBResults struct {
	UserEmail        string
	VerificationCode int
	ExpiryTime       time.Time
}

func handleUserVerification(c echo.Context, logger *zap.SugaredLogger, db *sql.DB) error {
	var EmailVerificationDetailsFromFrontend EmailVerificationDetailsFromFrontend
	logger.Info("binding data from frontend...")
	if err := c.Bind(&EmailVerificationDetailsFromFrontend); err != nil {
		logger.Errorf("failed to bind data from frontend err", err.Error())
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Invalid request body for email verification check"})
	}
	// converts verification code from frontend to a int if possible before comparing it against DB verification code entry
	verificationCodeFromFrontend, err := strconv.Atoi(EmailVerificationDetailsFromFrontend.VerificationCode)
	if err != nil {
		logger.Error("Verification code from frontend cannot be converted to integer, returning 500")
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid verification code"})
	}

	// the first sql statement checks if a record exists for registered user and the email verification is false first
	CheckRecordExistsQuery := `
	SELECT user_id, email_verification_code, expiry_time
	FROM email_verification_check
	WHERE user_id = $1
	AND email_verification_code = $2
	`
	rows, err := db.Query(CheckRecordExistsQuery, EmailVerificationDetailsFromFrontend.UserEmail, verificationCodeFromFrontend)
	if err != nil {
		logger.Error("Error in quering db for verification Email information")
		logger.Error(err)
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to run query for email verification on DB"})
	}

	defer rows.Close()

	var countReturnedRows int
	var EmailVerificationDBResults EmailVerificationDBResults

	for rows.Next() {
		err := rows.Scan(&EmailVerificationDBResults.UserEmail, &EmailVerificationDBResults.VerificationCode, &EmailVerificationDBResults.ExpiryTime)
		if err != nil {
			logger.Errorf("Error in rows.Scan for parsing rows into EmailVerificationDBResults")
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Error in rows.Scan for parsing rows into EmailVerificationDBResults"})
		}

		countReturnedRows++

	}

	if countReturnedRows > 1 {
		logger.Error("Rows returned more than 1 from DB for email verification check should only be one")
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Rows returned more than 1 from DB for email verification check should only be one"})
	}
	logger.Info("rows not more than 1")
	logger.Infof("Email verification code from db is %v", EmailVerificationDBResults.VerificationCode)
	logger.Infof("Email code from frontend is %s", EmailVerificationDetailsFromFrontend.VerificationCode)

	expiryTimeValid := time.Now().Before(EmailVerificationDBResults.ExpiryTime)

	// if verificationcode from db is 0 then it is because there is no result so it is a default value, so check to see if not 0
	if EmailVerificationDBResults.VerificationCode == verificationCodeFromFrontend && (EmailVerificationDBResults.VerificationCode != 0) && expiryTimeValid {
		// update verification flag in user database

		updateVerificationFlag := `
		UPDATE users
		SET verified_email =  $1
		WHERE user_id = $2
		`

		_, err := db.Exec(updateVerificationFlag, true, EmailVerificationDBResults.UserEmail)

		if err != nil {
			logger.Error("Failed to set email verification flag to true in DB")
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to set email verification flag to true in DB"})
		}
		logger.Info("Updated table successfully!")
		return c.JSON(http.StatusOK, map[string]string{"error": ""})

	}
	return c.JSON(http.StatusBadRequest, map[string]string{"error": "Email verification code recieved is incorrect!"})

}

type ResetVerificationCodeCheck struct {
	user_id        string
	verified_email bool
}

type EmailReverificationDetailsFromFrontend struct {
	UserEmail string `json:"userEmail"`
}

func handleResetUserVerification(c echo.Context, logger *zap.SugaredLogger, db *sql.DB) error {
	// this function allows user to reset the verification code
	// it gets the data from the frontend
	// it does this by first getting the data for the user from the user db and parsing it into a struct
	// it then checks the struct to see if the user exists, if the user doesnt exist it returns the correct response and asks user to create a login first
	// it then checks to see if a user is verified, if user is verified, it returns a response which asks the user to just login, if not it proceeds to the next part
	// we then generate a verification code for the user
	// any rows that match the users ID in the verification email db are then deleted, and a new row is inserted for the user with the new verification code
	// email is then sent again to the user
	// return OK to user, if OK is recieved from frontend redirect user to login page
	// it pulls the data for the user from the users db
	var (
		// UserEmail        string
		EmailReverificationDetailsFromFrontend EmailReverificationDetailsFromFrontend
		ResetVerificationCodeCheck             ResetVerificationCodeCheck
		verificationCode                       int
	)
	logger.Info("executed reset user verification")

	if err := c.Bind(&EmailReverificationDetailsFromFrontend); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Failed to bind user Email from frontend for resetting verification"})
	}
	logger.Info("after binding incoming data")

	queryUser := `
	SELECT user_id, verified_email FROM users
	WHERE user_id = $1;
	`
	// queryRow is used here as it is expected at most 1 row for user
	err := db.QueryRow(queryUser, EmailReverificationDetailsFromFrontend.UserEmail).Scan(&ResetVerificationCodeCheck.user_id, &ResetVerificationCodeCheck.verified_email)
	if err == sql.ErrNoRows {
		// return this when no rows are found i.e. user has not even registered
		return c.JSON(http.StatusNotFound, map[string]string{"error": "User has not registered for an account yet, please create account"})
	} else if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to query db with queryUser var in handleResetUserVerification"})
	}

	if ResetVerificationCodeCheck.verified_email {
		return c.JSON(http.StatusAlreadyReported, map[string]string{"error": "User is already verified, you can login"})
	}
	verificationCode, err = generateRandomCode(6)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to generate verification code in handleResetUserVerification"})
	}

	queryDeleteExistingCodes := `
	DELETE  FROM email_verification_check
	WHERE user_id = $1;`
	_, err = db.Exec(queryDeleteExistingCodes, EmailReverificationDetailsFromFrontend.UserEmail)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to delete rows from email_verification_check table in handleResetUserVerification"})
	}

	expiryTime := currentTimePlusHourPostgres()
	//TODO this is used twice, maybe its better to set this somewhere as a config var and call?
	insertVerificationCodeSQL := `
	INSERT INTO email_verification_check (
		user_id, email_verification_code, expiry_time
	) VALUES (
		$1, $2, $3
	);
	`
	_, err = db.Exec(insertVerificationCodeSQL, ResetVerificationCodeCheck.user_id, verificationCode, expiryTime)
	if err != nil {
		logger.Errorf("Failed to insert email verification code in DB, err %w", err)
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to upload email verification code to db"})
	}

	err = sendEmailVerification(c, verificationCode, EmailReverificationDetailsFromFrontend.UserEmail, logger)
	if err != nil {
		logger.Errorf("Error sending email, err: %s", err.Error())
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Send email verification failed"})
	}

	return c.JSON(http.StatusOK, map[string]string{"error": ""})

}

func handleLogin(c echo.Context, logger *zap.SugaredLogger, db *sql.DB, hmacSecret []byte) error {
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
		// logger.Infof("hashed password from db is: %s and verified email status is : %s", hashed_password_from_db, strconv.FormatBool(verified_email))
		if err != nil {
			logger.Info("rows.Next() error in getting hashed password + verified email flag")
			logger.Error(err)
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed in rows.Next() loop"})
		}
	}

	isPassHashed := checkPasswordHash(loginCredentials.UserPassword, hashed_password_from_db)

	if isPassHashed {
		logger.Infof("password is hashed correctly and login details match!")
		if verified_email {
			logger.Info("email verified and password correct")
			// Create a new token object, specifying signing method and the claims
			// you would like it to contain.

			claims := &Claims{
				UserType: "user",
				User:     loginCredentials.UserEmail,
				RegisteredClaims: jwt.RegisteredClaims{
					// In JWT, the expiry time is expressed as unix milliseconds

					// ExpiresAt: jwt.NewNumericDate(expirationTime),
					Issuer:    "tpm",
					ExpiresAt: jwt.NewNumericDate((time.Now().Add(24 * time.Hour))),
				},
			}

			token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

			// Sign and get the complete encoded token as a string using the secret

			tokenString, err := token.SignedString(hmacSecret)
			if err != nil {
				logger.Errorf("token.SignedString failed, err: ", err.Error())
			}
			cookie := &http.Cookie{
				Name:     "jwt",
				Value:    tokenString,
				Path:     "/",
				HttpOnly: true,
				SameSite: http.SameSiteStrictMode,
				Secure:   true, // set true if using HTTPS
				// look into getting tls certification to have frontend (ingress controller) call backend service securely using https
				Expires: time.Now().Add(24 * time.Hour),
			}

			c.SetCookie(cookie)

			return c.JSON(http.StatusOK, map[string]string{"error": ""})
		} else {
			logger.Info("email is not verified, password is correct")
			return c.JSON(http.StatusNotAcceptable, map[string]string{"error": "Email is not verified"})
		}

	} else {
		logger.Info("password hash is incorrect!")
		return c.JSON(http.StatusForbidden, map[string]string{"error": "User credentials invalid"})
	}

}

func handleGetAllStats(c echo.Context, logger *zap.SugaredLogger, db *sql.DB, hmacSecret []byte) error {
	logger.Info("handleGetAllStats hit!")
	// this function gets all the users stats to be surfaced to the frontend, it also calls functions
	// which will populate empty areas for the data e.g. weekly data gaps

	// get date value for sql query
	// today := time.Now()
	// last_week := time.Now().AddDate(0, 0, -7)
	// logger.Infof("date is %s, last week date is %s", today, last_week)

	// get user login data from jwt token for sql query
	tokenString := c.Request().Header.Get("Authorization")
	logger.Infof("token string is ... %s", tokenString)

	if tokenString == "" {
		logger.Error("token string empty")
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "No valid JWT token"})

	}
	// // initialise new instance of claims
	claims := &Claims{}

	// parse jwt claims into claims var using secret
	_, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (any, error) {
		return hmacSecret, nil
	})

	if err != nil {
		if err == jwt.ErrSignatureInvalid {
			logger.Error("Invalid jwt signature")
			return c.JSON(http.StatusUnauthorized, map[string]string{"error": "Unauthorized access, invalid jwt signature"})
		}
		logger.Error("Bad request returned")
		logger.Error(err.Error())
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Bad request parsing jwt claims"})
	}

	if claims.UserType != "user" {
		logger.Error(" user type is not - user")
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid claims"})
	}

	userEmail := claims.User
	logger.Infof("user email is ... %s", userEmail)
	if userEmail == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "User not logged in"})
	}

	last_week := time.Now().AddDate(0, 0, -7)

	full_week_sql_query := `
	SELECT 
	productive_val,
	first_prayer_name,
	second_prayer_name,
	first_prayer_time,
	second_prayer_time,
	ingestion_timestamp
	FROM user_submissions
	WHERE
	user_id = $1
	and 
	ingestion_timestamp >= $2
	`

	logger.Info("querying db for stats")
	rows, err := db.Query(full_week_sql_query, userEmail, last_week)
	if err != nil {
		logger.Errorf("Rows errored in get stats, err: %w", err)
	}
	defer rows.Close()
	count := 1
	type UserProductivitySubmissions struct {
		productive_val      bool
		first_prayer_name   string
		second_prayer_name  string
		first_prayer_time   time.Time
		second_prayer_time  time.Time
		ingestion_timestamp time.Time
	}

	type SingleRowSubmission map[string]string

	var allUserProductivitySubmissions []SingleRowSubmission

	for rows.Next() {
		var userProductivitySubmission UserProductivitySubmissions
		singleRowSubmission := make(SingleRowSubmission)

		count += 1
		err := rows.Scan(&userProductivitySubmission.productive_val, &userProductivitySubmission.first_prayer_name,
			&userProductivitySubmission.second_prayer_name, &userProductivitySubmission.first_prayer_time,
			&userProductivitySubmission.second_prayer_time, &userProductivitySubmission.ingestion_timestamp)

		if err != nil {
			// Handle the error, perhaps by logging it or returning it.
			logger.Error("failed to scan variables in get all stats db query!")
		}

		// Convert boolean to string
		productiveValString := strconv.FormatBool(userProductivitySubmission.productive_val)
		// Convert time.Time variables to string
		firstPrayerTimeString := userProductivitySubmission.first_prayer_time.Format(time.RFC3339)
		secondPrayerTimeString := userProductivitySubmission.second_prayer_time.Format(time.RFC3339)
		ingestionTimestampString := userProductivitySubmission.ingestion_timestamp.Format(time.RFC3339)

		singleRowSubmission["productive_val"] = productiveValString
		singleRowSubmission["first_prayer_name"] = userProductivitySubmission.first_prayer_name
		singleRowSubmission["first_prayer_time"] = firstPrayerTimeString
		singleRowSubmission["second_prayer_name"] = userProductivitySubmission.second_prayer_name
		singleRowSubmission["second_prayer_time"] = secondPrayerTimeString
		singleRowSubmission["ingestion_timestamp"] = ingestionTimestampString

		allUserProductivitySubmissions = append(allUserProductivitySubmissions, singleRowSubmission)

		// Print the scanned variables
		// logger.Infof("SCANNED VARIABLES:\nproductive_val: %s, first_prayer_name: %s, second_prayer_name: %s, first_prayer_time: %s, second_prayer_time: %s, ingestion_timestamp: %s",
		// 	productiveValString, userProductivitySubmission.first_prayer_name, userProductivitySubmission.second_prayer_name,
		// 	firstPrayerTimeString, secondPrayerTimeString, ingestionTimestampString)
	}

	// logger.Info(allUserProductivitySubmissions)
	logger.Infof("row number %s", strconv.Itoa(count))

	jsonData, err := json.Marshal(allUserProductivitySubmissions)
	if err != nil {
		logger.Error("failed to marshal user submissions to JSON!")
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to marshal data into json in get all stats"})
	}
	logger.Info(jsonData)

	allStats := getAllStats(c, logger, db, userEmail)

	return c.JSON(http.StatusOK, allStats)
}
