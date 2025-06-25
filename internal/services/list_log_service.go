package services

import (
	"fmt"
	"net/http"
	"os"

	"CyberusGolangShareLibrary/redis_db"
)

// Struct to map the expected JSON fields

func ListLogService(r *http.Request) []map[string]string {
	redisConnection := os.Getenv("BN_REDIS_URL")

	// Connect to Redis (you already do this)
	redis_db.ConnectRedis(redisConnection, "", 0)

	// Scan for keys
	keys := redis_db.ScanKey("*", 100)

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
