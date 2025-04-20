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
type ServiceUpdateDataReq struct {
	ID              int    `json:"id"`
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

func UpdateServiceService(r *http.Request) map[string]string {

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
	var serviceUpdateRequest ServiceUpdateDataReq
	err = json.Unmarshal(jsonData, &serviceUpdateRequest)
	if err != nil {
		//fmt.Println("Error map Json to Struct :" + err.Error())
		//fmt.Println("Error marshalling JSON:", err.Error())
		res := map[string]string{
			"code":    "-3",
			"message": "JSON Not match struct",
		}
		return res
	}

	clientServiceModel := models.ClientService{
		Keyword:         serviceUpdateRequest.Keyword,
		Shortcode:       serviceUpdateRequest.Shortcode,
		TelcoID:         serviceUpdateRequest.TelcoID,
		AdsID:           serviceUpdateRequest.AdsID,
		ClientPartnerID: serviceUpdateRequest.ClientPartnerID,
		WapAocRefID:     serviceUpdateRequest.WapAocRefID,
		WapAocID:        serviceUpdateRequest.WapAocID,
		WapAocMedia:     serviceUpdateRequest.WapAocMedia,
		PostbackURL:     serviceUpdateRequest.PostbackURL,
		DNURL:           serviceUpdateRequest.DNURL,
		PostbackCounter: serviceUpdateRequest.PostbackCounter,
	}

	//fmt.Println(clientRequest.ReqClientID)
	//fmt.Println(clientRequest.ReqNewPassword)

	// select * from where
	if err := db.First(&clientServiceModel, "id = ?", serviceUpdateRequest.ID).Error; err != nil {
		res := map[string]string{
			"code":    "-1",
			"message": "client not found",
		}
		return res
	}

	//Update field
	if err := db.Model(&clientServiceModel).Update("keyword", serviceUpdateRequest.Keyword).Error; err != nil {
		//if err := db.Save(&clientServiceModel).Error; err != nil {
		res := map[string]string{
			"code":    "-1",
			"message": "keyword failures",
		}
		return res
	}

	if err := db.Model(&clientServiceModel).Update("shortcode", serviceUpdateRequest.Shortcode).Error; err != nil {
		//if err := db.Save(&clientServiceModel).Error; err != nil {
		res := map[string]string{
			"code":    "-1",
			"message": "shortcode failures",
		}
		return res
	}

	if err := db.Model(&clientServiceModel).Update("telcoid", serviceUpdateRequest.TelcoID).Error; err != nil {
		//if err := db.Save(&clientServiceModel).Error; err != nil {
		res := map[string]string{
			"code":    "-1",
			"message": "telcoid failures",
		}
		return res
	}

	if err := db.Model(&clientServiceModel).Update("ads_id", serviceUpdateRequest.AdsID).Error; err != nil {
		//if err := db.Save(&clientServiceModel).Error; err != nil {
		res := map[string]string{
			"code":    "-1",
			"message": "ads_id failures",
		}
		return res
	}

	if err := db.Model(&clientServiceModel).Update("wap_aoc_refid", serviceUpdateRequest.WapAocRefID).Error; err != nil {
		//if err := db.Save(&clientServiceModel).Error; err != nil {
		res := map[string]string{
			"code":    "-1",
			"message": "wap_aoc_refid failures",
		}
		return res
	}

	if err := db.Model(&clientServiceModel).Update("wap_aoc_id", serviceUpdateRequest.WapAocID).Error; err != nil {
		//if err := db.Save(&clientServiceModel).Error; err != nil {
		res := map[string]string{
			"code":    "-1",
			"message": "wap_aoc_id failures",
		}
		return res
	}

	if err := db.Model(&clientServiceModel).Update("wap_aoc_media", serviceUpdateRequest.WapAocMedia).Error; err != nil {
		//if err := db.Save(&clientServiceModel).Error; err != nil {
		res := map[string]string{
			"code":    "-1",
			"message": "wap_aoc_media failures",
		}
		return res
	}

	if err := db.Model(&clientServiceModel).Update("postback_url", serviceUpdateRequest.PostbackURL).Error; err != nil {
		//if err := db.Save(&clientServiceModel).Error; err != nil {
		res := map[string]string{
			"code":    "-1",
			"message": "postback_url failures",
		}
		return res
	}

	if err := db.Model(&clientServiceModel).Update("dn_url", serviceUpdateRequest.DNURL).Error; err != nil {
		//if err := db.Save(&clientServiceModel).Error; err != nil {
		res := map[string]string{
			"code":    "-1",
			"message": "dn_url failures",
		}
		return res
	}

	if err := db.Model(&clientServiceModel).Update("postback_counter", serviceUpdateRequest.PostbackCounter).Error; err != nil {
		//if err := db.Save(&clientServiceModel).Error; err != nil {
		res := map[string]string{
			"code":    "-1",
			"message": "postback_counter failures",
		}
		return res
	}

	if err := db.Model(&clientServiceModel).Update("client_partner_id", serviceUpdateRequest.ClientPartnerID).Error; err != nil {
		//if err := db.Save(&clientServiceModel).Error; err != nil {
		res := map[string]string{
			"code":    "-1",
			"message": "client_partner_id failures",
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
