package notification

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
)

// SlackNotificationParams holds the incoming arguments.
type SlackNotificationParams struct {
	Title      string
	Message    string
	Channel    string
	WebhookUrl string
}

// SlackPayload is the body structure we send to Slack.
type SlackPayload struct {
	Channel string `json:"channel"`
	Text    string `json:"text"`
}

// SlackNotification is a rough Go translation of your TypeScript function.
func SendSlackNotification(params SlackNotificationParams) {
	log.Printf("Sending slack alert for %s", params.Title)

	payload := SlackPayload{
		Channel: params.Channel,
		Text:    fmt.Sprintf("*%s*\n\n%s", params.Title, params.Message),
	}

	bodyBytes, err := json.Marshal(payload)
	if err != nil {
		log.Printf("Error marshaling Slack payload: %v", err)
		return
	}

	req, err := http.NewRequest("POST", params.WebhookUrl, bytes.NewBuffer(bodyBytes))
	if err != nil {
		log.Printf("Error creating Slack request: %v", err)
		return
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Printf("Slack request failed: %v", err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		respBody, _ := io.ReadAll(resp.Body)
		log.Printf("Slack request failed with status %d, body: %s", resp.StatusCode, string(respBody))
		return
	}
}
