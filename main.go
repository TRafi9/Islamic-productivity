package main

import (
	"database/sql"
	"fmt"
	"net/http"
	"os"

	"github.com/go-redis/redis"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	_ "github.com/lib/pq"
	"go.uber.org/zap"
)

// ngrok http --domain=living-sacred-skunk.ngrok-free.app 8080
func readFile(filepath string) (string, error) {
	content, err := os.ReadFile(filepath)
	if err != nil {
		return "", err
	}
	contentString := string(content)
	return contentString, nil
}

func main() {
	// setup logger
	z, _ := zap.NewProduction()
	logger := z.Sugar()

	// read in password for redis connection
	pass, err := readFile("./pass.txt")
	if err != nil {
		// do something here
		return
	}
	// setup redis client
	client := redis.NewClient(&redis.Options{
		Addr:     "redis-13336.c304.europe-west1-2.gce.cloud.redislabs.com:13336",
		Password: pass,
		DB:       0,
	})

	// setup connection to postgresql db, need to run proxy first

	// when running via docker, need to mount these vars as env vars
	user := os.Getenv("USER")
	password := os.Getenv("PASSWORD")
	dbName := os.Getenv("DB_NAME")
	connectionString := fmt.Sprintf("host=localhost user=%s password=%s dbname=%s sslmode=disable", user, password, dbName)

	logger.Info(connectionString)

	db, err := sql.Open("postgres", connectionString)

	if err != nil {
		logger.Fatalf("Failed to open sql connection, err: %w", err)
	}
	defer db.Close()
	// verify connection to db by pinging it
	err = db.Ping()
	if err != nil {
		logger.Fatalf("Ping to db failed, fataling out, err: %w", err)
	}
	logger.Infof("Successfully connected to Postgre instance")

	e := echo.New()
	// cors not needed when running on docker containers that are on same network because of docker-compose file? - check TODO
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"http://tpm-frontend:3000"},
		AllowMethods: []string{echo.GET, echo.POST},
	}))
	//TODO make getPrayerTimes return the whole month
	// use a cron job to run get prayertimes
	// use the same cron job to add the prayer times to a redis in memory database for that month
	// that way the data will persist if the application goes down
	// add v1 GET api to make a call, given a date, to recieve all prayer times for that day, from the redis server

	location := "Europe/London"

	Pt, err := GetPrayerTimes(location, client, logger)
	if err != nil {
		logger.Errorf("error executing GetPrayerTimes, err %w", err)
	}

	//TODO add panic and recover if it fails to upload to memory

	api := e.Group("/api/v1")

	api.GET("/getPrayerTimes/:dateValue", func(c echo.Context) error {
		return todayPrayerHandler(c, Pt, logger)
	})

	api.POST("/userData", func(c echo.Context) error {
		return handlePostUserData(c, logger, db)
	})

	//TODO CONTINUE FROM HERE
	// also look at serialisation of prayerData so you dont have to convert and revert between time.Time and string values, can store serialized strings in redis that are
	// time.time values instead?

	// make Pt a global variable, then use a cloud function to trigger the handler function underneath
	// handler function should just rerun getPrayerTimes, which will run for the new day and push the month data to redis
	// therefore you want the cloud function to trigger just after midnight on the first of a new month to populate redis data
	// use mutex to lock Pt while you update it, and then open it up once mutex is done

	// you can also use an infinite for-loop instead that will sleep daily and trigger Pt to run if it is the first of the month
	// but this is a copout version (although very viable and efficient heuheh)

	//TODO abstract away Error messages into structs and call them instead
	type ErrorResponse struct {
		Message string `json:"message"`
		Error   string `json:"error"`
	}
	// updatePt needs to be triggered off a cloud function that is set to run every X days on cloud scheduler
	api.GET("/updatePt", func(c echo.Context) error {
		newPt, err := GetPrayerTimes(location, client, logger)
		if err != nil {
			logger.Errorf("Error running updatePt %w", err)
			errorResponse := ErrorResponse{
				Message: "Failed to update prayer times",
				Error:   err.Error(),
			}
			return c.JSON(http.StatusInternalServerError, errorResponse)
		}
		Pt = newPt
		successResponse := ErrorResponse{
			Message: "Successfully updated prayer times from cloud run",
			Error:   "",
		}
		return c.JSON(http.StatusOK, successResponse)

	})

	e.Start(":8080")
}
