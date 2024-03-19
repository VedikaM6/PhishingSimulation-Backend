package emails

import (
	"context"
	"fmt"
	"net/http"

	"example.com/m/src/db"
	"example.com/m/src/util"
	"go.mongodb.org/mongo-driver/bson"
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

	fmt.Printf("[ListAttacks][DEBUG] Successfully got %d emails!", len(allAttackEmails))

	// return the retrieved users
	respData := make(map[string][]AttackEmailObj)
	respData["emails"] = allAttackEmails
	util.JsonResponse(w, respData, http.StatusOK)
}
