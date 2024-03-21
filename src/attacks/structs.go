package attacks

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
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
	TriggerTime     time.Time          `json:"triggerTime" bson:"TriggerTime"`
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
	TriggerTime     time.Time          `json:"triggerTime" bson:"TriggerTime"`
	Results         AttackLogResults   `json:"results" bson:"Results"`
}

func (alo *AttackLogObj) LogAttack(attackLogColl *mongo.Collection) (*mongo.InsertOneResult, error) {
	// insert the object into the AttackLog collection
	return attackLogColl.InsertOne(context.TODO(), alo)
}
