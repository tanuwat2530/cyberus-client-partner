package services

import (
	"CyberusGolangShareLibrary/postgresql_db"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
)

type AdminDnRequesData struct {
	DateTime []DateTime `json:"date-time"`
}

type AdminDataDnSummary struct {
	DnDate         string `json:"dn-date"`
	StartTimeStamp string `json:"start-timestamp"`
	EndTimeStamp   string `json:"end-timestamp"`
	AisDn          int64  `json:"ais-dn"`
	DtacDn         int64  `json:"dtac-dn"`
	TmvhDn         int64  `json:"tmvh-dn"`
}

func AdminReportDnService(r *http.Request) []AdminDataDnSummary {
	var output []AdminDataDnSummary

	var adminDNRequesData AdminDnRequesData
	err := json.NewDecoder(r.Body).Decode(&adminDNRequesData)
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

	for _, sc := range adminDNRequesData.DateTime {
		var aisDnResult, dtacDnResult, tmvhDnResult int64

		// Query for AIS
		postgresDB.Table("ais_transaction_logs").
			Select("COUNT(id) AS total").
			Where("timestamp >= ? AND timestamp <= ?", sc.StartSeconds, sc.EndSeconds).
			Scan(&aisDnResult)

		// Query for DTAC
		postgresDB.Table("dtac_transaction_logs").
			Select("COUNT(id) AS total").
			Where("timestamp >= ? AND timestamp <= ?", sc.StartSeconds, sc.EndSeconds).
			Scan(&dtacDnResult)

		// Query for TMVH
		postgresDB.Table("tmvh_transaction_logs").
			Select("COUNT(id) AS total").
			Where("timestamp >= ? AND timestamp <= ?", sc.StartSeconds, sc.EndSeconds).
			Scan(&tmvhDnResult)

		// Append the results to the output slice
		output = append(output, AdminDataDnSummary{
			DnDate:         sc.DateString,
			StartTimeStamp: sc.StartSeconds,
			EndTimeStamp:   sc.EndSeconds,
			AisDn:          aisDnResult,
			DtacDn:         dtacDnResult,
			TmvhDn:         tmvhDnResult,
		})
	}

	return output
}
