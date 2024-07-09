package auth_service

import (
	"net/http"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/sarthakjdev/wapikit/api/services"
	"github.com/sarthakjdev/wapikit/database"
	"github.com/sarthakjdev/wapikit/internal"
	"github.com/sarthakjdev/wapikit/internal/api_types"
	"github.com/sarthakjdev/wapikit/internal/interfaces"
	"golang.org/x/crypto/bcrypt"

	. "github.com/go-jet/jet/v2/postgres"
	"github.com/sarthakjdev/wapikit/.db-generated/model"
	table "github.com/sarthakjdev/wapikit/.db-generated/table"
)

type AuthService struct {
	services.BaseService `json:"-,inline"`
}

func NewAuthService() *AuthService {
	return &AuthService{
		BaseService: services.BaseService{
			Name:        "Auth Service",
			RestApiPath: "/api/auth",
			Routes: []interfaces.Route{
				{
					Path:                    "/api/auth/login",
					Method:                  http.MethodPost,
					Handler:                 interfaces.HandlerWithoutSession(handleSignIn),
					IsAuthorizationRequired: false,
				},
				{
					Path:                    "/api/auth/register",
					Method:                  http.MethodPost,
					Handler:                 interfaces.HandlerWithoutSession(handleUserRegistration),
					IsAuthorizationRequired: false,
				},
				{
					Path:                    "/api/auth/verify-email",
					Method:                  http.MethodPost,
					Handler:                 interfaces.HandlerWithoutSession(verifyEmailAndCreateAccount),
					IsAuthorizationRequired: false,
				},
				{
					Path:                    "/api/auth/api-keys",
					Method:                  http.MethodGet,
					Handler:                 interfaces.HandlerWithSession(getApiKeys),
					IsAuthorizationRequired: true,
					MetaData: interfaces.RouteMetaData{
						PermissionRoleLevel: api_types.Admin,
						RateLimitConfig: interfaces.RateLimitConfig{
							MaxRequests:    10,
							WindowTimeInMs: 1000 * 60 * 60, // 1 hour
						},
					},
				},
				{
					Path:                    "/api/auth/api-keys/regenerate",
					Method:                  http.MethodGet,
					Handler:                 interfaces.HandlerWithSession(regenerateApiKey),
					IsAuthorizationRequired: true,
					MetaData: interfaces.RouteMetaData{
						PermissionRoleLevel: api_types.Admin,
					},
				},
				{
					Path:                    "/api/auth/oauth",
					Method:                  http.MethodPost,
					Handler:                 interfaces.HandlerWithoutSession(handleLoginWithOAuth),
					IsAuthorizationRequired: false,
				},
				{
					Path:                    "/api/auth/switch",
					Method:                  http.MethodPost,
					Handler:                 interfaces.HandlerWithSession(switchOrganization),
					IsAuthorizationRequired: true,
				},
				{
					Path:                    "/api/auth/join-organization",
					Method:                  http.MethodPost,
					Handler:                 interfaces.HandlerWithSession(acceptOrganizationInvite),
					IsAuthorizationRequired: true,
					MetaData: interfaces.RouteMetaData{
						PermissionRoleLevel: api_types.Admin,
						RateLimitConfig: interfaces.RateLimitConfig{
							MaxRequests:    10,
							WindowTimeInMs: 1000 * 60 * 60, // 1 hour
						},
					},
				},
			},
		},
	}
}

