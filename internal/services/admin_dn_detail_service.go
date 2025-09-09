package services

import (
	"CyberusGolangShareLibrary/postgresql_db"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
)

type AdminDnDetailReques struct {
	StartTime string `json:"start-time"`
	EndTime   string `json:"end-time"`
	Telco     string `json:"telco"`
}

type AdminDnDetailSummary struct {
	Id            string `json:"id"`
	Code          string `json:"code"`
	Description   string `json:"description"`
	Msisdn        string `json:"msisdn"`
	Operator      string `json:"operator"`
	Shortcode     string `json:"short_code"`
	Timestamp     string `json:"timestamp"`
	TranRef       string `json:"tran_ref"`
	CyberusReturn string `json:"cyberus_return"`
}

func AdminDnDetailService(r *http.Request) []AdminDnDetailSummary {
	var output []AdminDnDetailSummary

	var adminDnDetailReques AdminDnDetailReques
	err := json.NewDecoder(r.Body).Decode(&adminDnDetailReques)
	if err != nil {
		fmt.Println("admin_mo_detail_service : Invalid JSON #1:", err)
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

	// Base query with timestamp conditions
	query := postgresDB.Table("").Where("timestamp >= ? AND timestamp <= ?", adminDnDetailReques.StartTime, adminDnDetailReques.EndTime)

	// Select the correct table based on telco
	switch adminDnDetailReques.Telco {
	case "ais":
		query = query.Table("ais_transaction_logs")
	case "dtac":
		query = query.Table("dtac_transaction_logs")
	case "tmvh":
		query = query.Table("tmvh_transaction_logs")
	default:
		fmt.Println("Invalid telco specified:", adminDnDetailReques.Telco)
		return nil
	}

	// Use Find to populate the output slice directly
	if err := query.Find(&output).Error; err != nil {
		fmt.Println("Failed to query data DN :", err)
		return nil
	}

	return output
}
