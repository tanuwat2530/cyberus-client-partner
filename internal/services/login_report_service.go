package services

import (
	"CyberusGolangShareLibrary/postgresql_db"
	"CyberusGolangShareLibrary/redis_db"
	"encoding/json"
	"log"
	"os"
	"strconv"
	"time"

	"fmt"
	"net/http"

	"cyberus/client-partner/internal/models"
)

// Struct to map the expected JSON fields
type LoginReportRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Session  string `json:"session"`
}

func LoginReportService(r *http.Request) []map[string]string {
	redisConnection := os.Getenv("BN_REDIS_URL")
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
	var loginReportRequest LoginReportRequest
	err = json.Unmarshal(jsonData, &loginReportRequest)
	if err != nil {
		//fmt.Println("Error map Json to Struct :" + err.Error())
		fmt.Println("Error marshalling JSON : #1", err.Error())

	}
	//fmt.Println(loginRequest.Username)
	//fmt.Println(loginRequest.Password)

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

	var clients []models.ClientPartner

	fmt.Println("✅ Connected to PostgreSQL with connection pool")
	result := postgresDB.Where("username = ? and password = ?", loginReportRequest.Username, loginReportRequest.Password).First(&clients)
	if result.Error != nil {
		fmt.Println("not found")
		var loginRes []map[string]string

		m := map[string]string{
			"code":       "0",
			"partner_id": "0",
		}
		loginRes = append(loginRes, m)
		return loginRes
	}
	var loginRes []map[string]string
	for _, client := range clients {
		m := map[string]string{
			//"id":                strconv.Itoa(client.ID),
			//"keyword":           client.Keyword,
			"code":       "1",
			"partner_id": client.ID,
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
		loginRes = append(loginRes, m)
	}

	defer r.Body.Close()

	redis_db.ConnectRedis(redisConnection, "", 0)
	redis_key := loginReportRequest.Username + ":" + loginReportRequest.Session

	ttl := 1 * time.Hour // expires in 240 Hour

	timestamp := time.Now().Unix()
	// Set key with TTL
	if err := redis_db.SetWithTTL(redis_key, "Login at : "+strconv.FormatInt(timestamp, 10), ttl); err != nil {
		//write to file if Redis problem or forward request to AIS
		log.Fatalf("SetWithTTL error: %v", err)
	}

	return loginRes
}
