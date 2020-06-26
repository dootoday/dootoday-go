package service

import (
	js "apidootoday/jwtservice"

	"google.golang.org/api/oauth2/v2"
)

// IGoogleService : interface for google service
type IGoogleService interface {
	VerifyIDToken(idToken string) (*oauth2.Tokeninfo, error)
}

// IJWTService : interface for JWT service
type IJWTService interface {
	GetAccessToken(userID uint) string
	GetRefreshToken(userID uint) string
	IsTokenValid(token string, tokenType js.TokenType) (bool, error)
	GetUserIDFromToken(token string) (uint, error)
	GetInfoFromToken(token string, key string) (string, error)
}
