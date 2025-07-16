package services

import (
	"CyberusGolangShareLibrary/postgresql_db"
	"cyberus/client-partner/internal/models"
	"fmt"
	"net/http"
	"os"
)

func ListShortcodeAllService(r *http.Request) []map[string]string {
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

	var clients []models.ClientService

	fmt.Println("✅ Connected to PostgreSQL with connection pool")
	result := postgresDB.Select("DISTINCT shortcode").Find(&clients)
	if result.Error != nil {
		fmt.Println("not found")
		return nil
	}

	var clientsRes []map[string]string
	for _, client := range clients {
		m := map[string]string{
			//"id":                strconv.Itoa(client.ID),
			//"keyword":           client.Keyword,
			"shortcode": client.Shortcode,
			//"telcoid":           client.TelcoID,
			//"ads_id":            client.AdsID,
			//"client_partner_id": client.ClientPartnerID,
			//"wap_aoc_refid":     client.WapAocRefID,
			//"wap_aoc_id":        client.WapAocID,
			//"wap_aoc_media":     client.WapAocMedia,
			//"postback_url":      client.PostbackURL,
			//"dn_url":            client.DNURL,
			//"postback_counter":  strconv.Itoa(client.PostbackCounter),
		}
		clientsRes = append(clientsRes, m)
	}

	return clientsRes
}
