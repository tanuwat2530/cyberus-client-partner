package services

import (
	"CyberusGolangShareLibrary/postgresql_db"
	"cyberus/client-partner/internal/models"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
)

// Struct to map the expected JSON fields
type DataReq struct {
	ClientPartnerID string `json:"client_partner_id"`
}

func ListShortcodeClientService(r *http.Request) []map[string]string {
	dbConnection := os.Getenv("BN_DB_URL")
	var payload map[string]interface{}
	errPayload := json.NewDecoder(r.Body).Decode(&payload)
	if errPayload != nil {
		// Example: print the values
		fmt.Println("Error decode Json to map[string]interface{} :", errPayload.Error())
	}

	jsonData, err := json.MarshalIndent(payload, "", "  ")
	if err != nil {
		fmt.Println("Error marshalling JSON:", err.Error())

	}

	// // Unmarshal JSON into struct
	var dataReq DataReq
	err = json.Unmarshal(jsonData, &dataReq)
	if err != nil {
		//fmt.Println("Error map Json to Struct :" + err.Error())
		fmt.Println("Error marshalling JSON:", err.Error())

	}
	fmt.Println(dataReq.ClientPartnerID)

	// Init database

	postgresDB, sqlConfig, err := postgresql_db.PostgreSqlInstance(dbConnection)
	if err != nil {
		panic(err)
	}
	// Test connection
	err = sqlConfig.Ping()
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(dataReq.ClientPartnerID)

	var clients []models.ClientService

	fmt.Println("✅ Connected to PostgreSQL with connection pool")
	result := postgresDB.Select("DISTINCT shortcode").Where("client_partner_id = ?", dataReq.ClientPartnerID).Group("shortcode").Find(&clients)
	if result.Error != nil {
		fmt.Println("not found")
		return nil
	}

	var clientsRes []map[string]string
	for _, client := range clients {
		m := map[string]string{
			//"id":                strconv.Itoa(client.ID),
			//"keyword":           client.Keyword,
			"shortcode": client.Shortcode,
			//"telcoid":           client.TelcoID,
			//"ads_id":            client.AdsID,
			//"client_partner_id": client.ClientPartnerID,
			//"wap_aoc_refid":     client.WapAocRefID,
			//"wap_aoc_id":        client.WapAocID,
			//"wap_aoc_media":     client.WapAocMedia,
			//"postback_url":      client.PostbackURL,
			//"dn_url":            client.DNURL,
			//"postback_counter":  strconv.Itoa(client.PostbackCounter),
		}
		clientsRes = append(clientsRes, m)
	}

	return clientsRes
}
