package handlers

import (
	"github.com/golang-jwt/jwt"
	"github.com/sarthakjdev/wapikit/internal"
)

type AuthHandlerBodySchemaType struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type CreateNewUserHandlerBodySchemaType struct {
	Username    string `json:"username"`
	Email       string `json:"email"`
	Password    string `json:"password"`
	PhoneNumber string `json:"phone_number"`
}

type UpdateBusinessAccountDetailsHandlerBodySchemaType struct {
	BusinessAccountId string `json:"business_account_id"`
}

type GetUserByIdHandlerDataSchemaType struct {
	Id string `json:"id"`
}

type JwtPayload struct {
	internal.ContextUser `json:",inline"`
	jwt.StandardClaims   `json:",inline"`
}
