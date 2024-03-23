package attacks

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"example.com/m/src/db"
	"example.com/m/src/util"
	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// This function is called when we want to check for any scheduled attacks and trigger them (if they are due).
// It should be invoked routinely.
func TriggerPendingAttacks(w http.ResponseWriter, r *http.Request) {
	cli := db.GetClient()
	if cli == nil {
		fmt.Println("[TriggerPendingAttacks] Failed to connect to DB")
		util.JsonResponse(w, "Failed to connect to DB", http.StatusBadGateway)
		return
	}

	// get handles for some collections
	custDb := cli.Database(db.VedikaCorpDatabase)
	pendingAttacksColl := custDb.Collection(db.PendingAttacksCollection)
	attackEmailsColl := custDb.Collection(db.EmailsCollection)
	attackLogColl := custDb.Collection(db.AttackLogCollection)

	// set the query filter to match all pending attacks with a TriggerTime in the past
	filter := bson.D{
		{Key: "TriggerTime", Value: bson.D{
			{Key: "$lte", Value: time.Now().UTC()},
		}},
	}

	// get all pending attacks in the collection that are due to be triggered
	ctx := context.TODO()
	cur, err := pendingAttacksColl.Find(ctx, filter)
	if err != nil {
		fmt.Printf("[TriggerPendingAttacks] Failed to get emails from DB: %+v\n", err)
		util.JsonResponse(w, "Failed to trigger attacks", http.StatusBadGateway)
		return
	}

	defer cur.Close(ctx)

	// decode all the results into a slice
	pendingAttacks := make([]PendingAttackObj, 0)
	err = cur.All(ctx, &pendingAttacks)
	if err != nil {
		fmt.Printf("[TriggerPendingAttacks] Failed to decode results: %+v\n", err)
		util.JsonResponse(w, "Failed to trigger attacks", http.StatusBadGateway)
		return
	}

	// check if there are any attacks to trigger
	if len(pendingAttacks) == 0 {
		util.JsonResponse(w, "There are no attacks scheduled to be triggered.", http.StatusOK)
		return
	} else {
		util.JsonResponse(w, fmt.Sprintf("Triggering %d attacks!", len(pendingAttacks)), http.StatusAccepted)
	}

	// Trigger this process in a goroutine because it may take some time to handle all the attacks
	go func() {
		// loop through all pending attacks
		for _, pendAttack := range pendingAttacks {
			// call a helper function to execute the attack
			err := executeAttack(attackEmailsColl, attackLogColl, pendingAttacksColl, pendAttack)
			if err != nil {
				fmt.Printf("[TriggerPendingAttack] Failed to execute this attack: %+v | %+v\n", pendAttack, err)
			}
		}
	}()
}

