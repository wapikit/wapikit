package services

import (
	"fmt"
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
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("there was an error")
			}
			return []byte(app.Koa.String("jwt_secret")), nil
		})

		if err != nil {
			return echo.NewHTTPError(echo.ErrUnauthorized.Code, "Unauthorized access")
		}

		if parsedPayload.Valid {
			castedPayload := parsedPayload.Claims.(interfaces.JwtPayload)
			// * query the database to get the user
			userQuery := SELECT(
				table.User.AllColumns,
				table.Organization.AllColumns,
				table.OrganizationMember.AllColumns,
				table.RoleAssignment.AllColumns,
			).
				FROM(
					table.User.
						LEFT_JOIN(table.Organization, table.User.UniqueId.EQ(table.OrganizationMember.UserId)).
						LEFT_JOIN(table.OrganizationMember, table.OrganizationMember.OrganizationId.EQ(table.Organization.UniqueId).AND(table.OrganizationMember.UserId.EQ(table.User.UniqueId))).
						LEFT_JOIN(table.RoleAssignment, table.RoleAssignment.OrganizationMemberId.EQ(table.OrganizationMember.UniqueId)),
				).
				WHERE(
					table.User.Email.EQ(String(castedPayload.ContextUser.Email)).
						AND(
							table.User.UniqueId.EQ(String(castedPayload.ContextUser.UniqueId)),
						),
				).LIMIT(1)

			type UserWithOrgDetails struct {
				User          model.User `json:"-,inline"`
				Organizations []struct {
					Organization struct {
						model.Organization `json:"-,inline"`
						MemberDetails      model.OrganizationMember `json:"member_details"`
					}
					AssignedRoles []model.RoleAssignment `json:"assigned_roles"`
				} `json:"organizations"`
			}

			user := UserWithOrgDetails{}
			userQuery.Query(database.GetDbInstance(), &user)

			app.Logger.Info("user: ", user)

			if user.User.UniqueId.String() == "" || user.User.Status != "active" {
				return echo.NewHTTPError(echo.ErrUnauthorized.Code, "Unauthorized access")
			}

			for _, org := range user.Organizations {
				if org.Organization.UniqueId.String() == castedPayload.ContextUser.OrganizationId {
					// confirm the role access here

					metadata := ctx.Get("routeMeatData").(interfaces.RouteMetaData)
					if isAuthorized(api_types.UserRoleEnum(org.Organization.MemberDetails.AccessLevel), metadata.PermissionRoleLevel) {
						return next(interfaces.CustomContext{
							Context: ctx,
							App:     *app,
							Session: interfaces.ContextSession{
								Token: authToken,
								User: interfaces.ContextUser{
									UniqueId: user.User.UniqueId.String(),
									Username: user.User.Username,
									Email:    user.User.Email,
									Role:     api_types.UserRoleEnum(org.Organization.MemberDetails.AccessLevel),
									Name:     user.User.Name,
								},
							},
						})
					} else {
						return echo.NewHTTPError(echo.ErrUnauthorized.Code, "Unauthorized access")
					}

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

func (service *BaseService) InjectRouterMetaData(route interfaces.RouteMetaData) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			c.Set("routeMetaData", route)
			return next(c)
		}
	}
}

// Register function now uses the Routes field
func (service *BaseService) Register(server *echo.Echo) {
	for _, route := range service.Routes {
		server.Use(service.InjectRouterMetaData(route.MetaData))
		handler := interfaces.CustomHandler(route.Handler).Handle
		if route.IsAuthorizationRequired {
			handler = authMiddleware(interfaces.CustomHandler(route.Handler).Handle)
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
