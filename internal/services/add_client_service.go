package services

import (
	"CyberusGolangShareLibrary/postgresql_db"
	"crypto/md5"
	"os"
	"time"

	"encoding/hex"
	"encoding/json"
	"fmt"
	"net/http"

	"cyberus/client-partner/internal/models"
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

// Generate ID: first 8 chars of MD5(timestamp)
func generateShortMD5ID() string {
	timestamp := fmt.Sprintf("%d", time.Now().UnixNano())
	hash := md5.Sum([]byte(timestamp))
	return hex.EncodeToString(hash[:])[:8]
}

func AddClientService(r *http.Request) map[string]string {
	dbConnection := os.Getenv("BN_DB_URL")
	// config database pool
	// dsn := "host=localhost user=root password=11111111 dbname=cyberus_db port=5432 sslmode=disable TimeZone=Asia/Bangkok search_path=root@cyberus"
	// db, errDatabase := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	// if errDatabase != nil {
	// 	log.Fatal("Failed to connect to database:", errDatabase)
	// }
	// sqlDB, err := db.DB()
	// if err != nil {
	// 	log.Fatal("Failed to get generic database object:", err)
	// }
	// // Set connection pool settings
	// sqlDB.SetMaxOpenConns(5)                    // Maximum number of open connections
	// sqlDB.SetMaxIdleConns(1)                    // Maximum number of idle connections
	// sqlDB.SetConnMaxLifetime(180 * time.Second) // Connection max lifetime

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
	partnerID := generateShortMD5ID()
	clientPartnerInsert := models.ClientPartner{
		ID:       partnerID,
		Username: clientRequest.ReqUsername,
		Password: clientRequest.ReqPassword,
	}
	fmt.Println(clientPartnerInsert.ID)
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

	// Auto migrate (create table if not exists)
	postgresDB.AutoMigrate(&models.ClientPartner{})

	fmt.Println("✅ Connected to PostgreSQL with connection pool")
	result := postgresDB.Create(&clientPartnerInsert)
	if result.Error != nil {
		fmt.Println("ERROR INSERT : ")
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
