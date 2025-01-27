package notification_service

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"log/slog"
	"net/http"
)

type NotificationService struct {
	Logger      *slog.Logger
	SlackConfig *struct {
		SlackWebhookUrl string
		SlackChannel    string
	}
	EmailConfig *struct {
		Host     string
		Port     int
		Username string
		Password string
	}
}

// SlackNotificationParams holds the incoming arguments.
type SlackNotificationParams struct {
	Title   string
	Message string
}

// SlackPayload is the body structure we send to Slack.
type SlackPayload struct {
	Channel string `json:"channel"`
	Text    string `json:"text"`
}

func (ns *NotificationService) SendSlackNotification(params SlackNotificationParams) {
	log.Printf("Sending slack alert for %s", params.Title)

	payload := SlackPayload{
		Channel: ns.SlackConfig.SlackChannel,
		Text:    fmt.Sprintf("*%s*\n\n%s", params.Title, params.Message),
	}

	bodyBytes, err := json.Marshal(payload)
	if err != nil {
		log.Printf("Error marshaling Slack payload: %v", err)
		return
	}

	req, err := http.NewRequest("POST", ns.SlackConfig.SlackWebhookUrl, bytes.NewBuffer(bodyBytes))
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
