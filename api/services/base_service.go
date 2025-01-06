package services

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/golang-jwt/jwt"
	"github.com/labstack/echo/v4"
	"github.com/wapikit/wapikit/internal/api_types"
	"github.com/wapikit/wapikit/internal/interfaces"

	. "github.com/go-jet/jet/v2/postgres"
	"github.com/wapikit/wapikit/.db-generated/model"
	table "github.com/wapikit/wapikit/.db-generated/table"
)

type BaseService struct {
	Name        string `json:"name"`
	RestApiPath string `json:"rest_api_path"`
	Routes      []interfaces.Route
}

func (s *BaseService) GetServiceName() string {
	return s.Name
}

func (s *BaseService) GetRoutes() []interfaces.Route {
	return s.Routes
}

func (s *BaseService) GetRestApiPath() string {
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
					MemberDetails struct {
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
				table.RoleAssignment.AllColumns,
				table.OrganizationRole.AllColumns,
			).FROM(
				table.User.
					LEFT_JOIN(table.OrganizationMember, table.User.UniqueId.EQ(table.OrganizationMember.UserId)).
					LEFT_JOIN(table.Organization, table.Organization.UniqueId.EQ(table.OrganizationMember.OrganizationId)).
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

		if routerMetaData == nil {
			// skip rate limiting if no metadata is found
		} else {
			routerMetaData, ok := routerMetaData.(interfaces.RouteMetaData)

			if !ok {
				// skip rate limiting if metadata is not of the correct type
			} else {
				rateLimitConfig := routerMetaData.RateLimitConfig
				// ! TODO: redis rate limit here
				app.Logger.Info("rateLimitConfig: ", rateLimitConfig)
			}

		}
		return next(context)
	}
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

func (service *BaseService) Register(server *echo.Echo) {
	for _, route := range service.Routes {
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
