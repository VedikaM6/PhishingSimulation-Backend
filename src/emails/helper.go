package emails

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func GetEmailById(emailsColl *mongo.Collection, emailId primitive.ObjectID) (EmailObj, error) {
	// set the query filter to match the email
	filter := bson.D{
		{Key: "_id", Value: emailId},
	}

	ctx := context.TODO()
	var email EmailObj
	err := emailsColl.FindOne(ctx, filter).Decode(&email)
	return email, err
}
