package services

import (
	"crypto/md5"
	"log"
	"time"

	"encoding/hex"
	"encoding/json"
	"fmt"
	"net/http"

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
type ClientPartnerDataRequest struct {
	ReqUsername string `json:"username"`
	ReqPassword string `json:"password"`
}

// Table name on database
type ClientPartner struct {
	ID       string `gorm:"primaryKey"`
	Username string `gorm:"column:username"`
	Password string `gorm:"column:password"`
}

// Generate ID: first 8 chars of MD5(timestamp)
func generateShortMD5ID() string {
	timestamp := fmt.Sprintf("%d", time.Now().UnixNano())
	hash := md5.Sum([]byte(timestamp))
	return hex.EncodeToString(hash[:])[:8]
}

func AddClientService(r *http.Request) map[string]string {

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
			"code":    "-1",
			"message": "Null Json",
		}
		return res
	}
	// // Unmarshal JSON into struct
	var clientRequest ClientPartnerDataRequest
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

	clientPartnerInsert := ClientPartner{
		ID:       generateShortMD5ID(),
		Username: clientRequest.ReqUsername,
		Password: clientRequest.ReqPassword,
	}

	if errInsertDB := db.Create(&clientPartnerInsert).Error; errInsertDB != nil {
		fmt.Println("ERROR INSERT : " + errInsertDB.Error())
		res := map[string]string{
			"code":    "-1",
			"message": "failures",
		}
		return res
	}

	//redis_db.Set("aaa", "AAA", 300)
	res := map[string]string{
		"code":    "200",
		"message": "success",
	}

	// Read the request body
	// body, err := ioutil.ReadAll(r.Body)
	// if err != nil {
	// 	//http.Error(w, "Failed to read body", http.StatusBadRequest)
	// 	return res
	// }
	defer r.Body.Close()

	// Unmarshal JSON into struct
	// var requestData WapRedirectRequest
	// err = json.Unmarshal(body, &requestData)
	// if err != nil {
	// 	//http.Error(w, "Invalid JSON", http.StatusBadRequest)
	// 	return res
	// }

	// // Print the data to the console
	// fmt.Println("##### Received WAP Redirect Data #####")
	// fmt.Println("IDPartner : " + requestData.IdPartner)
	// fmt.Println("RefIDPartner : " + requestData.RefIdPartner)
	// fmt.Println("MediaPartner  : " + requestData.MediaPartner)
	// fmt.Println("NamePartner  : " + requestData.NamePartner)

	// Respond to client
	//w.WriteHeader(http.StatusOK)
	//w.Write([]byte("WAP Redirect received successfully"))

	return res
}
