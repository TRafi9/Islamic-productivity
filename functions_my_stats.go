package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
)

type TrueFalseSqlSingleRowSubmissions map[string]int

type ProductiveValPieChart struct {
	// this struct is marshalled in allStats before being parsed by frontend
	Productive   int
	Unproductive int
}

func pieChartStats(c echo.Context, logger *zap.SugaredLogger, db *sql.DB, userEmail string, timeInterval string) ProductiveValPieChart {
	var timeLength string
	switch strings.ToLower(timeInterval) {
	case "daily":
		timeLength = "CURRENT_DATE"
	case "weekly":
		timeLength = "CURRENT_DATE - INTERVAL '7 days'"
	case "monthly":
		timeLength = "CURRENT_DATE - INTERVAL '31 days'"
	}

	query := fmt.Sprintf(`
	SELECT 
   productive_val,
   ingestion_timestamp
	FROM user_submissions
	WHERE
    user_id = $1
	AND DATE(ingestion_timestamp) >= %s;
	`, timeLength)

	rows, err := db.Query(query, userEmail)
	if err != nil {
		logger.Errorf("Rows errored in get stats, err: %w", err)
	}
	defer rows.Close()
	count := 0

	var productiveValTrue int
	var productiveValFalse int

	for rows.Next() {
		var productiveVal bool
		var ingestion_timestamp time.Time

		count += 1
		err := rows.Scan(&productiveVal, &ingestion_timestamp)
		logger.Infof("result from db for prod val and ingestion timestamp are %s and %s", productiveVal, ingestion_timestamp)

		if err != nil {
			// Handle the error, perhaps by logging it or returning it.
			logger.Error("failed to scan variables in get all stats db query!")
			return ProductiveValPieChart{}
		}

		// Convert boolean to string
		if productiveVal {
			logger.Info("this rows productive val was true")
			productiveValTrue += 1
		} else {
			productiveValFalse += 1
		}

		if err := rows.Err(); err != nil {
			logger.Error("Error while iterating through rows:", err)
			return ProductiveValPieChart{}
		}

	}
	productiveValPieChart := ProductiveValPieChart{
		Productive:   productiveValTrue,
		Unproductive: productiveValFalse,
	}

	return productiveValPieChart
}

func weeklyStats(c echo.Context, logger *zap.SugaredLogger, db *sql.DB, userEmail string) ProductiveValPieChart {

	query := `
	SELECT 
   productive_val,
   ingestion_timestamp
	FROM user_submissions
	WHERE
    user_id = $1
	AND DATE(ingestion_timestamp) >= CURRENT_DATE - INTERVAL '7 days';
	`

	logger.Info("querying db for stats")
	logger.Infof("querying db with the following user email and current day value %s,%s", userEmail)

	rows, err := db.Query(query, userEmail)
	if err != nil {
		logger.Errorf("Rows errored in get stats, err: %w", err)
	}
	defer rows.Close()
	count := 0

	var productiveValTrue int
	var productiveValFalse int

	for rows.Next() {
		var productiveVal bool
		var ingestion_timestamp time.Time

		count += 1
		err := rows.Scan(&productiveVal, &ingestion_timestamp)
		logger.Infof("result from db for prod val and ingestion timestamp are %s and %s", productiveVal, ingestion_timestamp)

		if err != nil {
			// Handle the error, perhaps by logging it or returning it.
			logger.Error("failed to scan variables in get all stats db query!")
			return ProductiveValPieChart{}
		}

		// Convert boolean to string
		if productiveVal {
			logger.Info("this rows productive val was true")
			productiveValTrue += 1
		} else {
			productiveValFalse += 1
		}

		if err := rows.Err(); err != nil {
			logger.Error("Error while iterating through rows:", err)
			return ProductiveValPieChart{}
		}

	}
	productiveValPieChart := ProductiveValPieChart{
		Productive:   productiveValTrue,
		Unproductive: productiveValFalse,
	}

	return productiveValPieChart
}

func getAllStats(c echo.Context, logger *zap.SugaredLogger, db *sql.DB, userEmail string) string {
	dailyStats := pieChartStats(c, logger, db, userEmail, "daily")
	weeklyStats := pieChartStats(c, logger, db, userEmail, "weekly")
	monthlyStats := pieChartStats(c, logger, db, userEmail, "monthly")

	type AllStats struct {
		DailyStats   ProductiveValPieChart
		WeeklyStats  ProductiveValPieChart
		MonthlyStats ProductiveValPieChart
	}
	allStats := AllStats{
		DailyStats:   dailyStats,
		WeeklyStats:  weeklyStats,
		MonthlyStats: monthlyStats,
	}
	logger.Info("all stats non marshalled...")
	logger.Info(allStats)
	json, err := json.Marshal(allStats)
	if err != nil {
		c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to marshal in getAllStats"})
	}
	logger.Info("all stats marshalled and converted to string...")
	logger.Info(string(json))

	return (string(json))
}