// This function is called when we want to get a list of attacks executed in the past.
func ListPreviousAttacks(w http.ResponseWriter, r *http.Request) {
	// get the start and end times to search the history from the URL parameters
	startTimeStr := r.URL.Query().Get(util.URLQueryParameterStartTime)
	endTimeStr := r.URL.Query().Get(util.URLQueryParameterEndTime)

	// check if the strings are non-empty
	if startTimeStr == "" || endTimeStr == "" {
		// One of the times is missing
		fmt.Printf("[ListPreviousAttacks] The start or end time is missing: %s | %s\n", startTimeStr, endTimeStr)
		util.JsonResponse(w, "URL query parameters are missing", http.StatusBadRequest)
		return
	}

	// convert the time strings to time.Time values
	startTime, err := time.Parse(time.RFC3339, startTimeStr)
	if err != nil {
		fmt.Printf("[ListPreviousAttacks] The start time is invalid: %s | %+v\n", startTimeStr, err)
		util.JsonResponse(w, "Failed to parse startTime", http.StatusBadRequest)
		return
	}
	endTime, err := time.Parse(time.RFC3339, endTimeStr)
	if err != nil {
		fmt.Printf("[ListPreviousAttacks] The end time is invalid: %s | %+v\n", endTimeStr, err)
		util.JsonResponse(w, "Failed to parse endTime", http.StatusBadRequest)
		return
	}

	// validate times
	if !startTime.Before(endTime) {
		// The start time is on or after the end time.
		fmt.Printf("[ListPreviousAttacks] The start time must be before the end time: %s | %s\n", startTimeStr, endTimeStr)
		util.JsonResponse(w, "The start time must be before the end time.", http.StatusBadRequest)
		return
	}

	cli := db.GetClient()
	if cli == nil {
		fmt.Println("[ListPreviousAttacks] Failed to connect to DB")
		util.JsonResponse(w, "Failed to connect to DB", http.StatusBadGateway)
		return
	}

	// get a handle for the AttackLog collection
	attackLogColl := cli.Database(db.VedikaCorpDatabase).Collection(db.AttackLogCollection)

	// set the query filter to match all attacks sent between startTime and endTime
	filter := bson.D{
		{Key: "TriggerTime", Value: bson.D{
			{Key: "$gte", Value: startTime},
			{Key: "$lte", Value: endTime},
		}},
	}

	// submit the query
	ctx := context.TODO()
	cur, err := attackLogColl.Find(ctx, filter)
	if err != nil {
		fmt.Printf("[ListPreviousAttacks] Failed to get previous attacks: %+v\n", err)
		util.JsonResponse(w, "Failed to get previous attacks", http.StatusBadGateway)
		return
	}

	defer cur.Close(ctx)

	// decode all the results into a slice
	allLogs := make([]AttackLogObj, 0)
	err = cur.All(ctx, &allLogs)
	if err != nil {
		fmt.Printf("[ListPreviousAttacks] Failed to decode results: %+v\n", err)
		util.JsonResponse(w, "Failed to get attack history", http.StatusBadGateway)
		return
	}

	fmt.Printf("[ListPreviousAttacks][DEBUG] Succesfully retrieved %d logs.\n", len(allLogs))

	// prepare the response data and return it
	respData := make(map[string][]AttackLogObj)
	respData["previousAttacks"] = allLogs
	util.JsonResponse(w, respData, http.StatusOK)
}

// This function is called when we want to save an attack to be executed at a later date.
func ScheduleFutureAttack(w http.ResponseWriter, r *http.Request) {
	// decode the request data
	var newAttack PendingAttackObj
	err := json.NewDecoder(r.Body).Decode(&newAttack)
	if err != nil {
		fmt.Printf("[ScheduleFutureAttack] Failed to decode request data: %+v\n", err)
		util.JsonResponse(w, "Request data is invalid", http.StatusBadRequest)
		return
	}

	// validate the attack details
	validationErr := ""
	if newAttack.EmailId.IsZero() {
		// No email was specified for this attack
		validationErr = "You must specify an email for this attack."
	} else if newAttack.TargetRecipient.Name == "" || newAttack.TargetRecipient.Address == "" {
		// Recipient info is missing for this attack
		validationErr = "Recipient info is missing."
	} else if newAttack.TargetUserId.IsZero() {
		// No target user was specified for this attack
		validationErr = "The targeted user is missing."
	} else if newAttack.TriggerTime.Before(time.Now()) {
		// The trigger time is in the past.
		validationErr = "The trigger time cannot be in the past."
	}

	if validationErr != "" {
		util.JsonResponse(w, validationErr, http.StatusBadRequest)
		return
	}

	// connect to the database
	cli := db.GetClient()
	if cli == nil {
		fmt.Println("[ScheduleFutureAttack] Failed to connect to DB")
		util.JsonResponse(w, "Failed to connect to DB", http.StatusBadGateway)
		return
	}

	// get a handle for the PendingAttacks collection
	pendingAttacksColl := cli.Database(db.VedikaCorpDatabase).Collection(db.PendingAttacksCollection)

	res, err := newAttack.LogAttack(pendingAttacksColl)
	if err != nil {
		fmt.Printf("[ScheduleFutureAttack] Failed to insert new attack: %+v\n", err)
		util.JsonResponse(w, "Failed to schedule attack", http.StatusBadGateway)
		return
	}

	// prepare the response data and return it
	respData := make(map[string]interface{})
	respData["message"] = "Successfully scheduled attack"
	respData["attackId"] = res.InsertedID
	util.JsonResponse(w, respData, http.StatusOK)
}

