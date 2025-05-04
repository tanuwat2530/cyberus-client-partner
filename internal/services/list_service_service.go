package services

import (
	"CyberusGolangShareLibrary/postgresql_db"
	"encoding/json"
	"strconv"

	"fmt"
	"net/http"

	"cyberus/client-partner/internal/models"
)

// Struct to map the expected JSON fields
type dataRequest struct {
	ClientPartnerID string `json:"client_partner_id"`
}

func ListServiceService(r *http.Request) []map[string]string {
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
	var dataReq dataRequest
	err = json.Unmarshal(jsonData, &dataReq)
	if err != nil {
		//fmt.Println("Error map Json to Struct :" + err.Error())
		fmt.Println("Error marshalling JSON:", err.Error())

	}
	fmt.Println(dataReq.ClientPartnerID)

	// Init database
	dns := "host=localhost user=root password=11111111 dbname=cyberus_db port=5432 sslmode=disable TimeZone=Asia/Bangkok search_path=root@cyberus"

	postgresDB, sqlConfig, err := postgresql_db.PostgreSqlInstance(dns)
	if err != nil {
		panic(err)
	}
	// Test connection
	err = sqlConfig.Ping()
	if err != nil {
		fmt.Println(err)
	}

	var clients []models.ClientService

	fmt.Println("✅ Connected to PostgreSQL with connection pool")
	result := postgresDB.Where("client_partner_id = ?", dataReq.ClientPartnerID).Find(&clients)
	if result.Error != nil {
		fmt.Println("not found")

		return nil
	}

	var clientsRes []map[string]string
	for _, client := range clients {
		m := map[string]string{
			"id":                strconv.Itoa(client.ID),
			"keyword":           client.Keyword,
			"shortcode":         client.Shortcode,
			"telcoid":           client.TelcoID,
			"ads_id":            client.AdsID,
			"client_partner_id": client.ClientPartnerID,
			"wap_aoc_refid":     client.WapAocRefID,
			"wap_aoc_id":        client.WapAocID,
			"wap_aoc_media":     client.WapAocMedia,
			"postback_url":      client.PostbackURL,
			"dn_url":            client.DNURL,
			"postback_counter":  strconv.Itoa(client.PostbackCounter),
		}
		clientsRes = append(clientsRes, m)
	}
	defer r.Body.Close()

	return clientsRes
}
