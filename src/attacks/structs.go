package attacks

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type RecipientObj struct {
	Name    string `json:"name" bson:"Name"`
	Address string `json:"address" bson:"Address"`
}

// Represents documents in the PendingAttacks collection
type PendingAttackObj struct {
	ObjId           primitive.ObjectID `json:"_id" bson:"_id,omitempty"`
	EmailId         primitive.ObjectID `json:"emailId" bson:"EmailId,omitempty"`
	TargetRecipient RecipientObj       `json:"targetRecipient" bson:"TargetRecipient"`
	TargetUserId    primitive.ObjectID `json:"targetUserId" bson:"TargetUserId"`
	TriggerTime     time.Time          `json:"attackTime" bson:"AttackTime"`
}

type AttackLogResults struct {
	IsSuccessful bool      `json:"isSuccessful" bson:"IsSuccessful"`
	ClickTime    time.Time `json:"clickTime" bson:"ClickTime"`
}

// Represents documents in the AttackLog collection. It contains info about an attack that was executed.
type AttackLogObj struct {
	ObjId           primitive.ObjectID `json:"_id" bson:"_id,omitempty"`
	EmailId         primitive.ObjectID `json:"emailId" bson:"EmailId,omitempty"`
	TargetRecipient RecipientObj       `json:"targetRecipient" bson:"TargetRecipient"`
	TargetUserId    primitive.ObjectID `json:"targetUserId" bson:"TargetUserId"`
	TriggerTime     time.Time          `json:"attackTime" bson:"AttackTime"`
	Results         AttackLogResults   `json:"results" bson:"Results"`
}
