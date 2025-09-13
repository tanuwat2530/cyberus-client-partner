package services

import (
	"CyberusGolangShareLibrary/postgresql_db"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
)

type AdminRequesData struct {
	DateTime []DateTime `json:"date-time"`
}

type AdminDataMoSummary struct {
	MoDate         string `json:"mo-date"`
	StartTimeStamp string `json:"start-timestamp"`
	EndTimeStamp   string `json:"end-timestamp"`
	AisMo          int64  `json:"ais-mo"`
	DtacMo         int64  `json:"dtac-mo"`
	TmvhMo         int64  `json:"tmvh-mo"`
}

func AdminReportMoService(r *http.Request) []AdminDataMoSummary {
	var output []AdminDataMoSummary

	var adminRequesData AdminRequesData
	err := json.NewDecoder(r.Body).Decode(&adminRequesData)
	if err != nil {
		fmt.Println("admin_report_all_mo_service : Invalid JSON #1:", err)
		return nil
	}

	dbConnection := os.Getenv("BN_DB_URL")
	postgresDB, sqlConfig, err := postgresql_db.PostgreSqlInstance(dbConnection)
	if err != nil {
		fmt.Println("Database connection error:", err)
		return nil
	}
	if err := sqlConfig.Ping(); err != nil {
		fmt.Println("Database ping error:", err)
	}

	defer r.Body.Close()

	for _, sc := range adminRequesData.DateTime {
		var aisMoResult, dtacMoResult, tmvhMoResult int64

		// Query for AIS
		postgresDB.Table("ais_subscription_logs").
			Select("COUNT(id) AS total").
			Where("timestamp >= ? AND timestamp <= ?", sc.StartSeconds, sc.EndSeconds).
			Order("timestamp DESC"). // Add the Order method here
			Scan(&aisMoResult)

		// Query for DTAC
		postgresDB.Table("dtac_subscription_logs").
			Select("COUNT(id) AS total").
			Where("timestamp >= ? AND timestamp <= ?", sc.StartSeconds, sc.EndSeconds).
			Order("timestamp DESC"). // Add the Order method here
			Scan(&dtacMoResult)

		// Query for TMVH
		postgresDB.Table("tmvh_subscription_logs").
			Select("COUNT(id) AS total").
			Where("timestamp >= ? AND timestamp <= ?", sc.StartSeconds, sc.EndSeconds).
			Order("timestamp DESC"). // Add the Order method here
			Scan(&tmvhMoResult)

		// Append the results to the output slice
		output = append(output, AdminDataMoSummary{
			MoDate:         sc.DateString,
			StartTimeStamp: sc.StartSeconds,
			EndTimeStamp:   sc.EndSeconds,
			AisMo:          aisMoResult,
			DtacMo:         dtacMoResult,
			TmvhMo:         tmvhMoResult,
		})
	}

	return output
}
