package services

import (
	"CyberusGolangShareLibrary/postgresql_db"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
)

type ShortCode struct {
	ShortCode string `json:"shortcode"`
}

type DateTime struct {
	DateString   string `json:"dateString"`
	StartSeconds string `json:"startSeconds"`
	EndSeconds   string `json:"endSeconds"`
}

type RequestData struct {
	ListShortcode []ShortCode `json:"list-shortcode"`
	DateTime      []DateTime  `json:"date-time"`
}

type LogSummary struct {
	Date      string `json:"Date"`
	ShortCode string `json:"ShortCode"`
	Total     int64  `json:"Total"`
}

func ClientReportChartService(r *http.Request) map[string][]LogSummary {
	var output = map[string][]LogSummary{
		"dtac-register": {},
		"dtac-cancel":   {},
		"tmvh-register": {},
		"tmvh-cancel":   {},
	}

	var reqData RequestData
	err := json.NewDecoder(r.Body).Decode(&reqData)
	if err != nil {
		fmt.Println("client_report_chart_service : Invalid JSON #1")
		return output
	}

	dbConnection := os.Getenv("BN_DB_URL")
	postgresDB, sqlConfig, err := postgresql_db.PostgreSqlInstance(dbConnection)
	if err != nil {
		panic(err)
	}
	if err := sqlConfig.Ping(); err != nil {
		fmt.Println(err)
	}

	defer r.Body.Close()

	for _, sc := range reqData.ListShortcode {
		for _, reqList := range reqData.DateTime {

			var result struct {
				Total int64
			}

			// DTAC REGISTER
			postgresDB.Table("dtac_subscription_logs").
				Select("COUNT(id) AS total").
				Where("action = 'REGISTER' AND timestamp BETWEEN ? AND ? AND short_code = ?", reqList.StartSeconds, reqList.EndSeconds, sc.ShortCode).
				Scan(&result)
			output["dtac-register"] = append(output["dtac-register"], LogSummary{Date: reqList.DateString, ShortCode: sc.ShortCode, Total: result.Total})

			// DTAC CANCEL
			postgresDB.Table("dtac_subscription_logs").
				Select("COUNT(id) AS total").
				Where("action = 'CANCEL' AND timestamp BETWEEN ? AND ? AND short_code = ?", reqList.StartSeconds, reqList.EndSeconds, sc.ShortCode).
				Scan(&result)
			output["dtac-cancel"] = append(output["dtac-cancel"], LogSummary{Date: reqList.DateString, ShortCode: sc.ShortCode, Total: result.Total})

			// TMVH REGISTER
			postgresDB.Table("tmvh_subscription_logs").
				Select("COUNT(id) AS total").
				Where("action = 'REGISTER' AND timestamp BETWEEN ? AND ? AND short_code = ?", reqList.StartSeconds, reqList.EndSeconds, sc.ShortCode).
				Scan(&result)
			output["tmvh-register"] = append(output["tmvh-register"], LogSummary{Date: reqList.DateString, ShortCode: sc.ShortCode, Total: result.Total})

			// TMVH CANCEL
			postgresDB.Table("tmvh_subscription_logs").
				Select("COUNT(id) AS total").
				Where("action = 'CANCEL' AND timestamp BETWEEN ? AND ? AND short_code = ?", reqList.StartSeconds, reqList.EndSeconds, sc.ShortCode).
				Scan(&result)
			output["tmvh-cancel"] = append(output["tmvh-cancel"], LogSummary{Date: reqList.DateString, ShortCode: sc.ShortCode, Total: result.Total})
		}
	}

	return output
}
