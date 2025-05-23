package services

import (
	"CyberusGolangShareLibrary/postgresql_db"
	"cyberus/client-partner/internal/models"
	"fmt"
	"os"

	"encoding/json"
	"net/http"
)

// Request
type ClientUpdateDataReq struct {
	ReqClientID    string `json:"id"`
	ReqNewPassword string `json:"new_password"`
}

func UpdateClientService(r *http.Request) map[string]string {

	dbConnection := os.Getenv("BN_DB_URL")
	var payload map[string]interface{}
	errPayload := json.NewDecoder(r.Body).Decode(&payload)
	if errPayload != nil {
		// Example: print the values
		//fmt.Println("Error decode Json to map[string]interface{} :", errPayload.Error())
		res := map[string]string{
			"code":    "-1",
			"message": "JSON Invalid",
		}
		return res
	}

	jsonData, err := json.MarshalIndent(payload, "", "  ")
	if err != nil {
		//fmt.Println("Error marshalling JSON:", err.Error())
		res := map[string]string{
			"code":    "-2",
			"message": "Null Json",
		}
		return res
	}
	// // Unmarshal JSON into struct
	var clientRequest ClientUpdateDataReq
	err = json.Unmarshal(jsonData, &clientRequest)
	if err != nil {
		//fmt.Println("Error map Json to Struct :" + err.Error())
		//fmt.Println("Error marshalling JSON:", err.Error())
		res := map[string]string{
			"code":    "-3",
			"message": "JSON Not match struct",
		}
		return res
	}

	clientPartnerModel := models.ClientPartner{
		ID:       clientRequest.ReqClientID,
		Password: clientRequest.ReqNewPassword,
	}

	//fmt.Println(clientRequest.ReqClientID)
	//fmt.Println(clientRequest.ReqNewPassword)

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
	// select * from where
	if err := postgresDB.First(&clientPartnerModel, "id = ?", clientRequest.ReqClientID).Error; err != nil {
		res := map[string]string{
			"code":    "-1",
			"message": "client not found",
		}
		return res
	}

	//clientPartnerModel.Password = clientRequest.ReqNewPassword
	if err := postgresDB.Model(&clientPartnerModel).Update("password", clientRequest.ReqNewPassword).Error; err != nil {

		res := map[string]string{
			"code":    "-1",
			"message": "failures",
		}
		return res
	}

	res := map[string]string{
		"code":    "200",
		"message": "success",
	}
	defer r.Body.Close()

	return res
}
