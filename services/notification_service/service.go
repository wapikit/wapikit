package notification_service

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"log/slog"
	"net/http"
	"net/smtp"
	"strings"
)

type SlackConfig struct {
	SlackWebhookUrl string
	SlackChannel    string
}

type EmailConfig struct {
	Host     string
	Port     string
	Username string
	Password string
}

type NotificationService struct {
	Logger      *slog.Logger
	SlackConfig *SlackConfig
	EmailConfig *EmailConfig
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

func (ns *NotificationService) SendEmail(sendToEmail string, subject string, body string) error {
	if ns.EmailConfig == nil {
		return fmt.Errorf("email configuration is not set")
	}

	// Compose the email message
	header := fmt.Sprintf("From: %s\r\nTo: %s\r\nSubject: %s\r\n\r\n", ns.EmailConfig.Username, sendToEmail, subject)
	message := header + body

	// Authenticate the sender
	auth := smtp.PlainAuth("", ns.EmailConfig.Username, ns.EmailConfig.Password, ns.EmailConfig.Host)

	address := strings.Join([]string{ns.EmailConfig.Host, ns.EmailConfig.Port}, ":")

	// Send the email
	err := smtp.SendMail(
		address,
		auth,
		ns.EmailConfig.Username,
		[]string{sendToEmail},
		[]byte(message),
	)

	if err != nil {
		return fmt.Errorf("failed to send email: %w", err)
	}

	return nil
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
