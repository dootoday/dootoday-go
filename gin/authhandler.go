package service

import (
	gauthservice "apidootoday/googleauth"
	jwtservice "apidootoday/jwtservice"
	userservice "apidootoday/user"
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/golang/glog"
)

// AuthHandler :
type AuthHandler struct {
	UserService  *userservice.UserService
	TokenService *jwtservice.TokenService
	GAuthService *gauthservice.GoogleAuthService
}

// NewAuthHandler :
func NewAuthHandler(
	userService *userservice.UserService,
	tokenService *jwtservice.TokenService,
	gauthService *gauthservice.GoogleAuthService,
) *AuthHandler {
	return &AuthHandler{
		UserService:  userService,
		TokenService: tokenService,
		GAuthService: gauthService,
	}
}

// Login : login handler
func (ah *AuthHandler) Login(c *gin.Context) {
	type RequestBody struct {
		IDToken string `json:"id_token"`
	}
	type ResponseBody struct {
		AccessToken  string `json:"access_token"`
		RefreshToken string `json:"refresh_token"`
	}
	var request RequestBody
	err := c.BindJSON(&request)
	if err != nil {
		glog.Error(err)
		c.JSON(http.StatusBadRequest, err)
		return
	}
	if request.IDToken == "" {
		c.JSON(http.StatusBadRequest, errors.New("empty token id"))
		return
	}

	userID, err := ah.UserService.Login(request.IDToken)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err)
		return
	}

	resp := ResponseBody{
		AccessToken:  ah.TokenService.GetAccessToken(userID),
		RefreshToken: ah.TokenService.GetRefreshToken(userID),
	}

	c.JSON(http.StatusOK, resp)
}

// Refresh : refresh token
func (ah *AuthHandler) Refresh(c *gin.Context) {
	type RequestBody struct {
		RefreshToken string `json:"refresh_token"`
	}
	type ResponseBody struct {
		AccessToken string `json:"access_token"`
	}
	var request RequestBody
	err := c.BindJSON(&request)
	if err != nil || request.RefreshToken == "" {
		c.JSON(http.StatusBadRequest, err)
		return
	}
	valid, err := ah.TokenService.IsTokenValid(
		request.RefreshToken, jwtservice.RefreshTokenType,
	)
	if err != nil {
		glog.Error(err)
		c.JSON(http.StatusInternalServerError, err)
		return
	}
	if !valid {
		c.JSON(http.StatusForbidden, errors.New("Refresh token is not valid"))
		return
	}
	userID, err := ah.TokenService.GetUserIDFromToken(request.RefreshToken)
	if err != nil {
		glog.Error(err)
		c.JSON(http.StatusInternalServerError, err)
		return
	}
	resp := ResponseBody{
		AccessToken: ah.TokenService.GetAccessToken(userID),
	}
	c.JSON(http.StatusOK, resp)
}
