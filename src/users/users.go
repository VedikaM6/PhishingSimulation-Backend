package users

import (
	"context"
	"fmt"
	"net/http"

	"example.com/m/src/db"
	"example.com/m/src/util"
	"go.mongodb.org/mongo-driver/bson"
)

func ListUsers(w http.ResponseWriter, r *http.Request) {
	cli := db.GetClient()
	if cli == nil {
		fmt.Println("[ListUsers] Failed to connect to DB")
		util.JsonResponse(w, "Failed to connect to DB", http.StatusBadGateway)
		return
	}

	// get a handle for the Users collection
	usersColl := cli.Database(db.VedikaCorpDatabase).Collection(db.UsersCollection)

	// get all users in the collection
	ctx := context.TODO()
	cur, err := usersColl.Find(ctx, bson.D{})
	if err != nil {
		fmt.Printf("[ListUsers] Failed to get users from DB: %+v\n", err)
		util.JsonResponse(w, "Failed to get users", http.StatusBadGateway)
		return
	}

	defer cur.Close(ctx)

	// decode all the results into a slice
	allUsers := make([]UserObj, 0)
	err = cur.All(ctx, &allUsers)
	if err != nil {
		fmt.Printf("[ListUsers] Failed to decode results: %+v\n", err)
		util.JsonResponse(w, "Failed to get users", http.StatusBadGateway)
		return
	}

	fmt.Printf("[ListUsers][DEBUG] Successfully got %d users!", len(allUsers))

	// return the retrieved users
	respData := make(map[string][]UserObj)
	respData["users"] = allUsers
	util.JsonResponse(w, respData, http.StatusOK)

	// writer.Header().Set("Content-Type", "application/json")
	// writer.WriteHeader(http.StatusOK)
	// var human Bio
	// err := json.NewDecoder(request.Body).Decode(&human)
	// if err != nil {
	// log.Fatalln("There was an error decoding the request body into the struct")
	// }
	// BioData = append(BioData, human)
	// err = json.NewEncoder(writer).Encode(&human)
	// if err != nil {
	// log.Fatalln("There was an error encoding the initialized struct")
	// }
}
