package attacks

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"example.com/m/src/db"
	"example.com/m/src/emails"
	"example.com/m/src/util"
	"go.mongodb.org/mongo-driver/bson"
)

// This function is called when we want to check for any scheduled attacks.
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
