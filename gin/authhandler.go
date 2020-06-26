package service

import (
	gauthservice "apidootoday/googleauth"
	jwtservice "apidootoday/jwtservice"
	userservice "apidootoday/user"
	"net/http"

	"github.com/gin-gonic/gin"
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
	if err != nil || request.IDToken == "" {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}

	userID, err := ah.UserService.Login(request.IDToken)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}

	resp := ResponseBody{
		AccessToken:  ah.TokenService.GetAccessToken(userID),
		RefreshToken: ah.TokenService.GetRefreshToken(userID),
	}

	c.JSON(http.StatusOK, resp)
}
