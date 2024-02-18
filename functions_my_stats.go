package main

import (
	"database/sql"
	"encoding/json"
	"time"

	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
)

type TrueFalseSqlSingleRowSubmissions map[string]int

type ProductiveValPieChart struct {
	// this struct can be marshalled into a json value to return json string to frontend
	Productive   int
	Unproductive int
}

func dailyStats(c echo.Context, logger *zap.SugaredLogger, db *sql.DB, userEmail string) string {
	// returns string of json data
	// returns empty string when error
	currentDay := time.Now()
	currentDayFormatted := currentDay.Format("2006-01-02")

	query := `
	SELECT 
   productive_val,
   ingestion_timestamp
	FROM user_submissions
	WHERE
    user_id = $1
	AND DATE(ingestion_timestamp) = CURRENT_DATE;
	`

	logger.Info("querying db for stats")
	logger.Infof("querying db with the following user email and current day value %s,%s", userEmail, currentDayFormatted)
	// 2024-02-18 00:00:00 +0000 UTC is currentDay
	rows, err := db.Query(query, userEmail)
	if err != nil {
		logger.Errorf("Rows errored in get stats, err: %w", err)
	}
	defer rows.Close()
	count := 0
	// type DailyStatsSQL struct {
	// 	trueCount  int
	// 	falseCount int
	// }

	// var allUserProductivitySubmissions []TrueFalseSqlSingleRowSubmissions
	var productiveValTrue int
	var productiveValFalse int

	for rows.Next() {
		var productiveVal bool
		var ingestion_timestamp time.Time

		// var dailyStatsSQL DailyStatsSQL
		// singleRowSubmission := make(TrueFalseSqlSingleRowSubmissions)

		count += 1
		err := rows.Scan(&productiveVal, &ingestion_timestamp)
		logger.Infof("result from db for prod val and ingestion timestamp are %s and %s", productiveVal, ingestion_timestamp)

		if err != nil {
			// Handle the error, perhaps by logging it or returning it.
			logger.Error("failed to scan variables in get all stats db query!")
		}

		// Convert boolean to string
		if productiveVal {
			logger.Info("this rows productive val was true")
			productiveValTrue += 1
		} else {
			productiveValFalse += 1
		}
		// singleRowSubmission["true_count"] = dailyStatsSQL.trueCount
		// singleRowSubmission["false_count"] = dailyStatsSQL.falseCount

		// allUserProductivitySubmissions = append(allUserProductivitySubmissions, singleRowSubmission)

		// Print the scanned variables
		// 	logger.Infof("SCANNED VARIABLES:\ntrue count: %s, false count: %s, ",
		// 		dailyStatsSQL.trueCount, dailyStatsSQL.falseCount)
		// }
		if err := rows.Err(); err != nil {
			logger.Error("Error while iterating through rows:", err)
			return ""
		}

		// logger.Infof("total productive val for true is %s, and false is %s", productiveValTrue, productiveValFalse)
	}
	productiveValPieChart := ProductiveValPieChart{
		Productive:   productiveValTrue,
		Unproductive: productiveValFalse,
	}
	json, err := json.Marshal(productiveValPieChart)
	if err != nil {
		logger.Error("failed to marshal user submissions to JSON!")
		return ""
	}
	return string(json)

	//return fmt.Sprintf("total productive val for true is %v, and false is %v", productiveValTrue, productiveValFalse)

	// jsonData, err := json.Marshal(allUserProductivitySubmissions)
	// if err != nil {
	// 	logger.Error("failed to marshal user submissions to JSON!")
	// 	return ""
	// }
	// return string(jsonData)

}
