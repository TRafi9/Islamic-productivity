package main

import (
	"database/sql"
	"encoding/json"

	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
)

type TrueFalseSqlSingleRowSubmissions map[string]int

func dailyStats(c echo.Context, logger *zap.SugaredLogger, db *sql.DB, userEmail string) string {
	// returns string of json data
	// returns empty string when error

	query := `
	SELECT 
    COUNT(CASE WHEN productive_val = TRUE THEN 1 END) AS true_count,
    COUNT(CASE WHEN productive_val = FALSE THEN 1 END) AS false_count
	FROM user_submissions
	WHERE
    user_id = $1
    AND DATE(ingestion_timestamp) = CURRENT_DATE;
	`

	logger.Info("querying db for stats")
	rows, err := db.Query(query, userEmail)
	if err != nil {
		logger.Errorf("Rows errored in get stats, err: %w", err)
	}
	defer rows.Close()
	count := 0
	type DailyStatsSQL struct {
		trueCount  int
		falseCount int
	}

	var allUserProductivitySubmissions []TrueFalseSqlSingleRowSubmissions

	for rows.Next() {
		var dailyStatsSQL DailyStatsSQL
		singleRowSubmission := make(TrueFalseSqlSingleRowSubmissions)

		count += 1
		err := rows.Scan(&dailyStatsSQL.trueCount, &dailyStatsSQL.falseCount)

		if err != nil {
			// Handle the error, perhaps by logging it or returning it.
			logger.Error("failed to scan variables in get all stats db query!")
		}

		// Convert boolean to string

		singleRowSubmission["true_count"] = dailyStatsSQL.trueCount
		singleRowSubmission["false_count"] = dailyStatsSQL.falseCount

		allUserProductivitySubmissions = append(allUserProductivitySubmissions, singleRowSubmission)

		// Print the scanned variables
		// logger.Infof("SCANNED VARIABLES:\nproductive_val: %s, first_prayer_name: %s, second_prayer_name: %s, first_prayer_time: %s, second_prayer_time: %s, ingestion_timestamp: %s",
		// 	productiveValString, userProductivitySubmission.first_prayer_name, userProductivitySubmission.second_prayer_name,
		// 	firstPrayerTimeString, secondPrayerTimeString, ingestionTimestampString)
	}
	if err := rows.Err(); err != nil {
		logger.Error("Error while iterating through rows:", err)
		return ""
	}

	jsonData, err := json.Marshal(allUserProductivitySubmissions)
	if err != nil {
		logger.Error("failed to marshal user submissions to JSON!")
		return ""
	}
	return string(jsonData)

}
