package dashboard

import (
	"fmt"
	"net/http"

	"example.com/m/src/db"
	"example.com/m/src/util"
)

type dashboardDataResponseObj struct {
	AllGaugeData []GaugeData `json:"allGaugeData" bson:"AllGaugeData"`
}

func GetGaugeData(w http.ResponseWriter, r *http.Request) {
	// get a Mongo client
	cli := db.GetClient()
	if cli == nil {
		fmt.Println("[GetEmail] Failed to connect to DB")
		util.JsonResponse(w, "Failed to connect to DB", http.StatusBadGateway)
		return
	}

	// get a handle for the AttackLog collection
	logsColl := cli.Database(db.VedikaCorpDatabase).Collection(db.AttackLogCollection)

	// create an object for the response data
	respData := dashboardDataResponseObj{
		AllGaugeData: make([]GaugeData, 0),
	}

	// ------------- CALL ALL FUNCTIONS TO DO GAUGE DATA AGGREGATIONS -------------
	// ----- UserCountsGauge -----
	userCountsData, err := AggregateUserCountsData(logsColl)
	if err != nil {
		fmt.Printf("[GetGaugeData][UserCountsGauge] Failed to aggregate data: %+v", err)
		util.JsonResponse(w, "Failed to get data", http.StatusBadGateway)
		return
	} else {
		respData.AllGaugeData = append(respData.AllGaugeData, userCountsData)
	}

	// ----- TotalTeamPerformance -----
	teamPerfData, err := AggregateTeamTotalPerformanceData(logsColl)
	if err != nil {
		fmt.Printf("[GetGaugeData][TotalTeamPerformance] Failed to aggregate data: %+v", err)
		util.JsonResponse(w, "Failed to get data", http.StatusBadGateway)
		return
	} else {
		respData.AllGaugeData = append(respData.AllGaugeData, teamPerfData)
	}

	// ----- EmailCountsGauge -----
	emailCountsData, err := AggregateEmailCountsData(logsColl)
	if err != nil {
		fmt.Printf("[GetGaugeData][EmailCountsGauge] Failed to aggregate data: %+v", err)
		util.JsonResponse(w, "Failed to get data", http.StatusBadGateway)
		return
	} else {
		respData.AllGaugeData = append(respData.AllGaugeData, emailCountsData)
	}

	util.JsonResponse(w, respData, http.StatusOK)
}
