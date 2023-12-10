package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"regexp"
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
	apiDate = time.Date(apiDate.Year(), apiDate.Month(), 01, 0, 0, 0, 0, apiDate.Location())
	apiDateString := apiDate.Format("02-01-2006")
	redisDateFormat := apiDate.Format("2006-01-02")
	// redisDateFormat := "01-02-2006"
	fmt.Printf("date string: %s \n", apiDateString)

	res, err := client.Get(redisDateFormat).Result()

	if err != nil {
		// error could also be that redis.Nil aka there was no value returned
		logger.Errorf("error with redis get call, continue to get data from  %w", err)

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

		// get todays date
		// todaysDate := time.Now().Format("2006-01-02")

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

		for outerKey, innerMap := range newMap {
			innerMapJson, err := json.Marshal(innerMap)
			if err != nil {
				logger.Error("error marshalling inner map: %w", err)
				//TODO check if return is right or should add continue here
				return nil, err
			}
			err = client.Set(outerKey, innerMapJson, 0).Err()
			if err != nil {
				logger.Errorf("error uploading outerKey %s and innerMapJson %s, error: %s", outerKey, innerMapJson, err)
				return nil, err
			}
		}

		return prayerMonthMap, nil
	}

	monthlyDataRedis := make(map[string]map[string]time.Time)
	var prayerTimesRedis PrayerTimesRedis
	for day := apiDate; day.Month() == apiDate.Month(); day = day.AddDate(0, 0, 1) {
		redisDateKey := day.Format("2006-01-02")

		redisData, err := client.Get(redisDateKey).Result()
		logger.Info("SUCCESSFULLY DOWNLOADED DATA FROM REDIS")
		if err != nil {
			logger.Errorf("redis data get request caused error, err: %w", err)
		}

		// logger.Infof("Redis data for %s day is: %s", redisDateKey, redisData)
		err = json.Unmarshal([]byte(redisData), &prayerTimesRedis)
		if err != nil {
			logger.Errorf("failed to unmarshal redis data into struct, err: %w", err)
		}
		logger.Infof("PrayerTimesRedis Struct: %s", prayerTimesRedis)
		//TODO CONTINUE need to parse the prayerTimesRedis struct from string to Time.time values and then can add it to the dailyPrayerTimesMap
		// From there you can add it to monthlyDataRedis map and serve it as the return of this function
		// dailyPrayerTimesMap := map[string]time.Time{
		// 	"Asr":     prayerTimesRedis.Asr,
		// 	"Dhuhr":   prayerTimesRedis.Dhuhr,
		// 	"Fajr":    prayerTimesRedis.Fajr,
		// 	"Isha":    prayerTimesRedis.Isha,
		// 	"Maghrib": prayerTimesRedis.Maghrib,
		// }
		// monthlyDataRedis[redisDateKey] = dailyPrayerTimesMap
		// logger.Infof("DAILY PRAYER MAP FOR DATE: %s, IS %s", redisDateKey, monthlyDataRedis[redisData])
	}
	return monthlyDataRedis, nil
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
