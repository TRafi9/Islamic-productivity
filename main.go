package main

import (
	"context"
	"net/http"
	"os"
	"time"

	"cloud.google.com/go/bigquery"
	"github.com/go-redis/redis"
	"github.com/labstack/echo"
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

	// user := os.Getenv("USER")
	// password := os.Getenv("PASSWORD")
	// dbName := os.Getenv("DB_NAME")
	// connectionString := fmt.Sprintf("host=localhost user=%s password=%s dbname=%s sslmode=disable", user, password, dbName)

	// logger.Info(connectionString)

	// db, err := sql.Open("postgres", connectionString)

	// if err != nil {
	// 	logger.Fatalf("Failed to open sql connection, err: %w", err)
	// }
	// defer db.Close()
	// // verify connection to db by pinging it
	// err = db.Ping()
	// if err != nil {
	// 	logger.Fatalf("Ping to db failed, fataling out, err: %w", err)
	// }

	//BQ initialisation and uploading functionality
	ctx := context.Background()
	projectID := "starlit-booster-408007"
	// opt := option.WithCredentialsFile("")
	bqClient, err := bigquery.NewClient(ctx, projectID)
	if err != nil {
		logger.Errorf("failed to create bqClient, err: %w", err.Error())
	}

	myDataset := bqClient.Dataset("the_productive_muslim")
	// if err := myDataset.Create(ctx, nil); err != nil {
	// 	logger.Errorf("failed to connect to dataset in BQ, err: %w", err.Error())
	// }

	table := myDataset.Table("user-submissions")

	uploader := table.Inserter()

	type user_submissions struct {
		User_id             string    `bigquery:"user_id"`
		Productive_val      bool      `bigquery:"productive_val"`
		First_prayer_name   string    `bigquery:"first_prayer_name"`
		Second_prayer_name  string    `bigquery:"second_prayer_name"`
		First_prayer_time   time.Time `bigquery:"first_prayer_time"`
		Second_prayer_time  time.Time `bigquery:"second_prayer_time"`
		Ingestion_timestamp time.Time `bigquery:"ingestion_timestamp"`
	}
	// Item implements the ValueSaver interface.
	userSubmissionItems := &user_submissions{
		User_id:             "talha_1",
		Productive_val:      true,
		First_prayer_name:   "Fajr",
		Second_prayer_name:  "Dhuhr",
		First_prayer_time:   time.Date(2023, 12, 16, 15, 4, 5, 0, time.UTC),
		Second_prayer_time:  time.Date(2023, 12, 16, 20, 20, 5, 0, time.UTC),
		Ingestion_timestamp: time.Now(),
	}
	if err := uploader.Put(ctx, userSubmissionItems); err != nil {
		logger.Errorf("error uploading userSubmissionItems, err: %w", err.Error())
	}

	e := echo.New()

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

	api.POST("/userData", handlePostUserData)

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

	// api.GET("/sendUserInput/:value", func(c echo.Context) error {
	// 	return uploadUserInput(c, logger, db)
	// })

	e.Start(":8080")
}
