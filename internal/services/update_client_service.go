package services

import (
	"cyberus/client-partner/internal/models"
	"log"
	"time"

	"encoding/json"
	"net/http"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// Request
type ClientUpdateDataReq struct {
	ReqClientID    string `json:"id"`
	ReqNewPassword string `json:"new_password"`
}

func UpdateClientService(r *http.Request) map[string]string {

	// config database pool
	dsn := "host=localhost user=root password=11111111 dbname=cyberus_db port=5432 sslmode=disable TimeZone=Asia/Bangkok search_path=root@cyberus"
	db, errDatabase := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if errDatabase != nil {
		log.Fatal("Failed to connect to database:", errDatabase)
	}
	sqlDB, err := db.DB()
	if err != nil {
		log.Fatal("Failed to get generic database object:", err)
	}
	// Set connection pool settings
	sqlDB.SetMaxOpenConns(5)                    // Maximum number of open connections
	sqlDB.SetMaxIdleConns(1)                    // Maximum number of idle connections
	sqlDB.SetConnMaxLifetime(180 * time.Second) // Connection max lifetime

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

	// select * from where
	if err := db.First(&clientPartnerModel, "id = ?", clientRequest.ReqClientID).Error; err != nil {
		res := map[string]string{
			"code":    "-1",
			"message": "client not found",
		}
		return res
	}

	//clientPartnerModel.Password = clientRequest.ReqNewPassword
	if err := db.Model(&clientPartnerModel).Update("password", clientRequest.ReqNewPassword).Error; err != nil {

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
