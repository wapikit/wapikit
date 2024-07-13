package websocket_server

import (
	"context"
	"encoding/json"
	"errors"
	"log"
	"net/http"

	. "github.com/go-jet/jet/v2/postgres"
	"github.com/golang-jwt/jwt"
	"github.com/gorilla/websocket"
	"github.com/labstack/echo/v4"
	"github.com/sarthakjdev/wapikit/.db-generated/model"
	table "github.com/sarthakjdev/wapikit/.db-generated/table"
	"github.com/sarthakjdev/wapikit/database"
	"github.com/sarthakjdev/wapikit/internal/interfaces"
)

type WebsocketConnectionData struct {
	UserId         string                    `json:"userId"`
	Token          string                    `json:"token"`
	AccessLevel    model.UserPermissionLevel `json:"access_level"`
	Connection     *websocket.Conn           `json:"connection"`
	OrganizationId string                    `json:"organizationId"`
	Email          string                    `json:"email"`
	Username       string                    `json:"username"`
}

type WebSocketServer struct {
	upgrader    websocket.Upgrader
	connections map[string]*WebsocketConnectionData
	server      *echo.Echo
	app         interfaces.App
}

func newWebSocketServer(server *echo.Echo, app interfaces.App) *WebSocketServer {
	return &WebSocketServer{
		app:    app,
		server: server,
		upgrader: websocket.Upgrader{
			CheckOrigin: func(r *http.Request) bool {
				return true
			},
			// EnableCompression: true,
		},
		connections: make(map[string]*WebsocketConnectionData),
	}
}

func (server *WebSocketServer) authorizeConnectionRequest(ctx echo.Context) (*WebsocketConnectionData, error) {
	token := ctx.QueryParam("token")
	if token == "" {
		return nil, echo.NewHTTPError(http.StatusUnauthorized, "Missing token")
	}

	app := server.app

	parsedPayload, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		secretKey := app.Koa.String("app.jwt_secret")
		if secretKey == "" {
			app.Logger.Error("jwt secret key not configured")
			return "", nil
		}
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			echo.NewHTTPError(echo.ErrUnauthorized.Code, "Unauthorized access")
			return "", nil
		}
		return []byte(app.Koa.String("app.jwt_secret")), nil
	})

	if err != nil {
		return nil, echo.NewHTTPError(echo.ErrUnauthorized.Code, "Unauthorized access")
	}

	if parsedPayload.Valid {
		castedPayload := parsedPayload.Claims.(jwt.MapClaims)
		type UserWithOrgDetails struct {
			model.User
			Organizations []struct {
				model.Organization
				MemberDetails struct {
					model.OrganizationMember
					AssignedRoles []model.RoleAssignment
				}
			}
		}

		email := castedPayload["email"].(string)
		uniqueId := castedPayload["unique_id"].(string)
		organizationId := castedPayload["organization_id"].(string)

		if email == "" || uniqueId == "" {
			return nil, echo.NewHTTPError(echo.ErrUnauthorized.Code, "Unauthorized access")
		}

		user := UserWithOrgDetails{}
		userQuery := SELECT(
			table.User.AllColumns,
			table.OrganizationMember.AllColumns,
			table.Organization.AllColumns,
			table.RoleAssignment.AllColumns,
		).FROM(
			table.User.
				LEFT_JOIN(table.OrganizationMember, table.User.UniqueId.EQ(table.OrganizationMember.UserId)).
				LEFT_JOIN(table.Organization, table.Organization.UniqueId.EQ(table.OrganizationMember.OrganizationId)).
				LEFT_JOIN(table.RoleAssignment, table.OrganizationMember.UniqueId.EQ(table.RoleAssignment.OrganizationMemberId)),
		).WHERE(
			table.User.Email.EQ(String(email)),
		)

		userQuery.QueryContext(ctx.Request().Context(), database.GetDbInstance(), &user)

		if user.User.UniqueId.String() == "" || user.User.Status != model.UserAccountStatusEnum_Active {
			app.Logger.Info("user not found or inactive")
			return nil, echo.NewHTTPError(echo.ErrUnauthorized.Code, "Unauthorized access")
		}

		// ! TODO: fetch the integrations and enabled integration for the users and feed the booleans flags to the context

		if organizationId == "" {
			return nil, echo.NewHTTPError(echo.ErrUnauthorized.Code, "Unauthorized access")
		}

		for _, org := range user.Organizations {
			if org.Organization.UniqueId.String() == organizationId {
				accessLevel := model.UserPermissionLevel(org.MemberDetails.AccessLevel)
				connectionData := WebsocketConnectionData{
					UserId:         user.User.UniqueId.String(),
					Token:          token,
					AccessLevel:    accessLevel,
					OrganizationId: org.Organization.UniqueId.String(),
					Email:          user.User.Email,
					Username:       user.User.Username,
				}

				return &connectionData, nil
			}
		}

		return nil, echo.NewHTTPError(echo.ErrUnauthorized.Code, "Unauthorized access")
	} else {
		return nil, echo.NewHTTPError(echo.ErrUnauthorized.Code, "Unauthorized access")
	}
}

