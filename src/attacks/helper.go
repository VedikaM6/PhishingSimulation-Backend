package attacks

import (
	"context"

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
