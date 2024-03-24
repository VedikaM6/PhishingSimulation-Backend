package util

import (
	"encoding/json"
	"fmt"
	"net/http"
)

const (
	URLParameterEmailId        = "emailId"
	URLParameterAttackId       = "attackId"
	URLQueryParameterStartTime = "startTime"
	URLQueryParameterEndTime   = "endTime"
)

func JsonResponse(w http.ResponseWriter, data interface{}, statCode int) {
	// marshal the data
	dataMarsh, err := json.Marshal(data)
	if err != nil {
		fmt.Printf("[JsonResponse] Failed to marshal data: %+v\n", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.WriteHeader(statCode)

	_, err = w.Write(dataMarsh)
	if err != nil {
		fmt.Printf("[JsonResponse] Failed to write response data: %+v\n", err)
		return
	}
}