func (server *WebSocketServer) handleWebSocket(ctx echo.Context) error {
	token := ctx.QueryParam("token")
	if token == "" {
		return echo.NewHTTPError(http.StatusUnauthorized, "Missing token")
	}

	logger := server.app.Logger
	connectionData, err := server.authorizeConnectionRequest(ctx)

	if err != nil {
		return err
	}

	// Upgrade to WebSocket connection
	ws, err := server.upgrader.Upgrade(ctx.Response().Writer, ctx.Request(), nil)
	logger.Info("upgraded to websocket connection!!")
	// logger.Error("error", err.Error(), nil)
	if err != nil {
		log.Printf("Error connecting websocket: %v", err)
		return err
	}
	defer ws.Close()

	// ! TODO: create a subscription to the redis pubsub channel, to receive messages and broadcast them to concerned clients
	// * Store connection data
	connectionData.Connection = ws
	server.connections[connectionData.UserId] = connectionData
	defer delete(server.connections, connectionData.UserId)

	// * Create a dedicated channel for receiving messages from this connection
	messageChan := make(chan []byte)
	go func() {
		logger.Info("new channel created for message reception")
		for {
			_, messageData, err := ws.ReadMessage()
			if err != nil {
				// Signal that the connection is closed
				server.server.Logger.Error("closing channel: %v", err)
				close(messageChan)
				return
			}
			logger.Info("message received: %v", string(messageData))
			messageChan <- messageData
		}
	}()

	// Message processing loop
	for messageData := range messageChan {
		logger.Info("messageData:", string(messageData))

		event := new(WebsocketEvent)
		if err := json.Unmarshal(messageData, &event); err != nil {
			logger.Error("error unmarshalling message: %v\n", err)
			// Send an error message to the client (optional)
			server.sendMessageToClient(ws, []byte(`{"error": "Invalid message format"}`))
			continue
		}

		switch event.EventName {
		case WebsocketEventTypePing:
			if err := server.handlePingEvent(event.MessageId, event.Data, ws); err != nil {
				logger.Error("error handling ping: %v", err.Error())
			}
		default:
			logger.Warn("Unknown WebSocket event: %s", event.EventName)
		}
	}

	logger.Info("WebSocket connection closed")
	return nil
}

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

func (ws *WebSocketServer) broadcastToAll(message []byte) {
	logger := ws.app.Logger
	for _, conn := range ws.connections {
		err := ws.sendMessageToClient(conn.Connection, message)
		if err != nil {
			// Handle error (e.g., log, remove closed connection)
			logger.Info("error sending message to client: %v", err)
		}
	}
}

func (ws *WebSocketServer) sendMessageToClient(conn *websocket.Conn, message []byte) error {
	logger := ws.app.Logger
	err := conn.WriteMessage(websocket.TextMessage, message)
	if err != nil {
		// Handle error (e.g., log, remove closed connection)
		logger.Error("error sending message to client: %v", err)
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
	websocketServer := newWebSocketServer(echoServer, *app)

	websocketServer.server.GET("/ws", websocketServer.handleWebSocket)
	websocketServerAddress := koa.String("websocket_server_address")

	if websocketServerAddress == "" {
		websocketServerAddress = "localhost:5001"
	}

	func() {
		logger.Info("starting Websocket server server on %s", websocketServerAddress, nil) // Add a placeholder value as the final argument
		if err := websocketServer.server.Start("127.0.0.1:8081"); err != nil {
			if errors.Is(err, http.ErrServerClosed) {
				logger.Info("websocket server shut down")
			} else {
				logger.Error("error starting HTTP server: %v", err.Error(), nil)
			}
		}
	}()

	go HandleApiServerEvents(context.Background(), *app)
	return websocketServer
}
