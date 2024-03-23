package agents

import "time"

const (
	MsGraphV1URL  = "https://graph.microsoft.com/v1.0"
	SendMailRoute = "/me/sendMail"
)

type AgentObj struct {
	Name  string   `json:"name" bson:"Name"`
	Email string   `json:"email" bson:"Email"`
	Token TokenObj `json:"token" bson:"Token"`
}

type TokenObj struct {
	AccessToken  string    `json:"accessToken" bson:"AccessToken"`
	RefreshToken string    `json:"refreshToken" bson:"RefreshToken"`
	ExpiryTime   time.Time `json:"expiryTime" bson:"ExpiryTime"`
}

// -------------- SEND MAIL ------------------
type sendMailRequestData struct {
	Message         smMessageObj `json:"message"`
	SaveToSentItems bool         `json:"saveToSentItems"`
}

type smMessageObj struct {
	Subject      string           `json:"subject"`
	Body         smBodyObj        `json:"body"`
	ToRecipients []smRecipientObj `json:"toRecipients"`
}

type smBodyObj struct {
	ContentType string `json:"contentType"`
	Content     string `json:"content"`
}

type smRecipientObj struct {
	EmailAddress smEmailAddressObj `json:"emailAddress"`
}

type smEmailAddressObj struct {
	Address string `json:"address"`
}
