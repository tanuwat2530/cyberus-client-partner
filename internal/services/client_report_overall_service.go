package services

import (
	"CyberusGolangShareLibrary/postgresql_db"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
)

type RequestShortcode struct {
	ListShortcode []ShortCode `json:"list-shortcode"`
	// DateTime      []DateTime  `json:"date-time"`
}

type dataSummary struct {
	//Date      string `json:"Date"`
	ShortCode     string `json:"ShortCode"`
	RegisterTotal int64  `json:"RegisterTotal"`
	CancelTotal   int64  `json:"CancelTotal"`
}

func ClientReportOverallService(r *http.Request) map[string][]dataSummary {
	var output = map[string][]dataSummary{
		"data-summary": {},
	}

	var requestShortcode RequestShortcode
	err := json.NewDecoder(r.Body).Decode(&requestShortcode)
	if err != nil {
		fmt.Println("client_report_overall_service : Invalid JSON #1")
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

	for _, sc := range requestShortcode.ListShortcode {
		registerTotal := int64(0)
		cancelTotal := int64(0)
		var result struct {
			Total int64
		}

		// DTAC REGISTER
		postgresDB.Table("ais_subscription_logs").
			Select("COUNT(id) AS total").
			Where("action = 'REGISTER' AND short_code = ?", sc.ShortCode).
			Scan(&result)

		registerTotal += result.Total

		// DTAC CANCEL
		postgresDB.Table("ais_subscription_logs").
			Select("COUNT(id) AS total").
			Where("action = 'CANCEL' AND short_code = ?", sc.ShortCode).
			Scan(&result)
		cancelTotal += result.Total

		// DTAC REGISTER
		postgresDB.Table("dtac_subscription_logs").
			Select("COUNT(id) AS total").
			Where("action = 'REGISTER' AND short_code = ?", sc.ShortCode).
			Scan(&result)
		registerTotal += result.Total

		// DTAC CANCEL
		postgresDB.Table("dtac_subscription_logs").
			Select("COUNT(id) AS total").
			Where("action = 'CANCEL' AND short_code = ?", sc.ShortCode).
			Scan(&result)
		cancelTotal += result.Total

		// TMVH REGISTER
		postgresDB.Table("tmvh_subscription_logs").
			Select("COUNT(id) AS total").
			Where("action = 'REGISTER' AND short_code = ?", sc.ShortCode).
			Scan(&result)
		registerTotal += result.Total

		// TMVH CANCEL
		postgresDB.Table("tmvh_subscription_logs").
			Select("COUNT(id) AS total").
			Where("action = 'CANCEL' AND short_code = ?", sc.ShortCode).
			Scan(&result)
		cancelTotal += result.Total
		output["data-summary"] = append(output["data-summary"], dataSummary{ShortCode: sc.ShortCode, RegisterTotal: registerTotal, CancelTotal: cancelTotal})

	}

	return output
}
