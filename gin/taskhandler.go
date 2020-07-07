package service

import (
	taskservice "apidootoday/taskservice"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang/glog"
)

// TaskHandler :
type TaskHandler struct {
	TaskService *taskservice.TaskService
}

// NewTaskHandler :
func NewTaskHandler(ts *taskservice.TaskService) *TaskHandler {
	return &TaskHandler{
		TaskService: ts,
	}
}

// TaskResponse :
type TaskResponse struct {
	ID         uint       `json:"id"`
	Markdown   string     `json:"markdown"`
	IsDone     bool       `json:"is_done"`
	ColumnUUID string     `json:"column_id"`
	Date       *time.Time `json:"date"`
	Order      int        `json:"order"`
}

// CreateTask :
func (th *TaskHandler) CreateTask(c *gin.Context) {
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
		Markdown   string `json:"markdown"`
		IsDone     bool   `json:"is_done"`
		ColumnUUID string `json:"column_id"`
		Date       string `json:"date"`
	}

	var request RequestBody
	err := c.BindJSON(&request)
	if err != nil || request.Markdown == "" {
		glog.Error("Task content is missing ", err)
		c.JSON(
			http.StatusBadRequest,
			gin.H{"error": "Task content is missing"},
		)
		return
	}
	task, err := th.TaskService.CreateTask(
		request.Markdown,
		request.IsDone,
		userID.(uint),
		request.ColumnUUID,
		request.Date,
	)
	if err != nil {
		c.JSON(
			http.StatusBadRequest,
			gin.H{"error": err.Error()},
		)
		return
	}
	taskResp := TaskResponse{
		ID:         task.ID,
		Markdown:   task.Markdown,
		IsDone:     task.Done,
		ColumnUUID: request.ColumnUUID,
		Date:       task.Date,
		Order:      task.Order,
	}
	c.JSON(http.StatusOK, taskResp)
}
