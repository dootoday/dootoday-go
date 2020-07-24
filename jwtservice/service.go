package service

import (
	"strconv"
	"time"

	"apidootoday/config"

	"github.com/brianvoe/sjwt"
)

// TokenType :
type TokenType string

const (
	// AccessTokenType :
	AccessTokenType = TokenType("access")

	// RefreshTokenType :
	RefreshTokenType = TokenType("refresh")
)

// TokenService :
type TokenService struct{}

// NewTokenService : instantiate a new service
func NewTokenService() *TokenService {
	return &TokenService{}
}

// GetAccessToken :
func (j *TokenService) GetAccessToken(userID uint) string {
	claims := sjwt.New()
	claims.Set("user_id", userID)
	claims.SetIssuedAt(time.Now().UTC())
	claims.SetExpiresAt(time.Now().UTC().Add(time.Hour * 24 * 7)) // 7 days
	secretKey := []byte(config.AccessTokenSecret)
	jwt := claims.Generate(secretKey)
	return jwt
}

// GetRefreshToken :
func (j *TokenService) GetRefreshToken(userID uint) string {
	claims := sjwt.New()
	claims.Set("user_id", userID)
	claims.SetIssuedAt(time.Now().UTC())
	claims.SetExpiresAt(time.Now().UTC().Add(time.Hour * 24 * 30 * 12)) // 1 year
	secretKey := []byte(config.RefreshTokenSecret)
	jwt := claims.Generate(secretKey)
	return jwt
}

// IsTokenValid : check for validity
func (j *TokenService) IsTokenValid(token string, tokenType TokenType) (bool, error) {
	secretKey := []byte(config.AccessTokenSecret)
	if tokenType == RefreshTokenType {
		secretKey = []byte(config.RefreshTokenSecret)
	}
	verified := sjwt.Verify(token, secretKey)
	if !verified {
		return false, nil
	}
	claims, err := sjwt.Parse(token)
	if err != nil {
		return false, err
	}
	expiry, err := claims.GetExpiresAt()
	if err != nil {
		return false, err
	}
	if expiry < time.Now().UTC().Unix() {
		return false, nil
	}
	return true, nil
}

// GetUserIDFromToken : get the user id from the token
func (j *TokenService) GetUserIDFromToken(token string) (uint, error) {
	claims, err := sjwt.Parse(token)
	if err != nil {
		return uint(0), err
	}
	userIDstr, err := claims.GetStr("user_id")
	userID, err := strconv.ParseUint(userIDstr, 10, 32)
	if err != nil {
		return uint(0), err
	}
	return uint(userID), nil
}

// GetInfoFromToken : get info from token
func (j *TokenService) GetInfoFromToken(token string, key string) (string, error) {
	claims, err := sjwt.Parse(token)
	if err != nil {
		return "", err
	}
	return claims.GetStr(key)
}
