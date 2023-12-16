package main

import (
	"context"
	"net/http"
	"os"

	"cloud.google.com/go/bigquery"
	"github.com/go-redis/redis"
	"github.com/labstack/echo"
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
	//BQ initialisation and uploading functionality
	ctx := context.Background()
	projectID := "starlit-booster-408007"
	bqClient, err := bigquery.NewClient(ctx, projectID)
	if err != nil {
		logger.Errorf("failed to create bqClient, err: %w", err.Error())
	}

	myDataset := bqClient.Dataset("my_dataset")
	if err := myDataset.Create(ctx, nil); err != nil {
		logger.Errorf("failed to connect to dataset in BQ, err: %w", err.Error())
	}

	table := myDataset.Table("my_table")

	uploader := table.Inserter()

	type user_submissions struct {
		user_id             string
		productive_val      bool
		first_prayer_name   string
		second_prayer_name  string
		first_prayer_time   string
		second_prayer_time  string
		ingestion_timestamp string
	}
	// Item implements the ValueSaver interface.
	userSubmissionItems := &user_submissions{
		user_id:             "talha_1",
		productive_val:      true,
		first_prayer_name:   "Fajr",
		second_prayer_name:  "Dhuhr",
		first_prayer_time:   "2023-12-16T15:04:05Z",
		second_prayer_time:  "2023-12-16T20:20:05Z",
		ingestion_timestamp: "2023-12-16T20:21:00Z",
	}
	if err := uploader.Put(ctx, userSubmissionItems); err != nil {
		// TODO: Handle error.
	}

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