// This function is called when we want to execute an attack immediately (upon the user's request).
func TriggerAttackNow(w http.ResponseWriter, r *http.Request) {
	// decode the request data
	var attack PendingAttackObj
	err := json.NewDecoder(r.Body).Decode(&attack)
	if err != nil {
		fmt.Printf("[TriggerAttackNow] Failed to decode request data: %+v\n", err)
		util.JsonResponse(w, "Request data is invalid", http.StatusBadRequest)
		return
	}

	cli := db.GetClient()
	if cli == nil {
		fmt.Println("[TriggerAttackNow] Failed to connect to DB")
		util.JsonResponse(w, "Failed to connect to DB", http.StatusBadGateway)
		return
	}

	// get handles for some collections
	custDb := cli.Database(db.VedikaCorpDatabase)
	pendingAttacksColl := custDb.Collection(db.PendingAttacksCollection)
	attackEmailsColl := custDb.Collection(db.EmailsCollection)
	attackLogColl := custDb.Collection(db.AttackLogCollection)

	// execute the attack in a goroutine so we don't have to wait for it
	go func() {
		// call a helper function to execute the attack
		err = executeAttack(attackEmailsColl, attackLogColl, pendingAttacksColl, attack)
		if err != nil {
			fmt.Printf("[TriggerPendingAttack] Failed to execute this attack: %+v | %+v\n", attack, err)
		}
	}()

	// prepare the response data and return it
	respData := make(map[string]interface{})
	respData["message"] = "Successfully triggered attack"
	util.JsonResponse(w, respData, http.StatusOK)
}

// This function is called when a link in an attack email is clicked.
func RecordAttackResults(w http.ResponseWriter, r *http.Request) {
	// get the email ID from the URL parameters
	vars := mux.Vars(r)
	attackIdHex := vars[util.URLParameterAttackId]

	// check if the emailId is empty
	if len(attackIdHex) == 0 {
		fmt.Println("[RecordAttackResults] Attack ID is missing from the request")
		util.JsonResponse(w, "Request is missing attack ID", http.StatusBadRequest)
		return
	}

	// convert the given ObjectID hex into a primitive.ObjectID
	objId, err := primitive.ObjectIDFromHex(attackIdHex)
	if err != nil {
		fmt.Printf("[RecordAttackResults] Failed to convert attack ID: %+v", err)
		util.JsonResponse(w, "Attack ID is invalid", http.StatusBadRequest)
		return
	}

	cli := db.GetClient()
	if cli == nil {
		fmt.Println("[RecordAttackResults] Failed to connect to DB")
		util.JsonResponse(w, "Failed to connect to DB", http.StatusBadGateway)
		return
	}

	// get a handle for the AttackLog collection
	attackLogColl := cli.Database(db.VedikaCorpDatabase).Collection(db.AttackLogCollection)

	// set the query filter to match the _id
	// NOTE: When we execute an attack, we set the _id of the corresponding document in the AttackLog collection as the
	//       _id of the document in the PendingAttacks collection.
	filter := bson.D{
		{Key: "_id", Value: objId},
	}

	// set the update query
	update := bson.D{
		{Key: "$set", Value: bson.D{
			{Key: "Results.IsSuccessful", Value: true},
			{Key: "Results.ClickTime", Value: time.Now()},
		}},
	}

	// submit the query
	res, err := attackLogColl.UpdateOne(context.TODO(), filter, update)
	if err != nil {
		fmt.Printf("[RecordAttackResults] Failed to record attack results for '%s': %+v\n", objId.Hex(), err)
		util.JsonResponse(w, map[string]string{"error": "Failed to record attack results"}, http.StatusOK)
		return
	}

	fmt.Printf("[RecordAttackResults] Successfully recorded attack results for '%s': %d | %d\n", objId.Hex(), res.MatchedCount, res.ModifiedCount)
	util.JsonResponse(w, map[string]string{"message": "Successfully recorded attack results"}, http.StatusOK)
}
