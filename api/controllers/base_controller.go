package controller

import (
	"context"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/golang-jwt/jwt"
	"github.com/labstack/echo/v4"
	wapi "github.com/wapikit/wapi.go/pkg/client"
	"github.com/wapikit/wapikit/interfaces"
	"github.com/wapikit/wapikit/internal/api_types"
	"github.com/wapikit/wapikit/internal/services/ai_service"
	"github.com/wapikit/wapikit/internal/services/encryption_service"
	cache_service "github.com/wapikit/wapikit/internal/services/redis_service"

	. "github.com/go-jet/jet/v2/postgres"
	"github.com/wapikit/wapikit/.db-generated/model"
	table "github.com/wapikit/wapikit/.db-generated/table"
)

type BaseController struct {
	Name        string `json:"name"`
	RestApiPath string `json:"rest_api_path"`
	Routes      []interfaces.Route
}

func (s *BaseController) GetControllerName() string {
	return s.Name
}

func (s *BaseController) GetRoutes() []interfaces.Route {
	return s.Routes
}

func (s *BaseController) GetRestApiPath() string {
	return s.RestApiPath
}

func _noAuthContextInjectionMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(ctx echo.Context) error {
		app := ctx.Get("app").(*interfaces.App)
		return next(interfaces.ContextWithoutSession{
			Context: ctx,
			App:     *app,
		})
	}
}

func _isPermissionInList(requiredPermission api_types.RolePermissionEnum, userPermissions []api_types.RolePermissionEnum) bool {
	for _, permission := range userPermissions {
		if permission == requiredPermission {
			return true
		}
	}
	return false
}

func authMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(ctx echo.Context) error {
		app := ctx.Get("app").(*interfaces.App)
		headers := ctx.Request().Header
		authToken := headers.Get("x-access-token")

		if authToken == "" {
			return echo.NewHTTPError(echo.ErrUnauthorized.Code, "Unauthorized access")
		}
		// verify the jwt token
		parsedPayload, err := jwt.Parse(authToken, func(token *jwt.Token) (interface{}, error) {
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
			return echo.NewHTTPError(echo.ErrUnauthorized.Code, "Unauthorized access")
		}

		if parsedPayload.Valid {
			castedPayload := parsedPayload.Claims.(jwt.MapClaims)
			type UserWithOrgDetails struct {
				model.User
				Organizations []struct {
					model.Organization
					WhatsappBusinessAccount *model.WhatsappBusinessAccount
					MemberDetails           struct {
						model.OrganizationMember
						AssignedRoles []struct {
							model.RoleAssignment
							role model.OrganizationRole
						}
					}
				}
			}

			email := castedPayload["email"].(string)
			uniqueId := castedPayload["unique_id"].(string)
			organizationId := castedPayload["organization_id"].(string)

			if email == "" || uniqueId == "" {
				return echo.NewHTTPError(echo.ErrUnauthorized.Code, "Unauthorized access")
			}

			user := UserWithOrgDetails{}
			userQuery := SELECT(
				table.User.AllColumns,
				table.OrganizationMember.AllColumns,
				table.Organization.AllColumns,
				table.WhatsappBusinessAccount.AllColumns,
				table.RoleAssignment.AllColumns,
				table.OrganizationRole.AllColumns,
			).FROM(
				table.User.
					LEFT_JOIN(table.OrganizationMember, table.User.UniqueId.EQ(table.OrganizationMember.UserId)).
					LEFT_JOIN(table.Organization, table.Organization.UniqueId.EQ(table.OrganizationMember.OrganizationId)).
					LEFT_JOIN(table.WhatsappBusinessAccount, table.WhatsappBusinessAccount.OrganizationId.EQ(table.Organization.UniqueId)).
					LEFT_JOIN(table.RoleAssignment, table.OrganizationMember.UniqueId.EQ(table.RoleAssignment.OrganizationMemberId)).
					LEFT_JOIN(table.OrganizationRole, table.RoleAssignment.OrganizationRoleId.EQ(table.OrganizationRole.UniqueId)),
			).WHERE(
				table.User.Email.EQ(String(email)),
			)

			userQuery.QueryContext(ctx.Request().Context(), app.Db, &user)

			if user.User.UniqueId.String() == "" || user.User.Status != model.UserAccountStatusEnum_Active {
				app.Logger.Info("user not found or inactive")
				return echo.NewHTTPError(echo.ErrUnauthorized.Code, "Unauthorized access")
			}

			// ! TODO: fetch the integrations and enabled integration for the users and feed the booleans flags to the context

			if organizationId == "" {
				return next(interfaces.ContextWithSession{
					Context: ctx,
					App:     *app,
					Session: interfaces.ContextSession{
						Token: authToken,
						User: interfaces.ContextUser{
							UniqueId: user.User.UniqueId.String(),
							Username: user.User.Username,
							Email:    user.User.Email,
							Name:     user.User.Name,
						},
					},
				})
			}

			for _, org := range user.Organizations {
				if org.Organization.UniqueId.String() == organizationId {
					var routeMetadata interfaces.RouteMetaData
					metadata := ctx.Get("routeMetaData")
					if meta, ok := metadata.(interfaces.RouteMetaData); ok {
						routeMetadata = meta
					}

					var wapiClient *wapi.Client

					if org.WhatsappBusinessAccount != nil {
						wapiClient = wapi.New(&wapi.ClientConfig{
							BusinessAccountId: org.WhatsappBusinessAccount.AccountId,
							ApiAccessToken:    org.WhatsappBusinessAccount.AccessToken,
							WebhookSecret:     org.WhatsappBusinessAccount.WebhookSecret,
						})

						app.WapiClient = wapiClient
					}

					encryptionService := encryption_service.NewEncryptionService(&app.Logger, app.Koa.String("app.encryption_key"))
					app.EncryptionService = encryptionService

					if org.IsAiEnabled {
						// * initialize AI service
						ai_service := ai_service.NewAiService(&app.Logger, app.Redis, app.Db, org.AiApiKey)
						app.AiService = ai_service
					}

					// create a set of all permissions this user has
					if org.MemberDetails.AccessLevel == model.UserPermissionLevelEnum_Owner ||
						routeMetadata.RequiredPermission == nil ||
						len(routeMetadata.RequiredPermission) == 0 {
						return next(interfaces.ContextWithSession{
							Context: ctx,
							App:     *app,
							Session: interfaces.ContextSession{
								Token: authToken,
								User: interfaces.ContextUser{
									UniqueId:       user.User.UniqueId.String(),
									Username:       user.User.Username,
									Email:          user.User.Email,
									Role:           api_types.UserPermissionLevelEnum(org.MemberDetails.AccessLevel),
									Name:           user.User.Name,
									OrganizationId: org.Organization.UniqueId.String(),
								},
							},
						})
					}

					userCurrentOrgPermissions := []api_types.RolePermissionEnum{}
					permissionSet := make(map[api_types.RolePermissionEnum]struct{})

					// * extracting out mutually exclusive permissions from the assigned roles
					for _, roleAssignment := range org.MemberDetails.AssignedRoles {
						permissionArray := strings.Split(roleAssignment.role.Permissions, ",")
						for _, permission := range permissionArray {
							perm := api_types.RolePermissionEnum(permission)
							if _, exists := permissionSet[perm]; !exists {
								permissionSet[perm] = struct{}{}
								userCurrentOrgPermissions = append(userCurrentOrgPermissions, perm)
							}
						}
					}

					// * now check if user  has required permission in the list of permissions it has
					for _, requiredPermission := range routeMetadata.RequiredPermission {
						if !_isPermissionInList(requiredPermission, userCurrentOrgPermissions) {
							return echo.NewHTTPError(echo.ErrUnauthorized.Code, "You are not authorized to access this resource.")
						}
					}

					return next(interfaces.ContextWithSession{
						Context: ctx,
						App:     *app,
						Session: interfaces.ContextSession{
							Token: authToken,
							User: interfaces.ContextUser{
								UniqueId:       user.User.UniqueId.String(),
								Username:       user.User.Username,
								Email:          user.User.Email,
								Role:           api_types.UserPermissionLevelEnum(org.MemberDetails.AccessLevel),
								Name:           user.User.Name,
								OrganizationId: org.Organization.UniqueId.String(),
							},
						},
					})
				}
			}

			return echo.NewHTTPError(echo.ErrUnauthorized.Code, "Unauthorized access")
		} else {
			return echo.NewHTTPError(echo.ErrUnauthorized.Code, "Unauthorized access")
		}
	}
}

