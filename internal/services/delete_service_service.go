package services

import (
	"CyberusGolangShareLibrary/postgresql_db"
	"CyberusGolangShareLibrary/redis_db"
	"cyberus/client-partner/internal/models"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
)

// Request
type ServiceDeleteDataReq struct {
	ID              string `json:"id"`
	ClientPartnerID string `json:"client_partner_id"`
	Media           string `json:"media"`
}

func DeleteServiceService(r *http.Request) map[string]string {
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
	var serviceDeleteDataReq ServiceDeleteDataReq
	err = json.Unmarshal(jsonData, &serviceDeleteDataReq)
	if err != nil {
		//fmt.Println("Error map Json to Struct :" + err.Error())
		fmt.Println("Error marshalling JSON:", err.Error())

	}
	fmt.Println(serviceDeleteDataReq.ID)
	fmt.Println(serviceDeleteDataReq.ClientPartnerID)

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

	result := postgresDB.Where("id = ? AND client_partner_id = ?", serviceDeleteDataReq.ID, serviceDeleteDataReq.ClientPartnerID).
		Delete(&models.ClientService{})

	if result.Error != nil {
		fmt.Println(result.Error)
		res := map[string]string{
			"code":    "-1",
			"message": result.Error.Error(),
		}
		return res
	}

	if result.RowsAffected == 0 {
		fmt.Println("no matching record found to delete")
		res := map[string]string{
			"code":    "-2",
			"message": "NOT FOUND",
		}
		return res
	}

	redis_db.ConnectRedis(redisConnection, "", 0)

	// Set key with TTL
	redis_db.DelValue("SERVICE:" + serviceDeleteDataReq.ClientPartnerID + ":" + serviceDeleteDataReq.Media)

	res := map[string]string{
		"code":    "200",
		"message": "success",
	}
	defer r.Body.Close()

	return res
}
