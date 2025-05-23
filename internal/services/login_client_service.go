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
type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Session  string `json:"session"`
}

func LoginClientService(r *http.Request) map[string]string {
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
	var loginRequest LoginRequest
	err = json.Unmarshal(jsonData, &loginRequest)
	if err != nil {
		//fmt.Println("Error map Json to Struct :" + err.Error())
		fmt.Println("Error marshalling JSON:", err.Error())

	}
	fmt.Println(loginRequest.Username)
	fmt.Println(loginRequest.Password)

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
	result := postgresDB.Where("username = ? and password = ?", loginRequest.Username, loginRequest.Password).First(&clients)
	if result.Error != nil {
		fmt.Println("not found")
		loginRes := map[string]string{
			"code": "0",
		}
		return loginRes
	}

	loginRes := map[string]string{
		"code": "1",
	}
	defer r.Body.Close()

	redis_db.ConnectRedis(redisConnection, "", 0)
	redis_key := loginRequest.Username + ":" + loginRequest.Session

	ttl := 1 * time.Hour // expires in 240 Hour

	timestamp := time.Now().Unix()
	// Set key with TTL
	if err := redis_db.SetWithTTL(redis_key, "Login at : "+strconv.FormatInt(timestamp, 10), ttl); err != nil {
		//write to file if Redis problem or forward request to AIS
		log.Fatalf("SetWithTTL error: %v", err)
	}

	return loginRes
}
