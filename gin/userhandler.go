package service

import (
	userservice "apidootoday/user"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/golang/glog"
)

// UserHandler :
type UserHandler struct {
	UserService *userservice.UserService
}

// NewUserHandler :
func NewUserHandler(us *userservice.UserService) *UserHandler {
	return &UserHandler{
		UserService: us,
	}
}

// UpdateAutoTaskMove :
func (uh *UserHandler) UpdateAutoTaskMove(c *gin.Context) {
	userID, ok := c.Get("user_id")

	if !ok {
		glog.Error("Could not get the user id from context")
		c.JSON(
			http.StatusInternalServerError,
			gin.H{"error": "could not get the user id from context"},
		)
		return
	}
	type RequestBody struct {
		AllowAutoTaskMove bool `json:"allow_auto_task_move"`
	}
	var request RequestBody
	err := c.BindJSON(&request)
	if err != nil {
		glog.Error("allow_auto_task_move is missing ", err)
		c.JSON(
			http.StatusBadRequest,
			gin.H{"error": "allow_auto_task_move is missing"},
		)
		return
	}
	err = uh.UserService.UpdateAutoTaskMove(
		userID.(uint),
		request.AllowAutoTaskMove,
	)
	if err != nil {
		glog.Error("Could not set allow_auto_task_move - ", err)
		c.JSON(
			http.StatusInternalServerError,
			gin.H{"error": err.Error()},
		)
		return
	}
	c.Status(http.StatusOK)
}

// UpdateUserTimeZoneOffset :
func (uh *UserHandler) UpdateUserTimeZoneOffset(c *gin.Context) {
	userID, ok := c.Get("user_id")

	if !ok {
		glog.Error("Could not get the user id from context")
		c.JSON(
			http.StatusInternalServerError,
			gin.H{"error": "could not get the user id from context"},
		)
		return
	}
	type RequestBody struct {
		TimeZoneOffset int `json:"time_zone_offset"`
	}
	var request RequestBody
	err := c.BindJSON(&request)
	if err != nil || request.TimeZoneOffset == 0 {
		glog.Error("Time zone offset is missing ", err)
		c.JSON(
			http.StatusBadRequest,
			gin.H{"error": "Time zone offset is missing"},
		)
		return
	}
	err = uh.UserService.UpdateUserTimeZoneOffset(
		userID.(uint),
		request.TimeZoneOffset,
	)
	if err != nil {
		glog.Error("Could not set new tz offset - ", err)
		c.JSON(
			http.StatusInternalServerError,
			gin.H{"error": err.Error()},
		)
		return
	}
	c.Status(http.StatusOK)
}
