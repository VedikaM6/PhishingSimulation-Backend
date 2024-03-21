package attacks

import (
	"context"
	"fmt"
	"time"

	"example.com/m/src/emails"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func deletePendingAttack(pendingAttackColl *mongo.Collection, objId primitive.ObjectID) (*mongo.DeleteResult, error) {
	// set the query filter to match the document
	filter := bson.D{
		{Key: "_id", Value: objId},
	}

	// submit the query
	return pendingAttackColl.DeleteOne(context.TODO(), filter)
}

// This function does the work to 'execute an attack'. It does the following steps:
// 1. Get the email contents (AttackEmails collection)
// 2. Send the email via Microsoft Graph
// 3. Log the attack (AttackLog collection)
// 4. Delete the pending attack from the database so we don't execute it again (PEndingAttacks collection)
func executeAttack(attackEmailsColl, attackLogColl, pendingAttacksColl *mongo.Collection, pendAttack PendingAttackObj) error {
	// get the email to send from the DB
	email, err := emails.GetEmailById(attackEmailsColl, pendAttack.EmailId)
	if err != nil {
		fmt.Printf("[executeAttack] Failed to get email with ID '%s': %+v\n", pendAttack.EmailId.Hex(), err)
		return err
	}

	// TODO: send email
	fmt.Printf("[executeAttack] TODO: Send email: %+v\n", email)

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
		fmt.Printf("[executeAttack] Failed to remove document '%s': %+v\n", pendAttack.ObjId.Hex(), err)
	} else {
		fmt.Printf("[executeAttack] Successfully removed document: %d\n", res.DeletedCount)
	}
}
