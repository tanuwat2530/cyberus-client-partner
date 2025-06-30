package services

import (
	"CyberusGolangShareLibrary/postgresql_db"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
)

// Struct to map the expected JSON fields
// type WapRedirectRequest struct {
// 	IdPartner    string `json:"id_partner"`
// 	RefIdPartner string `json:"refid_partner"`
// 	MediaPartner string `json:"media_partner"`
// 	NamePartner  string `json:"name_partner"`
// }

type ShortCode struct {
	ShortCode string `json:"shortcode"`
}

type DateTime struct {
	DateString   string `json:"dateString"`
	StartSeconds string `json:"startSeconds"`
	EndSeconds   string `json:"endSeconds"`
}

type RequestData struct {
	ListShortcode []ShortCode `json:"list-shortcode"`
	DateTime      []DateTime  `json:"date-time"`
}

// Request
type ClientReportDataRequest struct {
	ShortCode    string `json:"shortcode"`
	DateString   string `json:"dateString"`
	StartSeconds string `json:"startSeconds"`
	EndSeconds   string `json:"endSeconds"`
}

type DtacRegisterLogCount struct {
	DateString string
	ShortCode  string
	Total      int64
}

var dtacRegisterLogCount []DtacRegisterLogCount

type DtacCancelLogCount struct {
	DateString string
	ShortCode  string
	Total      int64
}

var dtacCancelLogCount []DtacCancelLogCount

type TmvhRegisterLogCount struct {
	DateString string
	ShortCode  string
	Total      int64
}

var tmvhRegisterLogCount []TmvhRegisterLogCount

type TmvhCancelLogCount struct {
	DateString string
	ShortCode  string
	Total      int64
}

var tmvhCancelLogCount []TmvhCancelLogCount

type LogSummary struct {
	Date      string `json:"Date"`
	ShortCode string `json:"ShortCode"`
	Total     int64  `json:"Total"`
}

