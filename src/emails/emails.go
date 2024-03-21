package emails

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"example.com/m/src/db"
	"example.com/m/src/util"
	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

// This function is called when we want to get a list of emails from the database.
func ListEmails(w http.ResponseWriter, r *http.Request) {
	cli := db.GetClient()
	if cli == nil {
		fmt.Println("[ListEmails] Failed to connect to DB")
		util.JsonResponse(w, "Failed to connect to DB", http.StatusBadGateway)
		return
	}

	// get a handle for the AttackEmails collection
	attackEmailsColl := cli.Database(db.VedikaCorpDatabase).Collection(db.AttackEmailsCollection)

	// get all users in the collection
	ctx := context.TODO()
	cur, err := attackEmailsColl.Find(ctx, bson.D{})
	if err != nil {
		fmt.Printf("[ListEmails] Failed to get emails from DB: %+v\n", err)
		util.JsonResponse(w, "Failed to get emails", http.StatusBadGateway)
		return
	}

	defer cur.Close(ctx)

	// decode all the results into a slice
	allAttackEmails := make([]AttackEmailObj, 0)
	err = cur.All(ctx, &allAttackEmails)
	if err != nil {
		fmt.Printf("[ListEmails] Failed to decode results: %+v\n", err)
		util.JsonResponse(w, "Failed to get users", http.StatusBadGateway)
		return
	}

	fmt.Printf("[ListEmails][DEBUG] Successfully got %d emails!\n", len(allAttackEmails))

	// return the retrieved users
	respData := make(map[string][]AttackEmailObj)
	respData["emails"] = allAttackEmails
	util.JsonResponse(w, respData, http.StatusOK)
}

// This endpoint is called when we want to get a specific email from the database.
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

	// get the email from the database
	email, err := GetEmailById(attackEmailsColl, objId)
	if err != nil {
		fmt.Printf("[GetAttackEmail] Failed to get email: %+v\n", err)
		if err == mongo.ErrNoDocuments {
			util.JsonResponse(w, "The specified email was not found", http.StatusNotFound)
		} else {
			util.JsonResponse(w, "Failed to get email", http.StatusBadGateway)
		}
		return
	}

	// return the email
	util.JsonResponse(w, email, http.StatusOK)
}

// This endpoint is called when a user wants to create their own custom email.
func CreateNewEmail(w http.ResponseWriter, r *http.Request) {
	// decode the request data
	var newEmail AttackEmailObj
	err := json.NewDecoder(r.Body).Decode(&newEmail)
	if err != nil {
		fmt.Printf("[CreateNewEmail] Failed to decode request data: %+v\n", err)
		util.JsonResponse(w, "Request data is invalid", http.StatusBadRequest)
		return
	}

	// validate the email details
	validationErr := ""
	if newEmail.Subject == "" {
		// No subject was specified for this email
		validationErr = "You must specify a subject for this email."
	} else if newEmail.Body == "" {
		// No body was specified for this email
		validationErr = "You must specify a body for this email."
	} else if !newEmail.Id.IsZero() {
		// The request specifies an _id but it shouldn't do that.
		validationErr = "Request should not specify _id"
	}

	if validationErr != "" {
		util.JsonResponse(w, validationErr, http.StatusBadRequest)
		return
	}

	// connect to the database
	cli := db.GetClient()
	if cli == nil {
		fmt.Println("[CreateNewEmail] Failed to connect to DB")
		util.JsonResponse(w, "Failed to connect to DB", http.StatusBadGateway)
		return
	}

	// get a handle for the AttackEmails collection
	emailsColl := cli.Database(db.VedikaCorpDatabase).Collection(db.AttackEmailsCollection)

	// log the
	res, err := emailsColl.InsertOne(context.TODO(), newEmail)
	if err != nil {
		fmt.Printf("[CreateNewEmail] Failed to insert new email: %+v\n", err)
		util.JsonResponse(w, "Failed to create email", http.StatusBadGateway)
		return
	}

	// prepare the response data and return it
	respData := make(map[string]interface{})
	respData["message"] = "Successfully create new email"
	respData["emailId"] = res.InsertedID
	util.JsonResponse(w, respData, http.StatusOK)
}
