package service

import (
	"strconv"
	"time"

	"apidootoday/config"

	"github.com/brianvoe/sjwt"
)

// GetAccessToken :
func GetAccessToken(userID uint) string {
	claims := sjwt.New()
	claims.Set("user_id", userID)
	claims.SetIssuedAt(time.Now().UTC())
	claims.SetExpiresAt(time.Now().UTC().Add(time.Hour * 720))
	secretKey := []byte(config.AccessTokenSecret)
	jwt := claims.Generate(secretKey)
	return jwt
}

// GetRefreshToken :
func GetRefreshToken(userID uint) string {
	claims := sjwt.New()
	claims.Set("user_id", userID)
	claims.SetIssuedAt(time.Now().UTC())
	claims.SetExpiresAt(time.Now().UTC().Add(time.Hour * 1))
	secretKey := []byte(config.RefreshTokenSecret)
	jwt := claims.Generate(secretKey)
	return jwt
}

// TokenType :
type TokenType string

const (
	// AccessTokenType :
	AccessTokenType = TokenType("access")

	// RefreshTokenType :
	RefreshTokenType = TokenType("refresh")
)

// TokenService :
type TokenService struct {
	token     string
	tokenType TokenType
}

// NewTokenService : instantiate a new service
func NewTokenService(token string, ttype TokenType) *TokenService {
	return &TokenService{
		token:     token,
		tokenType: ttype,
	}
}

// IsTokenValid : check for validity
func (j *TokenService) IsTokenValid() (bool, error) {
	secretKey := []byte(config.AccessTokenSecret)
	if j.tokenType == RefreshTokenType {
		secretKey = []byte(config.RefreshTokenSecret)
	}
	verified := sjwt.Verify(j.token, secretKey)
	if !verified {
		return false, nil
	}
	claims, err := sjwt.Parse(j.token)
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

// GetUserID : get the user id from the token
func (j *TokenService) GetUserID() (uint, error) {
	claims, err := sjwt.Parse(j.token)
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
