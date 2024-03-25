package emails

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Represents documents in the Emails collection
type EmailObj struct {
	Id      primitive.ObjectID `json:"_id" bson:"_id,omitempty"`
	Name    string             `json:"name" bson:"Name"`
	Type    string             `json:"type" bson:"Type"`
	Company string             `json:"company" bson:"Company"`
	Subject string             `json:"subject" bson:"Subject"`
	Body    string             `json:"body" bson:"Body"`
}
