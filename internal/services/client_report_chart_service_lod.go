package services

// import (
// 	"CyberusGolangShareLibrary/postgresql_db"
// 	"encoding/json"
// 	"fmt"
// 	"net/http"
// 	"os"
// )

// // Struct to map the expected JSON fields
// // type WapRedirectRequest struct {
// // 	IdPartner    string `json:"id_partner"`
// // 	RefIdPartner string `json:"refid_partner"`
// // 	MediaPartner string `json:"media_partner"`
// // 	NamePartner  string `json:"name_partner"`
// // }

// // Request
// type ClientReportDataRequest struct {
// 	DateString   string `json:"dateString"`
// 	StartSeconds string `json:"startSeconds"`
// 	EndSeconds   string `json:"endSeconds"`
// }

// type DtacRegisterLogCount struct {
// 	DateString string
// 	ShortCode  string
// 	Total      int64
// }

// var dtacRegisterLogCount []DtacRegisterLogCount

// type DtacCancelLogCount struct {
// 	DateString string
// 	ShortCode  string
// 	Total      int64
// }

// var dtacCancelLogCount []DtacCancelLogCount

// type TmvhRegisterLogCount struct {
// 	DateString string
// 	ShortCode  string
// 	Total      int64
// }

// var tmvhRegisterLogCount []TmvhRegisterLogCount

// type TmvhCancelLogCount struct {
// 	DateString string
// 	ShortCode  string
// 	Total      int64
// }

// var tmvhCancelLogCount []TmvhCancelLogCount

// type LogSummary struct {
// 	Date      string `json:"Date"`
// 	ShortCode string `json:"ShortCode"`
// 	Total     int64  `json:"Total"`
// }

// func ClientReportChartService(r *http.Request) map[string][]LogSummary {
// 	var clientReportDataRequest []ClientReportDataRequest
// 	dbConnection := os.Getenv("BN_DB_URL")
// 	// Decode the JSON body into the struct
// 	decoder := json.NewDecoder(r.Body)
// 	err := decoder.Decode(&clientReportDataRequest)
// 	if err != nil {
// 		fmt.Println("Invalid JSON format : " + err.Error())

// 		return nil
// 	}

// 	// Print the received data
// 	// fmt.Println("Received data:")
// 	// for _, d := range clientReportDataRequest {
// 	// 	fmt.Printf("Date: %s, Start: %s, End: %s\n", d.DateString, d.StartSeconds, d.EndSeconds)
// 	// }

// 	// Init database
// 	postgresDB, sqlConfig, err := postgresql_db.PostgreSqlInstance(dbConnection)
// 	if err != nil {
// 		panic(err)
// 	}
// 	// Test connection
// 	err = sqlConfig.Ping()
// 	if err != nil {
// 		fmt.Println(err)
// 	}

// 	for _, reqList := range clientReportDataRequest {
// 		//################################## DTAC REGISTER ##################################
// 		var DtacRegisterResults []struct {
// 			ShortCode string
// 			Total     int64
// 		}
// 		postgresDB.Table("dtac_subscription_logs").
// 			Select("short_code, COUNT(id) AS total").
// 			Where("action = ? AND timestamp BETWEEN ? AND ?", "REGISTER", reqList.StartSeconds, reqList.EndSeconds).
// 			Group("short_code").
// 			Scan(&DtacRegisterResults)
// 		// Inject dateString into each result
// 		if len(DtacRegisterResults) > 0 {
// 			for _, dtacRegister := range DtacRegisterResults {
// 				dtacRegisterLogCount = append(dtacRegisterLogCount, DtacRegisterLogCount{
// 					DateString: reqList.DateString,
// 					ShortCode:  dtacRegister.ShortCode,
// 					Total:      dtacRegister.Total,
// 				})
// 			}
// 		} else {

// 			dtacRegisterLogCount = append(dtacRegisterLogCount, DtacRegisterLogCount{
// 				DateString: reqList.DateString,
// 				ShortCode:  "0",
// 				Total:      0,
// 			})
// 		}

