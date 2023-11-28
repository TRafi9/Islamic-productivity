package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

func GetPrayerTimes(location string) (map[string]map[string]time.Time, error) {

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
	fmt.Printf("date string: %s \n", apiDateString)

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
	return prayerMonthMap, nil

	// create an instance of our struct we want to unmarshal this string into

	// prayerTimes := PrayerTimes{
	// 	PrayerDate: prayerDate,
	// 	Fajr:       parseTime(prayerDate, FajrTime, location),
	// 	Dhuhr:      parseTime(prayerDate, DhuhrTime, location),
	// 	Asr:        parseTime(prayerDate, AsrTime, location),
	// 	Maghrib:    parseTime(prayerDate, MaghribTime, location),
	// 	Isha:       parseTime(prayerDate, IshaTime, location),
	// }

	// // return PrayerTimes, nil
	// return prayerTimes, nil
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
