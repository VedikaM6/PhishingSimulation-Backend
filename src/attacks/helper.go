package attacks

import (
	"context"
	"fmt"
	"strings"
	"time"

	"example.com/m/src/agents"
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
// 1. Get the email contents (Emails collection)
// 2. Send the email via Microsoft Graph
// 3. Log the attack (AttackLog collection)
// 4. Delete the pending attack from the database so we don't execute it again (PEndingAttacks collection)
func executeAttack(emailsColl, attackLogColl, pendingAttacksColl *mongo.Collection, pendAttack PendingAttackObj) error {
	// get the email to send from the DB
	email, err := emails.GetEmailById(emailsColl, pendAttack.EmailId)
	if err != nil {
		fmt.Printf("[executeAttack] Failed to get email with ID '%s': %+v\n", pendAttack.EmailId.Hex(), err)
		return err
	}

	// replace all \n with <br/>
	email.Body = strings.ReplaceAll(email.Body, "\n", "<br/>")

	// send the email
	for _, recip := range pendAttack.TargetRecipients {
		err = agents.SendEmailWithRandomAgent(email, recip.Address)
		if err != nil {
			fmt.Printf("[executeAttack] Failed to send email '%s': %+v\n", pendAttack.EmailId.Hex(), err)
			// TODO: Don't return here because we still need to implement the access token, so this will always fail.
		} else {
			fmt.Printf("[executeAttack][DEBUG] Successfully sent email '%s' to '%s'\n", pendAttack.EmailId.Hex(), recip.Address)
		}

		time.Sleep(time.Millisecond * 50)
	}

	// log the attack in the AttackLog
	_, err = pendAttack.LogAttack(attackLogColl, email.Name)
	if err != nil {
		fmt.Printf("[executeAttack] Failed to log attack '%s': %+v\n", pendAttack.ObjId.Hex(), err)
		return err
	}

	// remove the document from the PendingAttacks collection because it has been processed
	res, err := deletePendingAttack(pendingAttacksColl, pendAttack.ObjId)
	if err != nil {
		fmt.Printf("[executeAttack] Failed to remove document '%s': %+v\n", pendAttack.ObjId.Hex(), err)
		return err
	} else {
		fmt.Printf("[executeAttack] Successfully removed document: %d\n", res.DeletedCount)
	}

	return nil
}
