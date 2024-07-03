package services

import (
	"net/http"

	"github.com/golang-jwt/jwt"
	"github.com/labstack/echo/v4"
	"github.com/sarthakjdev/wapikit/database"
	"github.com/sarthakjdev/wapikit/internal/api_types"
	"github.com/sarthakjdev/wapikit/internal/interfaces"

	. "github.com/go-jet/jet/v2/postgres"
	"github.com/sarthakjdev/wapikit/.db-generated/model"
	table "github.com/sarthakjdev/wapikit/.db-generated/table"
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

func isAuthorized(role api_types.UserRoleEnum, routerPermissionLevel api_types.UserRoleEnum) bool {
	switch role {
	case api_types.Owner:
		return true
	case api_types.Admin:
		return routerPermissionLevel == api_types.Admin || routerPermissionLevel == api_types.Owner
	case api_types.Member:
		return routerPermissionLevel == api_types.Member
	default:
		return false
	}
}

func noAuthContextInjectionMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(ctx echo.Context) error {
		app := ctx.Get("app").(*interfaces.App)
		return next(interfaces.ContextWithoutSession{
			Context: ctx,
			App:     *app,
		})
	}

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
						AssignedRoles []model.RoleAssignment
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
				return echo.NewHTTPError(echo.ErrUnauthorized.Code, "Unauthorized access")
			}

			// ! TODO: fetch the integrations and enabled integration for the users and feed the booleans flags to the context

			app.Logger.Info("organization_id: ", organizationId)
			app.Logger.Info("user_id: ", user.User.UniqueId.String())

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
				app.Logger.Info("organization_id: ", org.Organization.UniqueId.String())
				if org.Organization.UniqueId.String() == organizationId {
					app.Logger.Info("org: ", org)
					// confirm the role access here
					// metadata := ctx.Get("routeMetaData").(interfaces.RouteMetaData)
					// app.Logger.Info("metadata: ", metadata)
					if isAuthorized(api_types.UserRoleEnum(org.MemberDetails.AccessLevel), api_types.Admin) {
						app.Logger.Info("is auth apporved")
						return next(interfaces.ContextWithSession{
							Context: ctx,
							App:     *app,
							Session: interfaces.ContextSession{
								Token: authToken,
								User: interfaces.ContextUser{
									UniqueId:       user.User.UniqueId.String(),
									Username:       user.User.Username,
									Email:          user.User.Email,
									Role:           api_types.UserRoleEnum(org.MemberDetails.AccessLevel),
									Name:           user.User.Name,
									OrganizationId: org.Organization.UniqueId.String(),
								},
							},
						})
					}
				} else {
					app.Logger.Info("organization not matched")
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
		routerMetaData := context.Get("routeMetaData").(interfaces.RouteMetaData)
		rateLimitConfig := routerMetaData.RateLimitConfig
		app.Logger.Info("rate limit config: ", rateLimitConfig)

		// rate limit the request
		// return error if rate limit is exceeded
		return next(context)
	}
}

func (service *BaseService) InjectRouterMetaData(routeMeta interfaces.RouteMetaData) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			c.Set("routeMetaData", routeMeta)
			return next(c)
		}
	}
}

// Register function now uses the Routes field
func (service *BaseService) Register(server *echo.Echo) {
	for _, route := range service.Routes {
		// server.Use(service.InjectRouterMetaData(route.MetaData))
		// handler := interfaces.CustomHandler(func(c interfaces.CustomContext) error {
		// Store metadata in context
		// 	c.Set("routeMetaData", route.MetaData)

		// Call the original handler
		// 	return route.Handler(c)
		// }).Handle

		handler := route.Handler.Handle
		// ! TODO: check meta thing here
		if route.IsAuthorizationRequired {
			handler = authMiddleware(handler)
		} else {
			handler = noAuthContextInjectionMiddleware(handler)
		}

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
