package emails

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func GetEmailById(attackEmailsColl *mongo.Collection, emailId primitive.ObjectID) (AttackEmailObj, error) {
	// set the query filter to match the email
	filter := bson.D{
		{Key: "_id", Value: emailId},
	}

	ctx := context.TODO()
	var email AttackEmailObj
	err := attackEmailsColl.FindOne(ctx, filter).Decode(&email)
	return email, err
}