func rateLimiter(next echo.HandlerFunc) echo.HandlerFunc {
	return func(context echo.Context) error {
		app := context.Get("app").(*interfaces.App)
		routerMetaData := context.Get("routeMetaData")
		redisService := app.Redis

		if routerMetaData == nil {
			// Skip rate limiting if no metadata is found
			return next(context)
		}

		routeMetaData, ok := routerMetaData.(interfaces.RouteMetaData)
		if !ok {
			// Skip rate limiting if metadata is not of the correct type
			return next(context)
		}

		rateLimitConfig := routeMetaData.RateLimitConfig

		if rateLimitConfig.MaxRequests <= 0 || rateLimitConfig.WindowTimeInMs <= 0 {
			// Skip rate limiting if configuration is invalid
			return next(context)
		}

		clientIP := context.RealIP()
		path := context.Path()

		redisKey := redisService.ComputeRateLimitKey(clientIP, path)
		windowDuration := time.Duration(rateLimitConfig.WindowTimeInMs) * time.Millisecond

		finalMaxRequestsAllowed := rateLimitConfig.MaxRequests

		if app.Constants.IsCommunityEdition {
			// do nothing same as before
		} else if app.Constants.IsCloudEdition {
			// ! get the user subscription
			isUserOnProOrScalePlan := false
			isUserOnEnterprisePlan := false
			if isUserOnProOrScalePlan {
				finalMaxRequestsAllowed = int(float64(rateLimitConfig.MaxRequests) * 1.5)
			} else if isUserOnEnterprisePlan {
				finalMaxRequestsAllowed = rateLimitConfig.MaxRequests * 3
			}
		} else {
			// ! this means it is a instance of enterprise self hosted licensed edition
			finalMaxRequestsAllowed = rateLimitConfig.MaxRequests * 3
		}

		allowed, remaining, reset, err := enforceRateLimit(redisService, redisKey, finalMaxRequestsAllowed, windowDuration)
		if err != nil {
			app.Logger.Error("Error in rate limiter: ", err)
			return context.JSON(500, map[string]string{
				"error": "Internal server error",
			})
		}

		if !allowed {
			return context.JSON(429, map[string]interface{}{
				"error":     "Rate limit exceeded",
				"remaining": remaining,
				"reset":     reset,
			})
		}

		return next(context)
	}
}

func getUserSubscriptionDetails() {
	// ! get the user subscription
}

// enforceRateLimit checks and updates the rate limit state in Redis
func enforceRateLimit(redisClient *cache_service.RedisClient, key string, limit int, window time.Duration) (bool, int, int64, error) {
	ctx := context.Background()

	currentTime := time.Now().Unix()
	expiry := int64(window.Seconds())

	pipe := redisClient.Pipeline()
	currentCountCmd := pipe.Get(ctx, key)
	pipe.Exec(ctx)

	currentCountStr, err := currentCountCmd.Result()
	if err != nil && err != redis.Nil {
		return false, 0, 0, err
	}

	currentCount := 0
	if err == nil {
		currentCount, err = strconv.Atoi(currentCountStr)
		if err != nil {
			return false, 0, 0, err
		}
	}

	// Check if the request is within the limit
	if currentCount >= limit {
		resetTimeCmd := redisClient.TTL(ctx, key)
		resetTime, err := resetTimeCmd.Result()
		if err != nil {
			return false, 0, 0, err
		}
		return false, limit - currentCount, currentTime + int64(resetTime.Seconds()), nil
	}

	// Increment the request count and set expiry if not set
	pipe.Incr(ctx, key)
	pipe.Expire(ctx, key, window)
	_, err = pipe.Exec(ctx)
	if err != nil {
		return false, 0, 0, err
	}

	remaining := limit - (currentCount + 1)
	return true, remaining, currentTime + expiry, nil
}

func _injectRouteMetaData(routeMeta interfaces.RouteMetaData) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			fmt.Println("Injecting route metadata", routeMeta)
			// Set the specific route metadata in the context
			c.Set("routeMetaData", routeMeta)
			return next(c)
		}
	}
}

func (controller *BaseController) Register(server *echo.Echo) {
	for _, route := range controller.Routes {
		// Create handler and inject route-specific metadata
		handler := route.Handler.Handle

		// Apply authorization middleware if required
		if route.IsAuthorizationRequired {
			handler = authMiddleware(handler)
		} else {
			handler = _noAuthContextInjectionMiddleware(handler)
		}

		handler = rateLimiter(handler)
		handler = _injectRouteMetaData(route.MetaData)(handler)

		// Register the route with the appropriate HTTP method
		switch route.Method {
		case http.MethodGet:
			server.GET(route.Path, handler)
		case http.MethodPost:
			server.POST(route.Path, handler)
		case http.MethodPut:
			server.PUT(route.Path, handler)
		case http.MethodDelete:
			server.DELETE(route.Path, handler)
		}
	}
}
