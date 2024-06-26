package attacks

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type AttackDBObj interface {
	LogAttack(coll *mongo.Collection) (*mongo.InsertOneResult, error)
}

type RecipientObj struct {
	Name        string    `json:"name" bson:"Name"`
	Address     string    `json:"address" bson:"Address"`
	IsClicked   bool      `json:"isClicked" bson:"IsClicked"`
	ClickedTime time.Time `json:"clickedTime" bson:"ClickedTime"`
}

// Represents documents in the PendingAttacks collection
type PendingAttackObj struct {
	ObjId            primitive.ObjectID   `json:"_id" bson:"_id,omitempty"`
	Name             string               `json:"name" bson:"Name"`
	Description      string               `json:"description" bson:"Description"`
	EmailId          primitive.ObjectID   `json:"emailId" bson:"EmailId,omitempty"`
	TargetRecipients []RecipientObj       `json:"targetRecipients" bson:"TargetRecipients"`
	TargetUserIds    []primitive.ObjectID `json:"targetUserIds" bson:"TargetUserIds"`
	TriggerTime      time.Time            `json:"triggerTime" bson:"TriggerTime"`
}

func (pao *PendingAttackObj) CreateAttack(pendingAttacksColl *mongo.Collection) (*mongo.InsertOneResult, error) {
	// insert the object into the PendingAttacks collection
	return pendingAttacksColl.InsertOne(context.TODO(), pao)
}

type AttackLogResults struct {
	IsSuccessful bool      `json:"isSuccessful" bson:"IsSuccessful"`
	ClickTime    time.Time `json:"clickTime" bson:"ClickTime"`
}

// An object used to indicate which email was used in an attack
type UsedEmailObj struct {
	Id   primitive.ObjectID `json:"_id" bson:"_id"`
	Name string             `json:"name" bson:"Name"`
}

// Represents documents in the AttackLog collection. It contains info about an attack that was executed.
type AttackLogObj struct {
	ObjId            primitive.ObjectID   `json:"_id" bson:"_id,omitempty"`
	Name             string               `json:"name" bson:"Name"`
	Description      string               `json:"description" bson:"Description"`
	UsedEmail        UsedEmailObj         `json:"usedEmail" bson:"UsedEmail,omitempty"`
	TargetRecipients []RecipientObj       `json:"targetRecipients" bson:"TargetRecipients"`
	TargetUserIds    []primitive.ObjectID `json:"targetUserIds" bson:"TargetUserIds"`
	TriggerTime      time.Time            `json:"triggerTime" bson:"TriggerTime"`
}

func (pao *PendingAttackObj) LogAttack(attackLogColl *mongo.Collection, emailTemplateName string) (*mongo.InsertOneResult, error) {
	// insert the object into the AttackLog collection
	alo := AttackLogObj{
		ObjId:       pao.ObjId,
		Name:        pao.Name,
		Description: pao.Description,
		UsedEmail: UsedEmailObj{
			Id:   pao.EmailId,
			Name: emailTemplateName,
		},
		TargetRecipients: pao.TargetRecipients,
		TargetUserIds:    pao.TargetUserIds,
		TriggerTime:      pao.TriggerTime,
	}

	return alo.LogAttack(attackLogColl)
}

func (alo *AttackLogObj) LogAttack(attackLogColl *mongo.Collection) (*mongo.InsertOneResult, error) {
	// insert the object into the AttackLog collection
	return attackLogColl.InsertOne(context.TODO(), alo)
}
