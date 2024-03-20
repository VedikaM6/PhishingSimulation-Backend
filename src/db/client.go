package db

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const (
	VedikaCorpDatabase = "Vedika_Corp"

	MailboxesCollection      = "Mailboxes"
	UsersCollection          = "Users"
	AttackEmailsCollection   = "AttackEmails"
	AttackLogCollection      = "AttackLog"
	PendingAttacksCollection = "PendingAttacks"
)

var client *mongo.Client
var cancelCallback context.CancelFunc
var ctx context.Context

func DisconnectClient() {
	// cancel the context
	cancelCallback()

	// disconnect the client
	if err := client.Disconnect(ctx); err != nil {
		panic(err)
	}
}

func InitClient() {
	// establish a connection to the local Mongo database
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	cancelCallback = cancel
	newClient, err := mongo.Connect(ctx, options.Client().ApplyURI("mongodb://localhost:27017"))
	if err != nil {
		panic(err)
	}

	client = newClient
}

func GetClient() *mongo.Client {
	return client
}
