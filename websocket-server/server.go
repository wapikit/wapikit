package websocket_server

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"sync"

	. "github.com/go-jet/jet/v2/postgres"
	"github.com/golang-jwt/jwt"
	"github.com/gorilla/websocket"
	"github.com/labstack/echo/v4"
	"github.com/wapikit/wapikit/.db-generated/model"
	table "github.com/wapikit/wapikit/.db-generated/table"

	"github.com/wapikit/wapikit/internal/interfaces"
)

// ! TODO:
// ! 1. we must be able to broadcast a message to all the connected clients of a organization
// ! 2. we must be able to send a message to a specific client
// ! 3. there must be a retry mechanism for sending message if in case the

type WebsocketConnectionData struct {
	UserId         string                        `json:"userId"`
	Token          string                        `json:"token"`
	AccessLevel    model.UserPermissionLevelEnum `json:"access_level"`
	Connection     *websocket.Conn               `json:"connection"`
	OrganizationId string                        `json:"organizationId"`
	Email          string                        `json:"email"`
	Username       string                        `json:"username"`
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

		userQuery.QueryContext(ctx.Request().Context(), app.Db, &user)

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
				accessLevel := model.UserPermissionLevelEnum(org.MemberDetails.AccessLevel)
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

	if err != nil {
		logger.Error("Error connecting websocket: %v", err.Error(), nil)
		return err
	}
	defer ws.Close()

	// * Store connection data
	connectionData.Connection = ws
	server.connections[connectionData.UserId] = connectionData
	defer delete(server.connections, connectionData.UserId)

	// * Create a dedicated channel for receiving websocket events from this connection
	websocketEventChannel := make(chan []byte)
	go func() {
		logger.Info("new channel created for event reception")
		for {
			_, eventInBinaryFormat, err := ws.ReadMessage()

			if err != nil {
				close(websocketEventChannel)
				return
			}

			var websocketEvent map[string]interface{}
			err = json.Unmarshal(eventInBinaryFormat, &websocketEvent)
			if err != nil {
				logger.Error("error decoding binary message: %v", err.Error(), nil)
				continue
			}

			websocketEventChannel <- eventInBinaryFormat
		}
	}()

	// Message processing loop
	for websocketEventData := range websocketEventChannel {
		event := new(WebsocketEvent)
		if err := json.Unmarshal(websocketEventData, &event); err != nil {
			logger.Error("error unmarshalling message: %v\n", err)
			// Send an error message to the client (optional)
			server.sendWebsocketEvent(ws, []byte(`{"error": "Invalid message format"}`))
			continue
		}

		switch event.EventName {
		case WebsocketEventTypePing:
			if err := server.handlePingEvent(event.EventId, event.Data, ws); err != nil {
				logger.Error("error handling ping: %v", err.Error(), nil)
			}
		case WebsocketEventTypeMessage:
			// ! TODO: user from the frontend has sent a new message to a contact

		default:
			logger.Warn("Unknown WebSocket event: %s", event.EventName)
		}
	}

	logger.Info("WebSocket connection closed")
	return nil
}

func (ws *WebSocketServer) broadcastToAll(message []byte) {
	logger := ws.app.Logger
	for _, conn := range ws.connections {
		err := ws.sendWebsocketEvent(conn.Connection, message)
		if err != nil {
			// Handle error (e.g., log, remove closed connection)
			logger.Info("error sending message to client: %v", err.Error(), nil)
		}
	}
}

func (ws *WebSocketServer) sendWebsocketEvent(conn *websocket.Conn, eventBytes []byte) error {

	var buffer bytes.Buffer

	buffer.Write(eventBytes)

	// ! TODO: implement a retry mechanism to send the message to the client, also as we know every message will be acknowledged, so we can wait for the acknowledgment and then retry if error

	logger := ws.app.Logger
	err := conn.WriteMessage(websocket.BinaryMessage, buffer.Bytes())
	if err != nil {
		logger.Error("error sending websocket event to client: %v", err)
		conn.Close()
		delete(ws.connections, conn.RemoteAddr().String()) // Cleanup
	}

	return err
}

func InitWebsocketServer(app *interfaces.App, wg *sync.WaitGroup) *WebSocketServer {
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

	go func() {
		logger.Info("starting Websocket server server on %s", websocketServerAddress, nil) // Add a placeholder value as the final argument
		if err := websocketServer.server.Start("127.0.0.1:8081"); err != nil {
			if errors.Is(err, http.ErrServerClosed) {
				logger.Info("websocket server shut down")
			} else {
				logger.Error("error starting HTTP server: %v", err.Error(), nil)
			}
		}
	}()

	go func() {
		defer wg.Done()
		websocketServer.HandleApiServerEvents(context.Background(), *app)
	}()

	fmt.Println("Websocket server started on: ", websocketServerAddress)

	return websocketServer
}
