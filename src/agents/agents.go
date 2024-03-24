package agents

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"

	"example.com/m/src/db"
	"example.com/m/src/emails"
	"go.mongodb.org/mongo-driver/bson"
)

func getRandomAgent() (AgentObj, error) {
	cli := db.GetClient()
	if cli == nil {
		fmt.Printf("[getRandomAgent] Failed to connect to DB\n")
		return AgentObj{}, errors.New("failed to connect to DB")
	}

	// get a handle for the Agents collection
	agentsColl := cli.Database(db.VedikaCorpDatabase).Collection(db.AgentsCollection)

	// create an aggregation pipeline to get a random document
	aggPipeline := bson.A{
		bson.D{
			{Key: "$sample", Value: 1},
		},
	}

	// submit the query
	ctx := context.TODO()
	cur, err := agentsColl.Aggregate(ctx, aggPipeline)
	if err != nil {
		fmt.Printf("[getRandomAgent] Failed to get randomg agent: %+v\n", err)
		return AgentObj{}, err
	}

	defer cur.Close(ctx)

	var agent AgentObj
	if cur.Next(ctx) {
		// Got the first document
		err := cur.Decode(&agent)
		if err != nil {
			// Since we're only getting 1 document, return this error
			fmt.Printf("[getRandomgAgent] Failed to decode document: %+v\n", err)
			return agent, err
		}
	}

	if cur.Err() != nil {
		fmt.Printf("[getRandomAgent] A Mongo cursor error occurred! %+v\n", cur.Err())
	}

	return agent, nil
}

func SendEmailWithRandomAgent(email emails.EmailObj, recipient string) error {
	// get an arbitrary agent from the database
	randAgent, err := getRandomAgent()
	if err != nil {
		fmt.Printf("[SendEmailWithRandomAgent][%s] Failed to get agent: %+v\n", recipient, err)
		return err
	}

	statCode, err := randAgent.SendEmail(email, recipient)
	if err != nil {
		fmt.Printf("[SendEmailWithRandomAgent][%s] Failed to send email: %d | %+v\n", recipient, statCode, err)
		return err
	}

	return nil
}

func (a *AgentObj) SendEmail(email emails.EmailObj, recipient string) (int, error) {
	// get an HTTP client
	cli := http.Client{}

	// prepare the request data
	reqData := sendMailRequestData{
		Message: smMessageObj{
			Subject: email.Subject,
			Body: smBodyObj{
				ContentType: "HTML",
				Content:     email.Body,
			},
			ToRecipients: []smRecipientObj{
				smRecipientObj{
					EmailAddress: smEmailAddressObj{
						Address: recipient,
					},
				},
			},
		},
		SaveToSentItems: true,
	}

	// marshal the request data
	reqDataMarsh, err := json.Marshal(reqData)
	if err != nil {
		fmt.Printf("[SendEmail][%s] Failed to marshal request data: %+v\n", a.Email, err)
		return http.StatusInternalServerError, err
	}

	// create the request object
	req, err := http.NewRequest(http.MethodPost, MsGraphV1URL+SendMailRoute, bytes.NewBuffer(reqDataMarsh))
	if err != nil {
		fmt.Printf("[SendEmail][%s] Failed to create request: %+v\n", a.Email, err)
		return http.StatusInternalServerError, err
	}

	// set the request headers
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", a.Token.AccessToken)

	// send the request
	resp, err := cli.Do(req)
	if err != nil {
		fmt.Printf("[SendEmail][%s] Failed to send request: %+v\n", a.Email, err)
		return http.StatusBadGateway, err
	}

	defer resp.Body.Close()

	// empty the buffer and then close the response body
	_, err = io.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("[SendEmail][%s] An error occurred when reading the response: %+v\n", a.Email, err)
	}

	return resp.StatusCode, nil
}