func acceptOrganizationInvite(context interfaces.ContextWithSession) error {

	payload := new(api_types.JoinOrganizationJSONRequestBody)

	if err := context.Bind(payload); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	if payload.InviteSlug == nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Organization invite slug is required")
	}

	// get the user
	user := context.Session.User

	userUuid, err := uuid.Parse(user.UniqueId)

	if err != nil {
		return echo.NewHTTPError(http.StatusUnauthorized, "Unauthorized access")
	}

	// get the organization invite
	var invite model.OrganizationMemberInvite

	invitationQuery := SELECT(table.OrganizationMemberInvite.AllColumns).
		FROM(table.OrganizationMemberInvite).
		WHERE(
			table.OrganizationMemberInvite.Email.EQ(String(user.Email)).
				AND(table.OrganizationMemberInvite.Slug.EQ(String(*payload.InviteSlug)))).
		LIMIT(1)

	err = invitationQuery.QueryContext(context.Request().Context(), context.App.Db, &invite)

	if err != nil {
		if err.Error() == "qrm: no rows in result set" {
			return echo.NewHTTPError(http.StatusNotFound, "Organization invite not found")
		} else {
			context.App.Logger.Error("database query error", err.Error())
			return echo.NewHTTPError(http.StatusInternalServerError, "Something went wrong while processing your request.")
		}
	}

	if invite.Status != model.OrganizationInviteStatusEnum_Pending {
		return echo.NewHTTPError(http.StatusForbidden, "Invite already accepted")
	}

	// check if the user is already a member of the organization
	var existingMember model.OrganizationMember

	existingMemberQuery := SELECT(table.OrganizationMember.AllColumns).
		FROM(table.OrganizationMember).
		WHERE(
			table.OrganizationMember.UserId.EQ(UUID(userUuid)).
				AND(table.OrganizationMember.OrganizationId.EQ(UUID(invite.OrganizationId)))).
		LIMIT(1)

	err = existingMemberQuery.QueryContext(context.Request().Context(), context.App.Db, &existingMember)

	if err != nil {
		if err.Error() == "qrm: no rows in result set" {
			// do nothing
		}
	} else {
		return echo.NewHTTPError(http.StatusForbidden, "User already a member of the organization")
	}

	var insertedOrgMember model.OrganizationMember

	err = table.OrganizationMember.INSERT().MODEL(model.OrganizationMember{
		AccessLevel:    invite.AccessLevel,
		OrganizationId: invite.OrganizationId,
		UserId:         userUuid,
		InviteId:       &invite.UniqueId,
	}).QueryContext(context.Request().Context(), context.App.Db, &insertedOrgMember)

	if err != nil {
		context.App.Logger.Error("database query error", err.Error())
		return echo.NewHTTPError(http.StatusInternalServerError, "Something went wrong while processing your request.")
	}

	// create the token

	claims := &interfaces.JwtPayload{
		ContextUser: interfaces.ContextUser{
			Username:       context.Session.User.Username,
			Email:          context.Session.User.Email,
			Role:           api_types.UserRoleEnum(invite.AccessLevel),
			UniqueId:       context.Session.User.UniqueId,
			OrganizationId: invite.OrganizationId.String(),
			Name:           context.Session.User.Name,
		},
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Hour * 24 * 60).Unix(), // 60-day expiration
			Issuer:    "wapikit",
		},
	}

	token, err := jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString([]byte(context.App.Koa.String("app.jwt_secret")))
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Error generating token")
	}
	response := api_types.JoinOrganizationResponseBodySchema{
		Token: token,
	}
	return context.JSON(http.StatusOK, response)

}

func handleSignIn(context interfaces.ContextWithoutSession) error {
	payload := new(api_types.LoginRequestBodySchema)

	if err := context.Bind(payload); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	if payload.Username == "" || payload.Password == "" {
		return echo.NewHTTPError(echo.ErrBadRequest.Code, "Username / password is required")
	}

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

	user := UserWithOrgDetails{}
	stmt := SELECT(
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
		table.User.Username.EQ(String(payload.Username)).
			OR(table.User.Email.EQ(String(payload.Username))),
	)

	stmt.QueryContext(context.Request().Context(), database.GetDbInstance(), &user)
	context.App.Logger.Info("User details:", user)

	// if no user found then return 404
	if user.User.UniqueId.String() == "" || user.User.Password == nil {
		return echo.NewHTTPError(http.StatusNotFound, "Invalid email / password")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(*user.User.Password), []byte(payload.Password)); err != nil {
		return echo.NewHTTPError(http.StatusNotFound, "Invalid email / password")
	}

	var isOnboardingCompleted bool
	var organizationIdToLoginWith string
	var roleToLoginWith api_types.UserRoleEnum
	var claims *interfaces.JwtPayload

	// if no organization found, then simply return the user with a flag saying isOnboardingCompleted
	if len(user.Organizations) == 0 {
		isOnboardingCompleted = false
		// onboarding to be completed by the user
		// return the user with onboarding flag

		claims = &interfaces.JwtPayload{
			ContextUser: interfaces.ContextUser{
				Username: user.User.Username,
				Email:    user.User.Email,
				UniqueId: user.User.UniqueId.String(),
				Name:     user.User.Name,
			},
			StandardClaims: jwt.StandardClaims{
				ExpiresAt: time.Now().Add(time.Hour * 24 * 60).Unix(), // 60-day expiration
				Issuer:    "wapikit",
			},
		}

	} else {
		isOnboardingCompleted = true

		// check for the owner org
		for _, org := range user.Organizations {
			if org.MemberDetails.AccessLevel == model.UserPermissionLevel_Owner {
				organizationIdToLoginWith = org.Organization.UniqueId.String()
				roleToLoginWith = api_types.Owner
				break
			}
		}

		// check for the admin org
		if organizationIdToLoginWith == "" {
			// no owner org found, login with the org having the highest role
			// here if no owner org found then look for the lower roles too
			for _, org := range user.Organizations {
				if org.MemberDetails.AccessLevel == model.UserPermissionLevel_Admin {
					organizationIdToLoginWith = org.Organization.UniqueId.String()
					roleToLoginWith = api_types.Admin
					break
				}
			}
		}

		// else login with the first org
		if organizationIdToLoginWith == "" {
			organizationIdToLoginWith = user.Organizations[0].Organization.UniqueId.String()
			roleToLoginWith = api_types.Member
		}

		claims = &interfaces.JwtPayload{
			ContextUser: interfaces.ContextUser{
				Username:       user.User.Username,
				Email:          user.User.Email,
				Role:           api_types.UserRoleEnum(roleToLoginWith),
				UniqueId:       user.User.UniqueId.String(),
				OrganizationId: organizationIdToLoginWith,
				Name:           user.User.Name,
			},
			StandardClaims: jwt.StandardClaims{
				ExpiresAt: time.Now().Add(time.Hour * 24 * 60).Unix(), // 60-day expiration
				Issuer:    "wapikit",
			},
		}

	}

	//Create the token
	token, err := jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString([]byte(context.App.Koa.String("app.jwt_secret")))

	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Error generating token")
	}

	return context.JSON(http.StatusOK, api_types.LoginResponseBodySchema{
		IsOnboardingCompleted: isOnboardingCompleted,
		Token:                 token,
	})
}

