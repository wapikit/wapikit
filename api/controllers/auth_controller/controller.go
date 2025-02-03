package auth_controller

import (
	"fmt"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/wapikit/wapikit/api/api_types"
	controller "github.com/wapikit/wapikit/api/controllers"
	"github.com/wapikit/wapikit/interfaces"
	"github.com/wapikit/wapikit/utils"
	"golang.org/x/crypto/bcrypt"

	"github.com/go-jet/jet/qrm"
	. "github.com/go-jet/jet/v2/postgres"
	"github.com/wapikit/wapikit/.db-generated/model"
	table "github.com/wapikit/wapikit/.db-generated/table"
)

type AuthController struct {
	controller.BaseController `json:"-,inline"`
}

func NewAuthController() *AuthController {
	return &AuthController{
		BaseController: controller.BaseController{
			Name:        "Auth Controller",
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
					Handler:                 interfaces.HandlerWithSession(getApiKey),
					IsAuthorizationRequired: true,
					MetaData: interfaces.RouteMetaData{
						PermissionRoleLevel: api_types.Member,
						RateLimitConfig: interfaces.RateLimitConfig{
							MaxRequests:    60,
							WindowTimeInMs: 1000 * 60 * 60, // 1 hour
						},
						RequiredPermission: []api_types.RolePermissionEnum{
							api_types.GetApiKey,
						},
					},
				},
				{
					Path:                    "/api/auth/api-keys/regenerate",
					Method:                  http.MethodGet,
					Handler:                 interfaces.HandlerWithSession(regenerateApiKey),
					IsAuthorizationRequired: true,
					MetaData: interfaces.RouteMetaData{
						PermissionRoleLevel: api_types.Member,
						RateLimitConfig: interfaces.RateLimitConfig{
							MaxRequests:    60,
							WindowTimeInMs: 1000 * 60 * 60, // 1 hour
						},
						RequiredPermission: []api_types.RolePermissionEnum{
							api_types.RegenerateApiKey,
						},
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
						PermissionRoleLevel: api_types.Member,
						RateLimitConfig: interfaces.RateLimitConfig{
							MaxRequests:    60,
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
		return context.JSON(http.StatusBadRequest, err.Error())
	}

	if payload.InviteSlug == nil {
		return context.JSON(http.StatusBadRequest, "Organization invite slug is required")
	}

	user := context.Session.User
	userUuid, err := uuid.Parse(user.UniqueId)

	if err != nil {
		return context.JSON(http.StatusUnauthorized, "Unauthorized access")
	}

	var invite model.OrganizationMemberInvite

	invitationQuery := SELECT(table.OrganizationMemberInvite.AllColumns).
		FROM(table.OrganizationMemberInvite).
		WHERE(
			table.OrganizationMemberInvite.Email.EQ(String(user.Email)).
				AND(table.OrganizationMemberInvite.Slug.EQ(String(*payload.InviteSlug)))).
		LIMIT(1)

	err = invitationQuery.QueryContext(context.Request().Context(), context.App.Db, &invite)

	if err != nil {
		if err.Error() == qrm.ErrNoRows.Error() {
			return context.JSON(http.StatusNotFound, "Organization invite not found")
		} else {
			context.App.Logger.Error("database query error", err.Error(), nil)
			return context.JSON(http.StatusInternalServerError, "Something went wrong while processing your request.")
		}
	}

	if invite.Status != model.OrganizationInviteStatusEnum_Pending {
		return context.JSON(http.StatusForbidden, "Invite already accepted")
	}

	var existingMember model.OrganizationMember

	existingMemberQuery := SELECT(table.OrganizationMember.AllColumns).
		FROM(table.OrganizationMember).
		WHERE(
			table.OrganizationMember.UserId.EQ(UUID(userUuid)).
				AND(table.OrganizationMember.OrganizationId.EQ(UUID(invite.OrganizationId)))).
		LIMIT(1)

	err = existingMemberQuery.QueryContext(context.Request().Context(), context.App.Db, &existingMember)

	if err != nil {
		if err.Error() == qrm.ErrNoRows.Error() {
			// do nothing
		}
	} else {
		return context.JSON(http.StatusForbidden, "User already a member of the organization")
	}

	var insertedOrgMember model.OrganizationMember

	err = table.OrganizationMember.INSERT(table.OrganizationMember.MutableColumns).MODEL(model.OrganizationMember{
		AccessLevel:    invite.AccessLevel,
		OrganizationId: invite.OrganizationId,
		UserId:         userUuid,
		InviteId:       &invite.UniqueId,
	}).QueryContext(context.Request().Context(), context.App.Db, &insertedOrgMember)

	if err != nil {
		context.App.Logger.Error("database query error", err.Error())
		return context.JSON(http.StatusInternalServerError, "Something went wrong while processing your request.")
	}

	aiChatsToCreate := model.AiChat{
		UniqueId:             uuid.New(),
		CreatedAt:            time.Now(),
		UpdatedAt:            time.Now(),
		Status:               model.AiChatStatusEnum_Active,
		OrganizationId:       insertedOrgMember.OrganizationId,
		OrganizationMemberId: insertedOrgMember.UniqueId,
		Title:                "Default Chat",
		Visibility:           model.AiChatVisibilityEnum_Public,
	}

	_, err = table.AiChat.INSERT(table.AiChat.AllColumns).
		MODELS(aiChatsToCreate).
		ExecContext(context.Request().Context(), context.App.Db)

	if err != nil {
		return context.JSON(http.StatusInternalServerError, err.Error())
	}

	// create the token
	claims := &interfaces.JwtPayload{
		ContextUser: interfaces.ContextUser{
			Username:       context.Session.User.Username,
			Email:          context.Session.User.Email,
			Role:           api_types.UserPermissionLevelEnum(invite.AccessLevel),
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
		return context.JSON(http.StatusInternalServerError, "Error generating token")
	}
	response := api_types.JoinOrganizationResponseBodySchema{
		Token: token,
	}
	return context.JSON(http.StatusOK, response)
}

func handleSignIn(context interfaces.ContextWithoutSession) error {
	payload := new(api_types.LoginRequestBodySchema)

	if err := context.Bind(payload); err != nil {
		return context.JSON(http.StatusBadRequest, err.Error())
	}

	if payload.Username == "" || payload.Password == "" {
		return context.JSON(echo.ErrBadRequest.Code, "Username / password is required")
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

	stmt.QueryContext(context.Request().Context(), context.App.Db, &user)

	// if no user found then return 404
	if user.User.UniqueId.String() == "" || user.User.Password == nil {
		return context.JSON(http.StatusNotFound, "Invalid email / password")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(*user.User.Password), []byte(payload.Password)); err != nil {
		return context.JSON(http.StatusNotFound, "Invalid email / password")
	}

	var isOnboardingCompleted bool
	var organizationIdToLoginWith string
	var roleToLoginWith api_types.UserPermissionLevelEnum
	var claims *interfaces.JwtPayload

	// if no organization found, then simply return the user with a flag saying isOnboardingCompleted
	if len(user.Organizations) == 0 {
		isOnboardingCompleted = false
		// onboarding to be completed by the user
		// return the user with onboarding flag

		claims = &interfaces.JwtPayload{
			ContextUser: interfaces.ContextUser{
				Username:       user.User.Username,
				Email:          user.User.Email,
				UniqueId:       user.User.UniqueId.String(),
				Name:           user.User.Name,
				OrganizationId: "",
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
			if org.MemberDetails.AccessLevel == model.UserPermissionLevelEnum_Owner {
				organizationIdToLoginWith = org.Organization.UniqueId.String()
				roleToLoginWith = api_types.Owner
				break
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
				Role:           api_types.UserPermissionLevelEnum(roleToLoginWith),
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
		return context.JSON(http.StatusInternalServerError, "Error generating token")
	}

	// Set the cookie in the response
	cookie := new(http.Cookie)
	cookie.Name = "__auth_token"
	cookie.Value = token
	cookie.Path = "/"
	cookie.HttpOnly = true
	cookie.Secure = true                                 // Set this to true in production for HTTPS
	cookie.Domain = ".wapikit.com"                       // Ensure the domain matches your app
	cookie.Expires = time.Now().Add(time.Hour * 24 * 60) // 60-day expiration
	context.SetCookie(cookie)

	return context.JSON(http.StatusOK, api_types.LoginResponseBodySchema{
		IsOnboardingCompleted: isOnboardingCompleted,
		Token:                 token,
	})
}

func handleLoginWithOAuth(context interfaces.ContextWithoutSession) error {
	return context.JSON(http.StatusInternalServerError, "Not implemented")
}

// this handler would validate the email and send an otp to it
func handleUserRegistration(context interfaces.ContextWithoutSession) error {
	redis := context.App.Redis

	payload := new(api_types.RegisterRequestBodySchema)

	if err := context.Bind(payload); err != nil {
		return context.JSON(http.StatusBadRequest, err.Error())
	}

	if payload.Username == "" || payload.Email == "" || payload.Password == "" {
		return context.JSON(http.StatusBadRequest, "Username, Email and Password are required")
	}

	otp := utils.GenerateOtp(context.App.Constants.IsProduction)
	cacheKey := redis.ComputeCacheKey("otp", payload.Email, "registration")
	err := redis.CacheData(cacheKey, otp, time.Minute*5)

	if err != nil {
		context.App.Logger.Error("error caching otp", err.Error(), nil)
		return context.JSON(http.StatusInternalServerError, "Something went wrong while processing your request.")
	}

	err = context.App.NotificationService.SendEmail(payload.Email, "Wapikit Registration OTP", fmt.Sprintf("Your OTP is %s", otp), context.App.Constants.IsProduction)

	if err != nil {
		context.App.Logger.Error("error sending email", err.Error(), nil)
		return context.JSON(http.StatusInternalServerError, "Something went wrong while processing your request.")
	}

	// return the response
	return context.JSON(http.StatusOK, api_types.RegisterRequestResponseBodySchema{
		IsOtpSent: true,
	})
}

func verifyEmailAndCreateAccount(context interfaces.ContextWithoutSession) error {
	redis := context.App.Redis
	payload := new(api_types.VerifyOtpJSONRequestBody)
	if err := context.Bind(payload); err != nil {
		return context.JSON(http.StatusBadRequest, err.Error())
	}

	if payload.Username == "" || payload.Email == "" || payload.Password == "" || payload.Otp == "" {
		return context.JSON(http.StatusBadRequest, "Invalid details!!")
	}

	cacheKey := redis.ComputeCacheKey("otp", payload.Email, "registration")
	cachedOtp, err := redis.GetCachedData(cacheKey)
	if err != nil {
		context.App.Logger.Error("Error getting cached otp", err.Error())
		return context.JSON(http.StatusInternalServerError, "Something went wrong while processing your request.")
	}

	if cachedOtp != payload.Otp {
		return context.JSON(http.StatusBadRequest, "Invalid OTP")
	}

	// check if the user already exists
	var user model.User

	userQuery := SELECT(table.User.AllColumns).
		FROM(table.User).
		WHERE(table.User.Email.EQ(String(payload.Email))).
		LIMIT(1)

	err = userQuery.QueryContext(context.Request().Context(), context.App.Db, &user)

	if err != nil {
		if err.Error() == qrm.ErrNoRows.Error() {
			// do nothing
		} else {
			context.App.Logger.Error("database query error", err.Error())
			return context.JSON(http.StatusInternalServerError, "Something went wrong while processing your request.")
		}
	}

	if user.UniqueId.String() != uuid.Nil.String() {
		return context.JSON(http.StatusBadRequest, "User already exists")
	}

	var invite model.OrganizationMemberInvite

	if payload.OrganizationInviteSlug != nil {
		invitationQuery := SELECT(table.OrganizationMemberInvite.AllColumns).
			FROM(table.OrganizationMemberInvite).
			WHERE(
				table.OrganizationMemberInvite.Email.EQ(String(payload.Email)).
					AND(table.OrganizationMemberInvite.Slug.EQ(String(*payload.OrganizationInviteSlug)))).
			LIMIT(1)

		err = invitationQuery.QueryContext(context.Request().Context(), context.App.Db, &invite)

		if err != nil {
			if err.Error() == qrm.ErrNoRows.Error() {
				// do  nothing just move on, we cant let the user not register if they do not have a valid invite
			} else {
				context.App.Logger.Error("database query error", err.Error())
				return context.JSON(http.StatusInternalServerError, "Something went wrong while processing your request.")
			}
		}
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(payload.Password), bcrypt.DefaultCost)
	if err != nil {
		return context.JSON(http.StatusInternalServerError, "Error hashing password")
	}

	passwordString := string(hashedPassword)

	var insertedUser model.User
	var insertedOrgMember model.OrganizationMember

	err = table.User.INSERT(table.User.MutableColumns).MODEL(model.User{
		Username:  payload.Username,
		Email:     payload.Email,
		Password:  &passwordString,
		Name:      payload.Name,
		Status:    model.UserAccountStatusEnum_Active,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}).RETURNING(table.User.AllColumns).QueryContext(context.Request().Context(), context.App.Db, &insertedUser)

	if err != nil {
		context.App.Logger.Error("database query error", err.Error())
		return context.JSON(http.StatusInternalServerError, "Something went wrong while processing your request.")
	}

	if invite.UniqueId.String() != "" {
		err = table.OrganizationMember.INSERT(
			table.OrganizationMember.MutableColumns,
		).MODEL(model.OrganizationMember{
			AccessLevel:    invite.AccessLevel,
			OrganizationId: invite.OrganizationId,
			UserId:         insertedUser.UniqueId,
			InviteId:       &invite.UniqueId,
			CreatedAt:      time.Now(),
			UpdatedAt:      time.Now(),
		}).QueryContext(context.Request().Context(), context.App.Db, &insertedOrgMember)
	}

	var role api_types.UserPermissionLevelEnum

	if insertedOrgMember.UniqueId.String() != uuid.Nil.String() {
		role = api_types.UserPermissionLevelEnum(insertedOrgMember.AccessLevel)
		aiChatsToCreate := model.AiChat{
			UniqueId:             uuid.New(),
			CreatedAt:            time.Now(),
			UpdatedAt:            time.Now(),
			Status:               model.AiChatStatusEnum_Active,
			OrganizationId:       insertedOrgMember.OrganizationId,
			OrganizationMemberId: insertedOrgMember.UniqueId,
			Title:                "Default Chat",
			Visibility:           model.AiChatVisibilityEnum_Public,
		}

		_, err = table.AiChat.INSERT(table.AiChat.AllColumns).
			MODELS(aiChatsToCreate).
			ExecContext(context.Request().Context(), context.App.Db)

		if err != nil {
			return context.JSON(http.StatusInternalServerError, err.Error())
		}
	}

	contextUser := interfaces.ContextUser{
		Username:       insertedUser.Username,
		Email:          insertedUser.Email,
		Role:           role,
		UniqueId:       insertedUser.UniqueId.String(),
		Name:           insertedUser.Name,
		OrganizationId: "",
	}

	if insertedOrgMember.UniqueId.String() != uuid.Nil.String() {
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
		return context.JSON(http.StatusInternalServerError, "Error generating token")
	}

	// Set the cookie in the response
	cookie := new(http.Cookie)
	cookie.Name = "__auth_token"
	cookie.Value = token
	cookie.Path = "/"
	cookie.HttpOnly = true
	cookie.Secure = true                                 // Set this to true in production for HTTPS
	cookie.Domain = ".wapikit.com"                       // Ensure the domain matches your app
	cookie.Expires = time.Now().Add(time.Hour * 24 * 60) // 60-day expiration
	context.SetCookie(cookie)

	return context.JSON(http.StatusOK, api_types.VerifyOtpResponseBodySchema{
		Token: token,
	})
}

func regenerateApiKey(context interfaces.ContextWithSession) error {

	user := context.Session.User
	var apiKey model.ApiKey

	userUuid, err := uuid.Parse(user.UniqueId)

	if err != nil {
		return context.JSON(http.StatusUnauthorized, "Unauthorized access")
	}

	orgUuid, err := uuid.Parse(user.OrganizationId)

	if err != nil {
		return context.JSON(http.StatusUnauthorized, "Unauthorized access")
	}

	var orgMember model.OrganizationMember

	organizationMemberQuery := SELECT(table.OrganizationMember.AllColumns).
		FROM(table.OrganizationMember).
		WHERE(table.OrganizationMember.UserId.EQ(UUID(userUuid)).AND(
			table.OrganizationMember.OrganizationId.EQ(UUID(orgUuid))),
		).
		LIMIT(1)

	err = organizationMemberQuery.QueryContext(context.Request().Context(), context.App.Db, &orgMember)

	if err != nil {
		if err.Error() == qrm.ErrNoRows.Error() {
			return context.JSON(http.StatusUnauthorized, "Unauthorized access")
		} else {
			return context.JSON(http.StatusInternalServerError, "Something went wrong while processing your request.")
		}
	}

	stmt := SELECT(table.ApiKey.AllColumns).
		FROM(table.ApiKey).
		WHERE(table.ApiKey.MemberId.EQ(UUID(orgMember.UniqueId))).
		LIMIT(1)

	stmt.Query(context.App.Db, &apiKey)

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
		return context.JSON(http.StatusInternalServerError, "Error generating token")
	}

	var updatedApiKey model.ApiKey

	err = table.ApiKey.UPDATE(table.ApiKey.Key).MODEL(model.ApiKey{
		Key: token,
	}).WHERE(table.ApiKey.UniqueId.EQ(UUID(apiKey.UniqueId))).RETURNING(table.ApiKey.AllColumns).
		QueryContext(context.Request().Context(), context.App.Db, &updatedApiKey)

	if err != nil {
		return context.JSON(http.StatusInternalServerError, "Something went wrong while processing your request.")
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

func getApiKey(context interfaces.ContextWithSession) error {
	user := context.Session.User
	var apiKey model.ApiKey

	userUuid, err := uuid.Parse(user.UniqueId)

	if err != nil {
		return context.JSON(http.StatusUnauthorized, "Unauthorized access")
	}

	orgUuid, err := uuid.Parse(user.OrganizationId)

	if err != nil {
		return context.JSON(http.StatusUnauthorized, "Unauthorized access")
	}

	var orgMember model.OrganizationMember

	organizationMemberQuery := SELECT(table.OrganizationMember.AllColumns).
		FROM(table.OrganizationMember).
		WHERE(table.OrganizationMember.UserId.EQ(UUID(userUuid)).AND(
			table.OrganizationMember.OrganizationId.EQ(UUID(orgUuid))),
		).
		LIMIT(1)

	err = organizationMemberQuery.QueryContext(context.Request().Context(), context.App.Db, &orgMember)

	if err != nil {
		if err.Error() == qrm.ErrNoRows.Error() {
			return context.JSON(http.StatusUnauthorized, "Unauthorized access")
		} else {
			return context.JSON(http.StatusInternalServerError, "Something went wrong while processing your request.")
		}
	}

	stmt := SELECT(table.ApiKey.AllColumns).
		FROM(table.ApiKey).
		WHERE(table.ApiKey.MemberId.EQ(UUID(orgMember.UniqueId))).
		LIMIT(1)

	stmt.Query(context.App.Db, &apiKey)

	uniqueId := apiKey.UniqueId.String()

	apiKeysToReturn := api_types.ApiKeySchema{
		CreatedAt: apiKey.CreatedAt,
		Key:       apiKey.Key,
		UniqueId:  uniqueId,
	}

	return context.JSON(http.StatusOK, api_types.GetApiKeysResponseSchema{
		ApiKey: apiKeysToReturn,
	})
}

func switchOrganization(context interfaces.ContextWithSession) error {
	// organization id
	payload := new(api_types.SwitchOrganizationJSONRequestBody)
	if err := context.Bind(payload); err != nil {
		return context.JSON(http.StatusBadRequest, err.Error())
	}

	// get the current organization id
	currentAuthedOrganizationId := context.Session.User.OrganizationId

	if currentAuthedOrganizationId == *payload.OrganizationId {
		// bad request
		return context.JSON(http.StatusBadRequest, "Already in the same organization")
	}

	newOrgUuid, err := uuid.Parse(*payload.OrganizationId)

	if err != nil {
		return context.JSON(http.StatusBadRequest, "Invalid organization id")
	}

	userUuid, err := uuid.Parse(context.Session.User.UniqueId)

	if err != nil {
		return context.JSON(http.StatusUnauthorized, "Unauthorized access")
	}

	newOrgQuery := SELECT(
		table.Organization.AllColumns,
		table.OrganizationMember.AllColumns,
	).
		FROM(table.Organization.
			LEFT_JOIN(table.OrganizationMember,
				table.OrganizationMember.OrganizationId.EQ(UUID(newOrgUuid)).
					AND(table.OrganizationMember.UserId.EQ(UUID(userUuid))))).
		WHERE(table.Organization.UniqueId.EQ(UUID(newOrgUuid)))

	var newOrgDetails struct {
		model.Organization
		MemberDetails struct {
			model.OrganizationMember
		}
	}

	newOrgQuery.Query(context.App.Db, &newOrgDetails)

	if newOrgDetails.UniqueId.String() == "" {
		return context.JSON(http.StatusNotFound, "Organization not found")
	}

	if newOrgDetails.MemberDetails.UniqueId.String() == "" {
		return context.JSON(http.StatusUnauthorized, "User not a member of the organization")
	}

	fmt.Println("newOrgDetails", newOrgDetails)

	// create the token
	claims := &interfaces.JwtPayload{
		ContextUser: interfaces.ContextUser{
			Username:       context.Session.User.Username,
			Email:          context.Session.User.Email,
			Role:           api_types.UserPermissionLevelEnum(newOrgDetails.MemberDetails.AccessLevel),
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
		return context.JSON(http.StatusInternalServerError, "Error generating token")
	}

	return context.JSON(http.StatusOK, api_types.SwitchOrganizationResponseSchema{
		Token: token,
	})
}
