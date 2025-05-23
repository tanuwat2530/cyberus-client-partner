package services

import (
	"CyberusGolangShareLibrary/postgresql_db"
	"os"

	"fmt"
	"net/http"

	"cyberus/client-partner/internal/models"
)

// Struct to map the expected JSON fields

func ListClientService(r *http.Request) []map[string]string {
	dbConnection := os.Getenv("BN_DB_URL")
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
