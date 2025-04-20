package services

import (
	"log"
	"time"

	"encoding/json"
	"fmt"
	"net/http"

	"cyberus/client-partner/internal/models"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// Struct to map the expected JSON fields
// type WapRedirectRequest struct {
// 	IdPartner    string `json:"id_partner"`
// 	RefIdPartner string `json:"refid_partner"`
// 	MediaPartner string `json:"media_partner"`
// 	NamePartner  string `json:"name_partner"`
// }

// Request
type ClientServiceDataRequest struct {
	Keyword         string `json:"keyword"`
	Shortcode       string `json:"shortcode"`
	TelcoID         string `json:"telcoid"`
	AdsID           string `json:"ads_id"`
	ClientPartnerID string `json:"client_partner_id"`
	WapAocRefID     string `json:"wap_aoc_refid"`
	WapAocID        string `json:"wap_aoc_id"`
	WapAocMedia     string `json:"wap_aoc_media"`
	PostbackURL     string `json:"postback_url"`
	DNURL           string `json:"dn_url"`
	PostbackCounter *int   `json:"postback_counter"`
}

func AddServiceService(r *http.Request) map[string]string {

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
	var serviceDataRequest ClientServiceDataRequest
	err = json.Unmarshal(jsonData, &serviceDataRequest)
	if err != nil {
		//fmt.Println("Error map Json to Struct :" + err.Error())
		//fmt.Println("Error marshalling JSON:", err.Error())
		res := map[string]string{
			"code":    "-3",
			"message": "JSON Not match struct",
		}
		return res
	}

	clientServiceInsert := models.ClientService{
		Keyword:         serviceDataRequest.Keyword,
		Shortcode:       serviceDataRequest.Shortcode,
		TelcoID:         serviceDataRequest.TelcoID,
		AdsID:           serviceDataRequest.AdsID,
		ClientPartnerID: serviceDataRequest.ClientPartnerID,
		WapAocRefID:     serviceDataRequest.WapAocRefID,
		WapAocID:        serviceDataRequest.WapAocID,
		WapAocMedia:     serviceDataRequest.WapAocMedia,
		PostbackURL:     serviceDataRequest.PostbackURL,
		DNURL:           serviceDataRequest.DNURL,
		PostbackCounter: serviceDataRequest.PostbackCounter,
	}

	if errInsertDB := db.Create(&clientServiceInsert).Error; errInsertDB != nil {
		fmt.Println("ERROR INSERT : " + errInsertDB.Error())
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