func handleLoginWithOAuth(context interfaces.ContextWithoutSession) error {
	return nil
}

// this handler would validate the email and send an otp to it
func handleUserRegistration(context interfaces.ContextWithoutSession) error {

	payload := new(api_types.RegisterRequestBodySchema)

	if err := context.Bind(payload); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	if payload.Username == "" || payload.Email == "" || payload.Password == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "Username, Email and Password are required")
	}

	otp := "123456"

	if context.App.Constants.IsProduction {
		otp = internal.GenerateOtp()
	}

	cacheKey := internal.ComputeCacheKey("otp", payload.Email, "registration")

	err := internal.CacheData(cacheKey, otp, time.Minute*5)

	if err != nil {
		context.App.Logger.Error("Error caching otp", err.Error())
		return echo.NewHTTPError(http.StatusInternalServerError, "Something went wrong while processing your request.")
	}

	// ! TODO: send the otp to the email

	// return the response
	return context.JSON(http.StatusOK, api_types.RegisterRequestResponseBodySchema{
		IsOtpSent: true,
	})
}

func verifyEmailAndCreateAccount(context interfaces.ContextWithoutSession) error {
	payload := new(api_types.VerifyOtpJSONRequestBody)
	if err := context.Bind(payload); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	if payload.Username == "" || payload.Email == "" || payload.Password == "" || payload.Otp == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid details!!")
	}

	cacheKey := internal.ComputeCacheKey("otp", payload.Email, "registration")
	cachedOtp, err := internal.GetCachedData(cacheKey)
	if err != nil {
		context.App.Logger.Error("Error getting cached otp", err.Error())
		return echo.NewHTTPError(http.StatusInternalServerError, "Something went wrong while processing your request.")
	}

	if cachedOtp != payload.Otp {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid OTP")
	}

	// check if the user already exists
	var user model.User

	userQuery := SELECT(table.User.AllColumns).
		FROM(table.User).
		WHERE(table.User.Email.EQ(String(payload.Email))).
		LIMIT(1)

	err = userQuery.QueryContext(context.Request().Context(), context.App.Db, &user)

	if err != nil {
		if err.Error() == "qrm: no rows in result set" {
			// do nothing
		} else {
			context.App.Logger.Error("database query error", err.Error())
			return echo.NewHTTPError(http.StatusInternalServerError, "Something went wrong while processing your request.")
		}
	}

	if user.UniqueId.String() != "" {
		return echo.NewHTTPError(http.StatusBadRequest, "User already exists")
	}

	var invite model.OrganizationMemberInvite
	invitationQuery := SELECT(table.OrganizationMemberInvite.AllColumns).
		FROM(table.OrganizationMemberInvite).
		WHERE(
			table.OrganizationMemberInvite.Email.EQ(String(payload.Email)).
				AND(table.OrganizationMemberInvite.Slug.EQ(String(*payload.OrganizationInviteSlug)))).
		LIMIT(1)

	err = invitationQuery.QueryContext(context.Request().Context(), context.App.Db, &invite)

	if err != nil {
		if err.Error() == "qrm: no rows in result set" {
			// do  nothing just move on, we cant let the user not register if they do not have a valid invite
		} else {
			context.App.Logger.Error("database query error", err.Error())
			return echo.NewHTTPError(http.StatusInternalServerError, "Something went wrong while processing your request.")
		}
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(payload.Password), bcrypt.DefaultCost)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Error hashing password")
	}

	passwordString := string(hashedPassword)

	var insertedUser model.User
	var insertedOrgMember model.OrganizationMember

	err = table.User.INSERT().MODEL(model.User{
		Username: payload.Username,
		Email:    payload.Email,
		Password: &passwordString,
		Name:     payload.Name,
		Status:   model.UserAccountStatusEnum_Active,
	}).RETURNING(table.User.AllColumns).QueryContext(context.Request().Context(), context.App.Db, &insertedUser)

	if err != nil {
		context.App.Logger.Error("database query error", err.Error())
		return echo.NewHTTPError(http.StatusInternalServerError, "Something went wrong while processing your request.")
	}

	if invite.UniqueId.String() == "" {
		err = table.OrganizationMember.INSERT().MODEL(model.OrganizationMember{
			AccessLevel:    invite.AccessLevel,
			OrganizationId: invite.OrganizationId,
			UserId:         insertedUser.UniqueId,
			InviteId:       &invite.UniqueId,
		}).QueryContext(context.Request().Context(), context.App.Db, &insertedOrgMember)
	}

	contextUser := interfaces.ContextUser{
		Username: insertedUser.Username,
		Email:    insertedUser.Email,
		Role:     api_types.UserRoleEnum(insertedOrgMember.AccessLevel),
		UniqueId: insertedUser.UniqueId.String(),
		Name:     insertedUser.Name,
	}

	if insertedOrgMember.UniqueId.String() != "" {
		contextUser.OrganizationId = invite.OrganizationId.String()
	}

	claims := &interfaces.JwtPayload{
		ContextUser: contextUser,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Hour * 24 * 60).Unix(), // 60-day expiration
			Issuer:    "wapikit",
		},
	}

	token, err := jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString([]byte(context.App.Koa.String("app.jwt_secret")))

	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Error generating token")
	}

	return context.JSON(http.StatusOK, api_types.VerifyOtpResponseBodySchema{
		Token: token,
	})
}

