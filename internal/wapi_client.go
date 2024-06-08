package internal

import (
	cloud_client "github.com/sarthakjdev/wapi.go/pkg/client"
)

var client *cloud_client.Client

func GetWapiCloudClient(phoneNumberId, whatsappAccountId, webhookSecret, apiAccessToken string) *cloud_client.Client {

	if client != nil {
		return client
	}

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