// 		//################################## DTAC CANCEL ##################################
// 		var DtacCancelResults []struct {
// 			ShortCode string
// 			Total     int64
// 		}
// 		postgresDB.Table("dtac_subscription_logs").
// 			Select("short_code, COUNT(id) AS total").
// 			Where("action = ? AND timestamp BETWEEN ? AND ?", "CANCEL", reqList.StartSeconds, reqList.EndSeconds).
// 			Group("short_code").
// 			Scan(&DtacCancelResults)
// 		if len(DtacCancelResults) > 0 {
// 			for _, dtacCancel := range DtacCancelResults {
// 				dtacCancelLogCount = append(dtacCancelLogCount, DtacCancelLogCount{
// 					DateString: reqList.DateString,
// 					ShortCode:  dtacCancel.ShortCode,
// 					Total:      dtacCancel.Total,
// 				})
// 			}
// 		} else {
// 			dtacCancelLogCount = append(dtacCancelLogCount, DtacCancelLogCount{
// 				DateString: reqList.DateString,
// 				ShortCode:  "0",
// 				Total:      0,
// 			})
// 		}
// 		//################################## TMVH REGISTER ##################################
// 		var TmvhRegisterResults []struct {
// 			ShortCode string
// 			Total     int64
// 		}
// 		postgresDB.Table("tmvh_subscription_logs").
// 			Select("short_code, COUNT(id) AS total").
// 			Where("action = ? AND timestamp BETWEEN ? AND ?", "REGISTER", reqList.StartSeconds, reqList.EndSeconds).
// 			Group("short_code").
// 			Scan(&TmvhRegisterResults)
// 			// Inject dateString into each result
// 		if len(TmvhRegisterResults) > 0 {
// 			for _, tmvhRegister := range TmvhRegisterResults {
// 				tmvhRegisterLogCount = append(tmvhRegisterLogCount, TmvhRegisterLogCount{
// 					DateString: reqList.DateString,
// 					ShortCode:  tmvhRegister.ShortCode,
// 					Total:      tmvhRegister.Total,
// 				})
// 			}
// 		} else {
// 			tmvhRegisterLogCount = append(tmvhRegisterLogCount, TmvhRegisterLogCount{
// 				DateString: reqList.DateString,
// 				ShortCode:  "0",
// 				Total:      0,
// 			})
// 		}

// 		//################################## TMVH CANCEL ##################################
// 		var TmvhCancelResults []struct {
// 			ShortCode string
// 			Total     int64
// 		}
// 		postgresDB.Table("tmvh_subscription_logs").
// 			Select("short_code, COUNT(id) AS total").
// 			Where("action = ? AND timestamp BETWEEN ? AND ?", "CANCEL", reqList.StartSeconds, reqList.EndSeconds).
// 			Group("short_code").
// 			Scan(&TmvhCancelResults)
// 			// Inject dateString into each result
// 		if len(TmvhCancelResults) > 0 {
// 			for _, tmvhCancel := range TmvhCancelResults {
// 				tmvhCancelLogCount = append(tmvhCancelLogCount, TmvhCancelLogCount{
// 					DateString: reqList.DateString,
// 					ShortCode:  tmvhCancel.ShortCode,
// 					Total:      tmvhCancel.Total,
// 				})
// 			}
// 		} else {
// 			tmvhCancelLogCount = append(tmvhCancelLogCount, TmvhCancelLogCount{
// 				DateString: reqList.DateString,
// 				ShortCode:  "0",
// 				Total:      0,
// 			})
// 		}

// 	}

// 	output := make(map[string][]LogSummary)
// 	// Print result
// 	for _, res := range tmvhRegisterLogCount {
// 		output["tmvh-register"] = append(output["tmvh-register"], LogSummary{
// 			Date:      res.DateString,
// 			ShortCode: res.ShortCode,
// 			Total:     res.Total,
// 		})
// 		//fmt.Printf("TRUEMOVE REGISTER Date: %s, ShortCode: %s, Total: %d\n", res.DateString, res.ShortCode, res.Total)
// 	}
// 	for _, res := range tmvhCancelLogCount {
// 		output["tmvh-cancel"] = append(output["tmvh-cancel"], LogSummary{
// 			Date:      res.DateString,
// 			ShortCode: res.ShortCode,
// 			Total:     res.Total,
// 		})
// 		//fmt.Printf("TRUEMOVE CANCEL Date: %s, ShortCode: %s, Total: %d\n", res.DateString, res.ShortCode, res.Total)
// 	}

// 	// Print result
// 	for _, res := range dtacRegisterLogCount {
// 		output["dtac-register"] = append(output["dtac-register"], LogSummary{
// 			Date:      res.DateString,
// 			ShortCode: res.ShortCode,
// 			Total:     res.Total,
// 		})
// 		//fmt.Printf("DTAC REGISTER Date: %s, ShortCode: %s, Total: %d\n", res.DateString, res.ShortCode, res.Total)
// 	}
// 	for _, res := range dtacCancelLogCount {
// 		output["dtac-cancel"] = append(output["dtac-cancel"], LogSummary{
// 			Date:      res.DateString,
// 			ShortCode: res.ShortCode,
// 			Total:     res.Total,
// 		})
// 		//fmt.Printf("DTAC CANCEL Date: %s, ShortCode: %s, Total: %d\n", res.DateString, res.ShortCode, res.Total)
// 	}

// 	defer r.Body.Close()

// 	return output
// }
