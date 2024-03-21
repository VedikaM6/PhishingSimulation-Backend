package attacks

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"example.com/m/src/db"
	"example.com/m/src/emails"
	"example.com/m/src/util"
	"go.mongodb.org/mongo-driver/bson"
)

// This function is called when we want to check for any scheduled attacks and trigger them (if they are due).
// It should be invoked routinely.
func TriggerAttacks(w http.ResponseWriter, r *http.Request) {
	cli := db.GetClient()
	if cli == nil {
		fmt.Println("[TriggerAttacks] Failed to connect to DB")
		util.JsonResponse(w, "Failed to connect to DB", http.StatusBadGateway)
		return
	}

	// get a handle for the PendingAttacks collection
	custDb := cli.Database(db.VedikaCorpDatabase)
	pendingAttacksColl := custDb.Collection(db.PendingAttacksCollection)
	attackEmailsColl := custDb.Collection(db.AttackEmailsCollection)
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
		fmt.Printf("[TriggerAttacks] Failed to get emails from DB: %+v\n", err)
		util.JsonResponse(w, "Failed to trigger attacks", http.StatusBadGateway)
		return
	}

	defer cur.Close(ctx)

	// decode all the results into a slice
	pendingAttacks := make([]PendingAttackObj, 0)
	err = cur.All(ctx, &pendingAttacks)
	if err != nil {
		fmt.Printf("[TriggerAttacks] Failed to decode results: %+v\n", err)
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
			// get the email to send from the DB
			email, err := emails.GetEmailById(attackEmailsColl, pendAttack.EmailId)
			if err != nil {
				fmt.Printf("[TriggerAttacks] Failed to get email with ID '%s': %+v\n", pendAttack.EmailId.Hex(), err)
				continue
			}

			// TODO: send email
			fmt.Printf("[TriggerAttacks] TODO: Send email: %+v\n", email)

			// log the attack in the AttackLog
			log := AttackLogObj{
				EmailId:         pendAttack.EmailId,
				TargetRecipient: pendAttack.TargetRecipient,
				TargetUserId:    pendAttack.TargetUserId,
				TriggerTime:     pendAttack.TriggerTime,
				Results: AttackLogResults{
					IsSuccessful: false, // The attack only becomes successful if the user clicks on the link in the email.
					ClickTime:    time.Time{},
				},
			}
			log.LogAttack(attackLogColl)

			// remove the document from the PendingAttacks collection because it has been processed
			res, err := deletePendingAttack(pendingAttacksColl, pendAttack.ObjId)
			if err != nil {
				fmt.Printf("[TriggerAttacks] Failed to remove document '%s': %+v\n", pendAttack.ObjId.Hex(), err)
			} else {
				fmt.Printf("[TriggerAttacks] Successfully removed document: %d\n", res.DeletedCount)
			}
		}
	}()
}

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
