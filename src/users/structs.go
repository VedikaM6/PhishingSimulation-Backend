package users

import "go.mongodb.org/mongo-driver/bson/primitive"

type UserObj struct {
	Id    primitive.ObjectID `json:"_id" bson:"_id,omitempty"`
	Name  string             `json:"name" bson:"Name"`
	Email string             `json:"email" bson:"Email"`
}
