package interfaces

import (
	"github.com/golang-jwt/jwt"
	"github.com/wapikit/wapikit/api/api_types"
)

type ContextUser struct {
	Name           string                            `json:"name"`
	UniqueId       string                            `json:"unique_id"`
	Username       string                            `json:"username"`
	Email          string                            `json:"email"`
	Role           api_types.UserPermissionLevelEnum `json:"role"`
	OrganizationId string                            `json:"organization_id"`
}

type ContextSession struct {
	Token string      `json:"token"`
	User  ContextUser `json:"user"`
}

type JwtPayload struct {
	ContextUser        `json:",inline"`
	jwt.StandardClaims `json:",inline"`
}

type Constants struct {
	IsDevelopment         bool
	IsProduction          bool
	RedisEventChannelName string `koanf:"redis_event_channel_name"`
	IsDebugModeEnabled    bool
	IsCommunityEdition    bool
	IsCloudEdition        bool
	IsSingleBinaryMode    bool
}
