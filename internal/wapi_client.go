package internal

import (
	cloud_client "github.com/sarthakjdev/wapi.go/pkg/client"
	"github.com/sarthakjdev/wapikit/internal/interfaces"
)

var client *cloud_client.Client

func GetWapiClient(app *interfaces.App) *cloud_client.Client {
	if client != nil {
		return client
	}

	config := app.Koa
	phoneNumberId := config.String("phoneNumberId")
	whatsappAccountId := config.String("whatsappAccountId")
	webhookSecret := config.String("webhookSecret")
	apiAccessToken := config.String("apiAccessToken")

	client, err := cloud_client.New(cloud_client.ClientConfig{
		PhoneNumberId:     phoneNumberId,
		ApiAccessToken:    apiAccessToken,
		BusinessAccountId: whatsappAccountId,
		WebhookSecret:     webhookSecret,
	})

	if err != nil {
		panic(err)
	}

	return client

}