func regenerateApiKey(context interfaces.ContextWithSession) error {
	orgUuid, err := uuid.Parse(context.Session.User.UniqueId)
	if err != nil {
		return echo.NewHTTPError(http.StatusUnauthorized, "Unauthorized access")
	}

	userUuid, err := uuid.Parse(context.Session.User.UniqueId)

	if err != nil {
		return echo.NewHTTPError(http.StatusUnauthorized, "Unauthorized access")
	}

	orgMemberQuery := SELECT(table.OrganizationMember.AllColumns).
		FROM(table.OrganizationMember).
		WHERE(table.OrganizationMember.UserId.EQ(UUID(userUuid)).
			AND(table.OrganizationMember.OrganizationId.EQ(UUID(orgUuid))),
		).AsTable("Member")

	apiKeyQuery := SELECT(table.ApiKey.AllColumns).
		FROM(table.ApiKey).
		WHERE(table.ApiKey.OrganizationId.EQ(UUID(orgUuid))).AsTable("ApiKeys")

	currentUserApiKey := SELECT(
		apiKeyQuery.AllColumns(),
		orgMemberQuery.AllColumns(),
	).FROM(apiKeyQuery.INNER_JOIN(
		table.OrganizationMember, table.OrganizationMember.UniqueId.EQ(table.ApiKey.MemberId),
	)).LIMIT(1)

	var dest model.ApiKey

	err = currentUserApiKey.QueryContext(context.Request().Context(), context.App.Db, &dest)

	if err != nil {
		// ! it can not be possible that no rows found for an API KEY
		if err.Error() == "qrm: no rows in result set" {
			return echo.NewHTTPError(http.StatusNotFound, "Exisitng API key not found, report this bug at contact@wapikit.com")
		} else {
			context.App.Logger.Error("database query error", err.Error())
			return echo.NewHTTPError(http.StatusInternalServerError, "Something went wrong while processing your request.")
		}
	}

	accessLevel := context.Session.User.Role

	claims := &interfaces.JwtPayload{
		ContextUser: interfaces.ContextUser{
			Username:       context.Session.User.Username,
			Email:          context.Session.User.Email,
			Role:           accessLevel,
			UniqueId:       context.Session.User.UniqueId,
			OrganizationId: orgUuid.String(),
			Name:           context.Session.User.Name,
		},
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Hour * 24 * 365 * 2).Unix(), // 60-day expiration
			Issuer:    "wapikit",
		},
	}

	//Create the token
	token, err := jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString([]byte(context.App.Koa.String("app.jwt_secret")))
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Error generating token")
	}

	var updatedApiKey model.ApiKey

	err = table.ApiKey.UPDATE(table.ApiKey.Key).MODEL(model.ApiKey{
		Key: token,
	}).RETURNING(table.ApiKey.AllColumns).
		QueryContext(context.Request().Context(), context.App.Db, &updatedApiKey)

	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Something went wrong while processing your request.")
	}

	response := api_types.RegenerateApiKeyResponseSchema{
		ApiKey: &api_types.ApiKeySchema{
			CreatedAt: updatedApiKey.CreatedAt,
			Key:       updatedApiKey.Key,
			UniqueId:  updatedApiKey.UniqueId.String(),
		},
	}

	return context.JSON(http.StatusOK, response)
}

