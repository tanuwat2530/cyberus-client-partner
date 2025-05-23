package services

import (
	"CyberusGolangShareLibrary/redis_db"
	"encoding/json"
	"os"

	"fmt"
	"net/http"
)

// Struct to map the expected JSON fields
type SessionLoginRequest struct {
	Username string `json:"username"`
	Session  string `json:"session"`
}

func SessionLoginService(r *http.Request) map[string]string {
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
	var sessionLoginRequest SessionLoginRequest
	err = json.Unmarshal(jsonData, &sessionLoginRequest)
	if err != nil {
		//fmt.Println("Error map Json to Struct :" + err.Error())
		fmt.Println("Error marshalling JSON:", err.Error())

	}

	defer r.Body.Close()

	redis_db.ConnectRedis(redisConnection, "", 0)
	session, err := redis_db.GetValue(sessionLoginRequest.Username + ":" + sessionLoginRequest.Session)
	if err != nil {
		res := map[string]string{
			"code":    "0",
			"message": "",
		}
		return res
	}
	res := map[string]string{
		"code":    "1",
		"message": string(session),
	}
	return res
}
