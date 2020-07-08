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

// ColumnResponse :
type ColumnResponse struct {
	UUID  string         `json:"id"`
	Name  string         `json:"name"`
	Tasks []TaskResponse `json:"tasks"`
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

// CreateColumn :
func (th *TaskHandler) CreateColumn(c *gin.Context) {
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
		Name string `json:"name"`
	}

	var request RequestBody
	err := c.BindJSON(&request)
	if err != nil || request.Name == "" {
		glog.Error("Column name is missing", err)
		c.JSON(
			http.StatusBadRequest,
			gin.H{"error": "Column name is missing"},
		)
		return
	}
	col, err := th.TaskService.CreateColumn(userID.(uint), request.Name)
	if err != nil {
		c.JSON(
			http.StatusBadRequest,
			gin.H{"error": err.Error()},
		)
		return
	}

	colresp := ColumnResponse{
		Name:  col.Name,
		UUID:  col.UUID,
		Tasks: []TaskResponse{},
	}

	c.JSON(http.StatusOK, colresp)
}

// UpdateColumn :
func (th *TaskHandler) UpdateColumn(c *gin.Context) {
	userID, ok := c.Get("user_id")

	if !ok {
		glog.Error("Could not get the user id from context")
		c.JSON(
			http.StatusInternalServerError,
			gin.H{"error": "could not get the user id from context"},
		)
		return
	}

	colUUID := c.Param("col_id")

	type RequestBody struct {
		Name string `json:"name"`
	}

	var request RequestBody
	err := c.BindJSON(&request)
	if err != nil || request.Name == "" {
		glog.Error("Column name is missing", err)
		c.JSON(
			http.StatusBadRequest,
			gin.H{"error": "Column name is missing"},
		)
		return
	}
	err = th.TaskService.UpdateColumn(colUUID, request.Name, userID.(uint))
	if err != nil {
		c.JSON(
			http.StatusBadRequest,
			gin.H{"error": err.Error()},
		)
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": "ok"})
}