func getApiKeys(context interfaces.ContextWithSession) error {
	user := context.Session.User
	var apiKeys []model.ApiKey
	stmt := SELECT(table.ApiKey.AllColumns).
		FROM(table.ApiKey.
			RIGHT_JOIN(
				table.OrganizationMember,
				table.OrganizationMember.UserId.EQ(String(user.UniqueId)).
					AND(table.OrganizationMember.UniqueId.EQ(table.ApiKey.MemberId)).
					AND(table.Organization.UniqueId.EQ(table.ApiKey.OrganizationId))))

	stmt.Query(database.GetDbInstance(), &apiKeys)
	apiKeysToReturn := make([]api_types.ApiKeySchema, 0)
	for _, apiKey := range apiKeys {
		uniqueId := apiKey.UniqueId.String()
		apiKeysToReturn = append(apiKeysToReturn, api_types.ApiKeySchema{
			CreatedAt: apiKey.CreatedAt,
			Key:       apiKey.Key,
			UniqueId:  uniqueId,
		})
	}
	return context.JSON(http.StatusOK, api_types.GetApiKeysResponseSchema{
		ApiKeys: apiKeysToReturn,
	})
}

func switchOrganization(context interfaces.ContextWithSession) error {
	// organization id
	payload := new(api_types.SwitchOrganizationJSONRequestBody)
	if err := context.Bind(payload); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	// get the current organization id
	currentAuthedOrganizationId := context.Session.User.OrganizationId

	if currentAuthedOrganizationId == *payload.OrganizationId {
		// bad request
		return echo.NewHTTPError(http.StatusBadRequest, "Already in the same organization")
	}

	newOrgQuery := SELECT(table.Organization.AllColumns, table.OrganizationMember.AllColumns).FROM(table.Organization.LEFT_JOIN(table.OrganizationMember, table.OrganizationMember.OrganizationId.EQ(String(*payload.OrganizationId)).AND(table.OrganizationMember.UniqueId.EQ(String(context.Session.User.UniqueId))))).WHERE(table.Organization.UniqueId.EQ(String(*payload.OrganizationId)))

	var newOrgDetails struct {
		model.Organization `json:"-,inline"`
		MemberDetails      model.OrganizationMember `json:"member_details"`
	}

	newOrgQuery.Query(database.GetDbInstance(), &newOrgDetails)

	if newOrgDetails.UniqueId.String() == "" {
		return echo.NewHTTPError(http.StatusNotFound, "Organization not found")
	}

	if newOrgDetails.MemberDetails.UniqueId.String() == "" {
		return echo.NewHTTPError(http.StatusUnauthorized, "User not a member of the organization")
	}

	// create the token

	claims := &interfaces.JwtPayload{
		ContextUser: interfaces.ContextUser{
			Username:       context.Session.User.Username,
			Email:          context.Session.User.Email,
			Role:           api_types.UserRoleEnum(newOrgDetails.MemberDetails.AccessLevel),
			UniqueId:       context.Session.User.UniqueId,
			OrganizationId: *payload.OrganizationId,
			Name:           context.Session.User.Name,
		},
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Hour * 24 * 60).Unix(), // 60-day expiration
			Issuer:    "wapikit",
		},
	}

	//Create the token

	token, err := jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString([]byte(context.App.Koa.String("app.jwt_secret")))
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Error generating token")
	}

	return context.JSON(http.StatusOK, api_types.SwitchOrganizationResponseSchema{
		Token: token,
	})
}
