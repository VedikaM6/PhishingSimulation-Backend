package emails

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Represents documents in the AttackEmails collection
type AttackEmailObj struct {
	Id      primitive.ObjectID `json:"_id" bson:"_id,omitempty"`
	Subject string             `json:"subject" bson:"Subject"`
	Body    string             `json:"body" bson:"Body"`
}

type RecipientObj struct {
	Name    string `json:"name" bson:"Name"`
	Address string `json:"address" bson:"Address"`
}

// Represents documents in the PendingAttacks collection
type PendingAttacks struct {
	EmailId         primitive.ObjectID `json:"emailId" bson:"EmailId,omitempty"`
	TargetRecipient RecipientObj       `json:"targetRecipient" bson:"TargetRecipient"`
	TargetUserId    primitive.ObjectID `json:"targetUserId" bson:"TargetUserId"`
	AttackTime      time.Time          `json:"attackTime" bson:"AttackTime"`
}
