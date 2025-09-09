package services

import (
	"CyberusGolangShareLibrary/postgresql_db"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
)

type AdminMoDetailReques struct {
	StartTime string `json:"start-time"`
	EndTime   string `json:"end-time"`
	Telco     string `json:"telco"`
}

type AdminMoDetailSummary struct {
	Id            string `json:"id"`
	Action        string `json:"action"`
	Code          string `json:"code"`
	CyberusReturn string `json:"cyberus_return"`
	Description   string `json:"description"`
	Media         string `json:"media"`
	Msisdn        string `json:"msisdn"`
	Operator      string `json:"operator"`
	RefId         string `json:"ref_id"`
	Shortcode     string `json:"shortcode"`
	Timestamp     string `json:"timestamp"`
	Token         string `json:"token"`
	TranRef       string `json:"tran_ref"`
}

func AdminMoDetailService(r *http.Request) []AdminMoDetailSummary {
	var output []AdminMoDetailSummary

	var adminMoDetailReques AdminMoDetailReques
	err := json.NewDecoder(r.Body).Decode(&adminMoDetailReques)
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
	query := postgresDB.Table("").Where("timestamp >= ? AND timestamp <= ?", adminMoDetailReques.StartTime, adminMoDetailReques.EndTime)

	// Select the correct table based on telco
	switch adminMoDetailReques.Telco {
	case "ais":
		query = query.Table("ais_subscription_logs")
	case "dtac":
		query = query.Table("dtac_subscription_logs")
	case "tmvh":
		query = query.Table("tmvh_subscription_logs")
	default:
		fmt.Println("Invalid telco specified:", adminMoDetailReques.Telco)
		return nil
	}

	// Use Find to populate the output slice directly
	if err := query.Find(&output).Error; err != nil {
		fmt.Println("Failed to query data:", err)
		return nil
	}

	return output
}
