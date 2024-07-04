package websocket_server

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
	"github.com/labstack/echo/v4"
	"github.com/sarthakjdev/wapikit/internal"
	"github.com/sarthakjdev/wapikit/internal/interfaces"
)

type WebsocketConnectionToken struct {
	UserId string `json:"userId"`
	Token  string `json:"token"`
	Role   string `json:"role"`
}

type WebSocketServer struct {
	upgrader    websocket.Upgrader
	connections map[string]*websocket.Conn
	server      *echo.Echo
}

func newWebSocketServer(server *echo.Echo) *WebSocketServer {
	return &WebSocketServer{
		server: server,
		upgrader: websocket.Upgrader{
			CheckOrigin: func(r *http.Request) bool {
				return true
			},
			// EnableCompression: true,
		},
		connections: make(map[string]*websocket.Conn),
	}
}

func (s *WebSocketServer) handleWebSocket(c echo.Context) error {
	token := c.QueryParam("token")
	if token == "" {
		return echo.NewHTTPError(http.StatusUnauthorized, "Missing token")
	}

	// Upgrade to WebSocket connection
	ws, err := s.upgrader.Upgrade(c.Response().Writer, c.Request(), nil)
	fmt.Println("Upgraded to websocket connection")
	fmt.Println("error", err)
	if err != nil {
		log.Printf("Error connecting websocket: %v", err)
		return err
	}
	defer ws.Close()

	// Store connection
	connectionId := internal.GenerateUniqueId() // You'll need a function to generate unique IDs
	fmt.Println("connectionId: ", connectionId)
	fmt.Println("connectionId:", s.connections)
	s.connections[connectionId] = ws
	// defer delete(s.connections, connectionID)

	// Create a dedicated channel for receiving messages from this connection
	messageChan := make(chan []byte)
	go func() {
		fmt.Println("new channel")
		for {
			mt, messageData, err := ws.ReadMessage()
			fmt.Println("message type:", mt)
			if err != nil {
				// Signal that the connection is closed
				s.server.Logger.Error("closing channel: %v", err)
				close(messageChan)
				return
			}
			fmt.Printf("message received: %v", string(messageData))
			messageChan <- messageData
		}
	}()

	// Message processing loop
	for messageData := range messageChan {
		fmt.Println("messageData:", string(messageData))

		event := new(WebsocketEvent)
		if err := json.Unmarshal(messageData, &event); err != nil {
			fmt.Printf("Error unmarshalling message: %v\n", err)
			// Send an error message to the client (optional)
			s.sendMessageToClient(ws, []byte(`{"error": "Invalid message format"}`))
			continue
		}

		switch event.EventName {
		case WebsocketEventTypePing:
			if err := s.handlePingEvent(event.MessageId, event.Data, ws); err != nil {
				fmt.Println("Error handling ping: %v", err)
			}
		// ... (add cases for other event types) ...
		default:
			s.server.Logger.Warnf("Unknown WebSocket event: %s", event.EventName)
		}
	}

	fmt.Println("WebSocket connection closed")
	return nil
}

func (s *WebSocketServer) handlePingEvent(messageId string, data json.RawMessage, connection *websocket.Conn) error {

	fmt.Println("got a ping event")

	var eventData PingEventData
	if err := json.Unmarshal(data, &eventData); err != nil {
		log.Printf("Error unmarshalling event data: %v", err)
		return err
	}

	

	ackBytes := NewAcknowledgementEvent(messageId, "Pong").toJson()
	fmt.Println("message to send is", string(ackBytes))
	err := s.sendMessageToClient(connection, ackBytes)
	if err != nil {
		fmt.Println("error sending message to client: %v", err)
	}
	return err
}

func (ws *WebSocketServer) broadcastToAll(message []byte) {
	for _, conn := range ws.connections {
		err := ws.sendMessageToClient(conn, message)
		if err != nil {
			// Handle error (e.g., log, remove closed connection)
			fmt.Println("error sending message to client: %v", err)
		}
	}
}

func (ws *WebSocketServer) sendMessageToClient(conn *websocket.Conn, message []byte) error {
	err := conn.WriteMessage(websocket.TextMessage, message)
	if err != nil {
		// Handle error (e.g., log, remove closed connection)
		fmt.Println("Error sending message to client: %v", err)
		conn.Close()
		delete(ws.connections, conn.RemoteAddr().String()) // Cleanup
	}

	return err
}

func InitWebsocketServer(app *interfaces.App) *WebSocketServer {
	logger := app.Logger
	koa := app.Koa
	logger.Info("initializing websocket server")
	echoServer := echo.New()
	websocketServer := newWebSocketServer(echoServer)

	// corsOrigins := []string{}
	// corsOrigins = append(corsOrigins, "http://localhost:3000")

	websocketServer.server.GET("/ws", websocketServer.handleWebSocket)
	websocketServerAddress := koa.String("websocket_server_address")

	if websocketServerAddress == "" {
		websocketServerAddress = "localhost:5001"
	}

	func() {
		logger.Info("starting Websocket server server on %s", websocketServerAddress, nil) // Add a placeholder value as the final argument
		if err := websocketServer.server.Start("127.0.0.1:8081"); err != nil {
			if errors.Is(err, http.ErrServerClosed) {
				fmt.Println("websocket server shut down")
			} else {
				logger.Error("error starting HTTP server: %v", err)
			}
		}
	}()

	return websocketServer
}
