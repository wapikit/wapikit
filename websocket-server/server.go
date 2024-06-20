package websocket_server

import (
	"encoding/json"
	"errors"
	"fmt"
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
		server:      server,
		upgrader:    websocket.Upgrader{},
		connections: make(map[string]*websocket.Conn),
	}
}

func (s *WebSocketServer) handleWebSocket(c echo.Context) error {
	// Extract token from query parameter
	token := c.QueryParam("token")
	if token == "" {
		return echo.NewHTTPError(http.StatusUnauthorized, "Missing token")
	}
	// ! TODO: validate the token here

	// _, err := s.auth.ValidateJwt(token)
	// if err != nil {
	// 	return echo.NewHTTPError(http.StatusUnauthorized, "Invalid token")
	// }

	// Upgrade to WebSocket connection
	ws, err := s.upgrader.Upgrade(c.Response(), c.Request(), nil)
	if err != nil {
		return err
	}
	defer ws.Close()
	// Store connection
	connectionID := internal.GenerateUniqueId() // You'll need a function to generate unique IDs
	s.connections[connectionID] = ws
	defer delete(s.connections, connectionID)

	// Start reading messages from the client
	for {
		_, messageData, err := ws.ReadMessage()
		if err != nil {
			// handle error or connection close
			break
		}

		var event WebsocketEvent
		if err := json.Unmarshal(messageData, &event); err != nil {
			// handle invalid JSON
			continue
		}

		switch event.EventName {

		case WebsocketEventTypePing:
			// send back the pong
			ackBytes := NewAcknowledgementEvent(event.MessageId, "Pong!!").toJson()
			s.sendMessageToClient(ws, ackBytes)

		case "MessageEvent":
			// Handle message event
		case "NotificationReadEvent":
			// Handle notification read event
		// ... other event types ...
		default:
			// Handle unknown event type
		}
	}
	return nil
}

func (ws *WebSocketServer) broadcastToAll(message []byte) {
	for _, conn := range ws.connections {
		ws.sendMessageToClient(conn, message)
	}
}

func (ws *WebSocketServer) sendMessageToClient(conn *websocket.Conn, message []byte) {
	if err := conn.WriteMessage(websocket.TextMessage, message); err != nil {
		// Handle error (e.g., log, remove closed connection)
		ws.server.Logger.Errorf("Error sending message to client: %v", err)
		conn.Close()
		delete(ws.connections, conn.RemoteAddr().String()) // Cleanup
	}
}

func InitWebsocketServer(app *interfaces.App) *WebSocketServer {
	logger := app.Logger
	koa := app.Koa
	logger.Info("initializing websocket server")
	echoServer := echo.New()
	websocketServer := newWebSocketServer(echoServer)

	websocketServer.server.GET("/ws", websocketServer.handleWebSocket)

	websocketServerAddress := koa.String("websocket_server_address")

	if websocketServerAddress == "" {
		websocketServerAddress = "localhost:5001"
	}

	func() {
		logger.Info("starting Websocket server server on %s", websocketServerAddress, nil) // Add a placeholder value as the final argument
		if err := websocketServer.server.Start(":8081"); err != nil {
			if errors.Is(err, http.ErrServerClosed) {
				fmt.Println("websocket server shut down")
			} else {
				logger.Error("error starting HTTP server: %v", err)
			}
		}
	}()

	return websocketServer
}
