package main

import (
	"database/sql"
	"fmt"
	"os"

	"github.com/go-co-op/gocron/v2"
	"github.com/go-redis/redis"

	// "github.com/labstack/echo"
	// "github.com/labstack/echo/middleware"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
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

	// logger.Info(connectionString)

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
		AllowOrigins:     []string{"http://tpm-frontend:3000", "http://localhost:3000"},
		AllowMethods:     []string{echo.GET, echo.POST},
		AllowHeaders:     []string{"Authorization", "Content-Type", "Set-Cookie"},
		ExposeHeaders:    []string{"Authorization, Set-Cookie"}, // Add this line
		AllowCredentials: true,
	}))

	// provides protection against cross-site scripting (XSS) attack, content type sniffing,
	// clickjacking, insecure connection and other code injection attacks.
	e.Use(middleware.Secure())
	// pass env var for secret here
	hmacSecret := []byte("RandomSecretStringHere")
	//TODO fix middleware implementation to check if incoming JWT is valid
	// e.Use(echojwt.JWT(hmacSecret))

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
	//TODO PULL SAMPLE SECRET KEY HERE FOR SIGNING JWT

	api.GET("/getPrayerTimes/:dateValue", func(c echo.Context) error {
		return todayPrayerHandler(c, Pt, logger)
	})

	api.POST("/userData", func(c echo.Context) error {
		return handlePostUserData(c, logger, db, hmacSecret)
	})

	api.POST("/createUser", func(c echo.Context) error {
		return handleCreateUser(c, logger, db)
	})

	api.POST("/userVerification", func(c echo.Context) error {
		return handleUserVerification(c, logger, db)
	})

	api.POST("/login", func(c echo.Context) error {
		return handleLogin(c, logger, db, hmacSecret)
	})

	api.POST("/resetUserVerification", func(c echo.Context) error {
		return handleResetUserVerification(c, logger, db)
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

	// updatePt needs to be triggered off a cloud function that is set to run every X days on cloud scheduler

	// run a cron job daily to check if the prayer values for the current date exist in Pt, if not, update Pt values for latest month & add to redis
	s, err := gocron.NewScheduler()
	if err != nil {
		fmt.Printf("Error creating scheduler: %v\n", err)
		return
	}
	// Add a job to the scheduler
	j, err := s.NewJob(
		// runs 1 minute past midnight every day
		gocron.CronJob("01 00 * * *", false),
		gocron.NewTask(func() {
			latestPt, err := prayerTimesCronJob(client, logger, location, Pt)
			if err != nil {
				// Do not update prayer times
				logger.Errorf("Error returned: %s", err.Error())
			} else {
				// Update prayer times
				Pt = latestPt
			}
		}),
	)
	if err != nil {
		fmt.Printf("Error creating job: %v\n", err)
		return
	}

	// Print the job ID
	fmt.Println(j.ID())

	// Start the scheduler
	s.Start()

	// api.GET("/updatePt", func(c echo.Context) error {
	// 	newPt, err := GetPrayerTimes(location, client, logger)
	// 	if err != nil {
	// 		logger.Errorf("Error running updatePt %w", err)
	// 		errorResponse := ErrorResponse{
	// 			Message: "Failed to update prayer times",
	// 			Error:   err.Error(),
	// 		}
	// 		return c.JSON(http.StatusInternalServerError, errorResponse)
	// 	}
	// 	Pt = newPt
	// 	successResponse := ErrorResponse{
	// 		Message: "Successfully updated prayer times from cloud run",
	// 		Error:   "",
	// 	}
	// 	return c.JSON(http.StatusOK, successResponse)

	// })

	e.Start(":8080")
}
