package attacks

import (
	"context"
	"fmt"
	"regexp"
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

	// take a copy of the email body while it still contains the placeholder
	emailBodyWithPlaceholder := email.Body

	// send the email
	for _, recip := range pendAttack.TargetRecipients {
		// replace the phishing link placeholder with an <a> tag
		email.Body = insertPhishingAnchorTag(emailBodyWithPlaceholder, getPhishingLink(pendAttack.ObjId.Hex(), recip.Address))

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

// Generate the phishing link with the given details
func getPhishingLink(attackId, userEmail string) string {
	return fmt.Sprintf("http://localhost:80/attacks/clicked/%s/%s", attackId, userEmail)
}

// Find and replace the placeholder for the phishing link in the email content
func insertPhishingAnchorTag(emailBody string, phishingLink string) string {
	// create a Regexp object
	regexObj, err := regexp.Compile(`\[\[[A-Za-z]+\]\]`)
	if err != nil {
		fmt.Printf("[InsertPhishingAnchorTag] Failed to compile regex pattern: %+v", err)
		return emailBody
	}

	// create the anchor tag text
	indexOfOpenDoubleBrackets := strings.Index(emailBody, "[[")
	indexOfCloseDoubleBrackets := strings.Index(emailBody, "]]")
	linkText := emailBody[indexOfOpenDoubleBrackets+2 : indexOfCloseDoubleBrackets]
	anchorTagText := fmt.Sprintf("<a href=\"%s\">%s</a>", phishingLink, linkText)

	// replace all instances of the placeholder with the anchor tag
	return regexObj.ReplaceAllString(emailBody, anchorTagText)

}
