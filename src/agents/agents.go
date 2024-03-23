package agents

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"example.com/m/src/emails"
)

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
