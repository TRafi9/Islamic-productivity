package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

func GetPrayerTimes(location string) (PrayerTimes, error) {

	// Gets prayer times monthly

	apiLoc := "london"
	london, err := time.LoadLocation("Europe/London")
	if err != nil {
		fmt.Println("Error loading location:", err)
		return PrayerTimes{}, err
	}

	apiDate := time.Now().In(london)
	apiDate = time.Date(apiDate.Year(), apiDate.Month(), 01, 0, 0, 0, 0, apiDate.Location())
	// apiDateString := apiDate.Format("02-01-2006")

	prayerTimesURL := fmt.Sprintf("https://muslimsalat.com/%s/monthly/%s/true.json?key==dd00aaed7ee591ead148b3af566d88f1", apiLoc, "01-01-2023")

	response, err := http.Get(prayerTimesURL)
	if err != nil {
		return PrayerTimes{}, fmt.Errorf("failed to get response from API call to prayerTimesURL, err: %w", err)
	}
	//TODO check why we defer this specifically
	defer response.Body.Close()

	resBody, err := io.ReadAll(response.Body)
	if err != nil {
		return PrayerTimes{}, fmt.Errorf("failed to read response body, err: %w", err)
	}
	jsonString := string(resBody)
	var ResponseStruct ResponseStruct

	err = json.Unmarshal([]byte(jsonString), &ResponseStruct)
	if err != nil {
		return PrayerTimes{}, fmt.Errorf("error unmarshalling json into struct, err: %w", err)
	}

	fmt.Println(ResponseStruct.Items[0].Fajr)
	fmt.Println(ResponseStruct.Items[0].DateFor)

	// prayerDate := ResponseStruct.Items[0].DateFor

	// if len(prayerDate) == 9 {
	// 	lastDigit := prayerDate[len(prayerDate)-1]
	// 	prayerDate = prayerDate[:len(prayerDate)-1]
	// 	prayerDate = prayerDate + "0" + string(lastDigit)

	// }

	// get today string
	// todayDay := "01"

	var prayerDate string
	for i, _ := range ResponseStruct.Items {
		prayerDate = ResponseStruct.Items[i].DateFor
		if len(prayerDate) > 8 {
			// this can be or even 2023-1-23 so need to account for both possibilites 2023-11-10
			// check if last two digits contain a hyphen, if so need to parse to add 0 in day value
			fmt.Println(prayerDate[4:5])
			if strings.Contains(prayerDate[len(prayerDate)-2:], "-") {
				lastDigit := prayerDate[len(prayerDate)-1]
				prayerDate = prayerDate[:len(prayerDate)-1]
				prayerDate = prayerDate + "0" + string(lastDigit)
			}
			if !strings.Contains(prayerDate[5:7], "-") {
				fmt.Println(prayerDate)
			}
			if prayerDate[4:5] == "-" {
				fmt.Sprintf("contains -, need to fix, %s", prayerDate)
			}

			// adds 0 into date val e.g. 2023-01-2 turns into 2023-01-02

			fmt.Printf("new prayer date %s \n", prayerDate)
		} else if len(prayerDate) == 8 {
			prayerDate = prayerDate[:5] + "0" + prayerDate[5:7] + "0" + prayerDate[7:]
			fmt.Printf("new prayer date %s \n", prayerDate)

		}

	}

	var (
		FajrTime    = ResponseStruct.Items[0].Fajr
		DhuhrTime   = ResponseStruct.Items[0].Dhuhr
		AsrTime     = ResponseStruct.Items[0].Asr
		MaghribTime = ResponseStruct.Items[0].Maghrib
		IshaTime    = ResponseStruct.Items[0].Isha
	)

	// create an instance of our struct we want to unmarshal this string into

	prayerTimes := PrayerTimes{
		Fajr:    parseTime(prayerDate, FajrTime, location),
		Dhuhr:   parseTime(prayerDate, DhuhrTime, location),
		Asr:     parseTime(prayerDate, AsrTime, location),
		Maghrib: parseTime(prayerDate, MaghribTime, location),
		Isha:    parseTime(prayerDate, IshaTime, location),
	}

	// return PrayerTimes, nil
	return prayerTimes, nil
}

func insertChar(s string, pos int, char rune) string {
	return strings.Join([]string{s[:pos], string(char), s[pos:]}, "")
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
	Fajr    time.Time
	Dhuhr   time.Time
	Asr     time.Time
	Maghrib time.Time
	Isha    time.Time
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
