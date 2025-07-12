package services

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"CyberusGolangShareLibrary/redis_db"
)

// Struct to map the expected JSON fields
// Struct to map the expected JSON fields
type HistoryRequest struct {
	PartnerId     string `json:"partnerId"`
	PrefixPathern string `json:"pathern"`
}

func RequestHistoryService(r *http.Request) []map[string]string {
	redisConnection := os.Getenv("BN_REDIS_URL")
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
	var historyRequest HistoryRequest
	err = json.Unmarshal(jsonData, &historyRequest)
	if err != nil {
		//fmt.Println("Error map Json to Struct :" + err.Error())
		fmt.Println("Error marshalling JSON:", err.Error())

	}
	fmt.Println(historyRequest.PartnerId)
	fmt.Println(historyRequest.PrefixPathern)

	// Connect to Redis (you already do this)
	redis_db.ConnectRedis(redisConnection, "", 0)

	// Scan for keys
	keys := redis_db.ScanKey(historyRequest.PrefixPathern+":"+historyRequest.PartnerId+"*", 1000)
	fmt.Println(keys)
	// Prepare result slice
	var result []map[string]string

	// Loop through keys and get values
	for _, key := range keys {
		val, err := redis_db.GetValue(key)
		if err != nil {
			fmt.Printf("Key %s does not exist or error: %v\n", key, err)
			continue
		}

		// Add each key-value pair as a map to the result slice
		result = append(result, map[string]string{
			"key":   key,
			"value": val,
		})
	}

	return result
}
