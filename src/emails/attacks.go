package emails

import (
	"context"
	"fmt"
	"net/http"

	"example.com/m/src/db"
	"example.com/m/src/util"
	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func ListAttacks(w http.ResponseWriter, r *http.Request) {
	cli := db.GetClient()
	if cli == nil {
		fmt.Println("[ListAttacks] Failed to connect to DB")
		util.JsonResponse(w, "Failed to connect to DB", http.StatusBadGateway)
		return
	}

	// get a handle for the AttackEmails collection
	attackEmailsColl := cli.Database(db.VedikaCorpDatabase).Collection(db.AttackEmailsCollection)

	// get all users in the collection
	ctx := context.TODO()
	cur, err := attackEmailsColl.Find(ctx, bson.D{})
	if err != nil {
		fmt.Printf("[ListAttacks] Failed to get emails from DB: %+v\n", err)
		util.JsonResponse(w, "Failed to get emails", http.StatusBadGateway)
		return
	}

	defer cur.Close(ctx)

	// decode all the results into a slice
	allAttackEmails := make([]AttackEmailObj, 0)
	err = cur.All(ctx, &allAttackEmails)
	if err != nil {
		fmt.Printf("[ListAttacks] Failed to decode results: %+v\n", err)
		util.JsonResponse(w, "Failed to get users", http.StatusBadGateway)
		return
	}

	fmt.Printf("[ListAttacks][DEBUG] Successfully got %d emails!\n", len(allAttackEmails))

	// return the retrieved users
	respData := make(map[string][]AttackEmailObj)
	respData["emails"] = allAttackEmails
	util.JsonResponse(w, respData, http.StatusOK)
}

func GetAttackEmail(w http.ResponseWriter, r *http.Request) {
	// get the email ID from the URL parameters
	vars := mux.Vars(r)
	emailIdHex := vars[util.URLParameterEmailId]

	// check if the emailId is empty
	if len(emailIdHex) == 0 {
		fmt.Println("[GetAttackEmail] email ID is missing from request")
		util.JsonResponse(w, "Request is missing email ID", http.StatusBadRequest)
		return
	}

	// convert the given ObjectID hex into a primitive.ObjectID
	objId, err := primitive.ObjectIDFromHex(emailIdHex)
	if err != nil {
		fmt.Printf("[GetAttackEmail] Failed to convert email ID: %+v", err)
		util.JsonResponse(w, "Email ID is invalid", http.StatusBadRequest)
		return
	}

	cli := db.GetClient()
	if cli == nil {
		fmt.Println("[GetAttackEmail] Failed to connect to DB")
		util.JsonResponse(w, "Failed to connect to DB", http.StatusBadGateway)
		return
	}

	// get a handle for the AttackEmails collection
	attackEmailsColl := cli.Database(db.VedikaCorpDatabase).Collection(db.AttackEmailsCollection)

	// set the query filter to match the email
	filter := bson.D{
		{Key: "_id", Value: objId},
	}

	ctx := context.TODO()
	var email AttackEmailObj
	err = attackEmailsColl.FindOne(ctx, filter).Decode(&email)
	if err != nil {
		fmt.Printf("[GetAttackEmail] Failed to get email from DB: %+v\n", err)
		util.JsonResponse(w, "Failed to get emails", http.StatusBadGateway)
		return
	}

	// return the email
	util.JsonResponse(w, email, http.StatusOK)
}