func ClientReportChartService(r *http.Request) map[string][]LogSummary {
	dbConnection := os.Getenv("BN_DB_URL")

	var reqData RequestData

	err := json.NewDecoder(r.Body).Decode(&reqData)
	if err != nil {
		fmt.Println("Invalid JSON")

	}

	// Print the received data
	// fmt.Println("Received data:")
	// for _, d := range clientReportDataRequest {
	// 	fmt.Printf("Date: %s, Start: %s, End: %s\n", d.DateString, d.StartSeconds, d.EndSeconds)
	// }

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

	for _, sc := range reqData.ListShortcode {
		for _, reqList := range reqData.DateTime {
			//fmt.Printf("%s - %s: %s to %s\n", sc.ShortCode, reqList.DateString, reqList.StartSeconds, reqList.EndSeconds)
			//################################## DTAC REGISTER ##################################
			var DtacRegisterResults []struct {
				ShortCode string
				Total     int64
			}
			postgresDB.Table("dtac_subscription_logs").
				Select("COUNT(id) AS total").
				Where("action = 'REGISTER' AND timestamp BETWEEN ? AND ? AND short_code = ?", reqList.StartSeconds, reqList.EndSeconds, sc.ShortCode).
				Scan(&DtacRegisterResults)
			// Inject dateString into each result
			if len(DtacRegisterResults) > 0 {
				for _, dtacRegister := range DtacRegisterResults {
					if dtacRegister.Total > 0 {
						dtacRegisterLogCount = append(dtacRegisterLogCount, DtacRegisterLogCount{
							DateString: reqList.DateString,
							ShortCode:  sc.ShortCode,
							Total:      dtacRegister.Total,
						})
					} else {
						dtacRegisterLogCount = append(dtacRegisterLogCount, DtacRegisterLogCount{
							DateString: reqList.DateString,
							ShortCode:  "",
							Total:      0,
						})
					}
				}
			} else {
				dtacRegisterLogCount = append(dtacRegisterLogCount, DtacRegisterLogCount{
					DateString: reqList.DateString,
					ShortCode:  "",
					Total:      0,
				})
			}

			//################################## DTAC CANCEL ##################################
			var DtacCancelResults []struct {
				ShortCode string
				Total     int64
			}
			postgresDB.Table("dtac_subscription_logs").
				Select("COUNT(id) AS total").
				Where("action = 'CANCEL' AND timestamp BETWEEN ? AND ? AND short_code = ?", reqList.StartSeconds, reqList.EndSeconds, sc.ShortCode).
				Scan(&DtacCancelResults)
			if len(DtacCancelResults) > 0 {
				for _, dtacCancel := range DtacCancelResults {
					if dtacCancel.Total > 0 {
						dtacCancelLogCount = append(dtacCancelLogCount, DtacCancelLogCount{
							DateString: reqList.DateString,
							ShortCode:  sc.ShortCode,
							Total:      dtacCancel.Total,
						})
					} else {
						dtacCancelLogCount = append(dtacCancelLogCount, DtacCancelLogCount{
							DateString: reqList.DateString,
							ShortCode:  "",
							Total:      0,
						})
					}

				}
			} else {
				dtacCancelLogCount = append(dtacCancelLogCount, DtacCancelLogCount{
					DateString: reqList.DateString,
					ShortCode:  "0",
					Total:      0,
				})
			}
			//################################## TMVH REGISTER ##################################
			var TmvhRegisterResults []struct {
				ShortCode string
				Total     int64
			}
			postgresDB.Table("tmvh_subscription_logs").
				Select("COUNT(id) AS total").
				Where("action = 'REGISTER' AND timestamp BETWEEN ? AND ? AND short_code = ?", reqList.StartSeconds, reqList.EndSeconds, sc.ShortCode).
				Scan(&TmvhRegisterResults)
				// Inject dateString into each result
			if len(TmvhRegisterResults) > 0 {
				for _, tmvhRegister := range TmvhRegisterResults {
					if tmvhRegister.Total > 0 {
						tmvhRegisterLogCount = append(tmvhRegisterLogCount, TmvhRegisterLogCount{
							DateString: reqList.DateString,
							ShortCode:  sc.ShortCode,
							Total:      tmvhRegister.Total,
						})
					} else {
						tmvhRegisterLogCount = append(tmvhRegisterLogCount, TmvhRegisterLogCount{
							DateString: reqList.DateString,
							ShortCode:  "",
							Total:      0,
						})
					}
				}
			} else {
				tmvhRegisterLogCount = append(tmvhRegisterLogCount, TmvhRegisterLogCount{
					DateString: reqList.DateString,
					ShortCode:  "",
					Total:      0,
				})
			}

			//################################## TMVH CANCEL ##################################
			var TmvhCancelResults []struct {
				ShortCode string
				Total     int64
			}
			postgresDB.Table("tmvh_subscription_logs").
				Select("COUNT(id) AS total").
				Where("action = 'CANCEL' AND timestamp BETWEEN ? AND ? AND short_code = ?", reqList.StartSeconds, reqList.EndSeconds, sc.ShortCode).
				Scan(&TmvhCancelResults)
				// Inject dateString into each result
			if len(TmvhCancelResults) > 0 {
				for _, tmvhCancel := range TmvhCancelResults {
					if tmvhCancel.Total > 0 {
						tmvhCancelLogCount = append(tmvhCancelLogCount, TmvhCancelLogCount{
							DateString: reqList.DateString,
							ShortCode:  sc.ShortCode,
							Total:      tmvhCancel.Total,
						})
					} else {
						tmvhCancelLogCount = append(tmvhCancelLogCount, TmvhCancelLogCount{
							DateString: reqList.DateString,
							ShortCode:  "",
							Total:      0,
						})
					}
				}
			} else {
				tmvhCancelLogCount = append(tmvhCancelLogCount, TmvhCancelLogCount{
					DateString: reqList.DateString,
					ShortCode:  "",
					Total:      0,
				})
			}

		}
	}
	output := make(map[string][]LogSummary)
	// Print result
	for _, res := range tmvhRegisterLogCount {
		output["tmvh-register"] = append(output["tmvh-register"], LogSummary{
			Date:      res.DateString,
			ShortCode: res.ShortCode,
			Total:     res.Total,
		})
		//fmt.Printf("TRUEMOVE REGISTER Date: %s, ShortCode: %s, Total: %d\n", res.DateString, res.ShortCode, res.Total)
	}
	for _, res := range tmvhCancelLogCount {
		output["tmvh-cancel"] = append(output["tmvh-cancel"], LogSummary{
			Date:      res.DateString,
			ShortCode: res.ShortCode,
			Total:     res.Total,
		})
		//fmt.Printf("TRUEMOVE CANCEL Date: %s, ShortCode: %s, Total: %d\n", res.DateString, res.ShortCode, res.Total)
	}

	// Print result
	for _, res := range dtacRegisterLogCount {
		output["dtac-register"] = append(output["dtac-register"], LogSummary{
			Date:      res.DateString,
			ShortCode: res.ShortCode,
			Total:     res.Total,
		})
		//fmt.Printf("DTAC REGISTER Date: %s, ShortCode: %s, Total: %d\n", res.DateString, res.ShortCode, res.Total)
	}
	for _, res := range dtacCancelLogCount {
		output["dtac-cancel"] = append(output["dtac-cancel"], LogSummary{
			Date:      res.DateString,
			ShortCode: res.ShortCode,
			Total:     res.Total,
		})
		//fmt.Printf("DTAC CANCEL Date: %s, ShortCode: %s, Total: %d\n", res.DateString, res.ShortCode, res.Total)
	}

	defer r.Body.Close()

	return output
}
