package emails

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Represents documents in the AttackEmails collection
type AttackEmailObj struct {
	Id      primitive.ObjectID `json:"_id" bson:"_id,omitempty"`
	Subject string             `json:"subject" bson:"Subject"`
	Body    string             `json:"body" bson:"Body"`
}
