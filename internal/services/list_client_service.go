package services

import (
	"CyberusGolangShareLibrary/postgresql_db"

	"fmt"
	"net/http"

	"cyberus/client-partner/internal/models"
)

// Struct to map the expected JSON fields

func ListClientService(r *http.Request) []map[string]string {

	// Init database
	dns := "host=localhost user=root password=11111111 dbname=cyberus_db port=5432 sslmode=disable TimeZone=Asia/Bangkok search_path=root@cyberus"

	postgresDB, sqlConfig, err := postgresql_db.PostgreSqlInstance(dns)
	if err != nil {
		panic(err)
	}
	// Test connection
	err = sqlConfig.Ping()
	if err != nil {
		fmt.Println(err)
	}

	var clients []models.ClientPartner

	fmt.Println("✅ Connected to PostgreSQL with connection pool")
	result := postgresDB.Find(&clients)
	if result.Error != nil {
		fmt.Println("not found")

		return nil
	}

	var clientsRes []map[string]string
	for _, client := range clients {
		m := map[string]string{
			"id":       client.ID,
			"username": client.Username,
			//"password": client.Password,
		}
		clientsRes = append(clientsRes, m)
	}
	defer r.Body.Close()

	return clientsRes
}
