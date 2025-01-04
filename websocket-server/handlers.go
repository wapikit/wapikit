package websocket_server

import (
	"encoding/json"

	"github.com/gorilla/websocket"
)

// * these are event handlers for the events received from the client

func (s *WebSocketServer) handlePingEvent(messageId string, data json.RawMessage, connection *websocket.Conn) error {
	logger := s.app.Logger
	var eventData PingEventData
	if err := json.Unmarshal(data, &eventData); err != nil {
		logger.Error("error unmarshalling event data: %v", err)
		return err
	}
	ackBytes := NewAcknowledgementEvent(messageId, "Pong").toJson()
	err := s.sendMessageToClient(connection, ackBytes)
	if err != nil {
		logger.Error("error sending message to client: %v", err)
	}
	return err
}

func (server *WebSocketServer) handleMessageEvent(messageId string, data json.RawMessage, connection *websocket.Conn) error {
	logger := server.app.Logger
	var eventData MessageEventData

	// ! TODO: get the contact details
	// ! TODO: get the logged in organization

	if err := json.Unmarshal(data, &eventData); err != nil {
		logger.Error("error unmarshalling event data: %v", err)
		return err
	}
	ackBytes := NewAcknowledgementEvent(messageId, "Message received").toJson()
	err := server.sendMessageToClient(connection, ackBytes)
	if err != nil {
		logger.Error("error sending message to client: %v", err)
	}
	return err
}
