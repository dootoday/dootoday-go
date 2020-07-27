package service

import (
	gauthservice "apidootoday/googleauth"
	jwtservice "apidootoday/jwtservice"
	subscriptionservice "apidootoday/subscription"
	taskservice "apidootoday/taskservice"
	userservice "apidootoday/user"
	"errors"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang/glog"
)

// AuthHandler :
type AuthHandler struct {
	UserService         *userservice.UserService
	TokenService        *jwtservice.TokenService
	GAuthService        *gauthservice.GoogleAuthService
	TaskService         *taskservice.TaskService
	SubscriptionService *subscriptionservice.SubscriptionService
}

// NewAuthHandler :
func NewAuthHandler(
	userService *userservice.UserService,
	tokenService *jwtservice.TokenService,
	gauthService *gauthservice.GoogleAuthService,
	taskService *taskservice.TaskService,
	subService *subscriptionservice.SubscriptionService,
) *AuthHandler {
	return &AuthHandler{
		UserService:         userService,
		TokenService:        tokenService,
		GAuthService:        gauthService,
		TaskService:         taskService,
		SubscriptionService: subService,
	}
}

// AuthMiddleware :
func (ah *AuthHandler) AuthMiddleware(c *gin.Context) {
	// setting header for auth
	type AuthHeader struct {
		Authorization string `header:"Authorization"`
	}
	var reqHeader AuthHeader
	err := c.BindHeader(&reqHeader)
	if reqHeader.Authorization == "" {
		glog.Error("Auth header missing")
		c.AbortWithStatusJSON(
			http.StatusUnauthorized,
			gin.H{"error": "authorization header missing"},
		)
		return
	}
	token := strings.Split(reqHeader.Authorization, " ")
	tok := token[len(token)-1]
	valid, err := ah.TokenService.IsTokenValid(
		tok, jwtservice.AccessTokenType,
	)
	if !valid {
		glog.Error("Invalid auth token")
		c.AbortWithStatusJSON(
			http.StatusUnauthorized,
			gin.H{"error": "invalid auth token"},
		)
		return
	}
	userID, err := ah.TokenService.GetUserIDFromToken(tok)
	if err != nil {
		glog.Error(err)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}
	_, err = ah.UserService.GetUserByID(userID)
	if err != nil {
		glog.Error(err)
		c.JSON(
			http.StatusUnauthorized,
			gin.H{"error": err.Error()},
		)
		return
	}
	c.Set("user_id", userID)
	c.Next()
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

	userID, isNewUser, err := ah.UserService.Login(request.IDToken)
	if err != nil {
		glog.Error(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "invalid token"})
		return
	}

	if isNewUser {
		// Subscibe to the initial plan
		initialPlanID, err := ah.SubscriptionService.GetSignupPlanID()
		if err != nil {
			glog.Error(err)
			c.JSON(http.StatusInternalServerError, err)
			return
		}
		err = ah.SubscriptionService.CreateSubscripton(
			userID, initialPlanID, uint(0), false,
		)
		if err != nil {
			glog.Error(err)
			c.JSON(http.StatusInternalServerError, err)
			return
		}
		// Create task presets
		err = ah.TaskService.CreatePresetForNewUser(userID)
		if err != nil {
			glog.Error(err)
		}
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
		c.JSON(http.StatusUnauthorized, errors.New("Refresh token is not valid"))
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

// GetUser : get ser details
func (ah *AuthHandler) GetUser(c *gin.Context) {
	type ResponseBody struct {
		FirstName string `json:"first_name"`
		LastName  string `json:"last_name"`
		Email     string `json:"email"`
		Avatar    string `json:"avatar"`
		LeftDays  int    `json:"left_days"`
	}

	// this was set in context from auth middleware
	userID, ok := c.Get("user_id")

	if !ok {
		glog.Error("Could not get the user id from context")
		c.JSON(
			http.StatusInternalServerError,
			gin.H{"error": "could not get the user id from context"},
		)
		return
	}

	user, err := ah.UserService.GetUserByID(userID.(uint))
	if err != nil {
		glog.Error(err)
		c.JSON(
			http.StatusBadRequest,
			gin.H{"error": err.Error()},
		)
		return
	}
	leftDays, err := ah.SubscriptionService.DaysLeftForUser(user.ID)
	if err != nil {
		glog.Error(err)
		// We should not return error in case subscription ends
		// c.JSON(http.StatusInternalServerError, err)
		// return
	}
	resp := ResponseBody{
		FirstName: user.FirstName,
		LastName:  user.LastName,
		Email:     user.Email,
		Avatar:    user.Avatar,
		LeftDays:  leftDays,
	}
	status := http.StatusOK
	if leftDays < 1 {
		status = http.StatusPartialContent
	}
	c.JSON(status, resp)
}
