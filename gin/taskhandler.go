package service

import (
	redisclient "apidootoday/redisclient"
	taskservice "apidootoday/taskservice"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/golang/glog"
)

// TaskHandler :
type TaskHandler struct {
	TaskService *taskservice.TaskService
	RedisClient *redisclient.RedisClient
}

// NewTaskHandler :
func NewTaskHandler(
	ts *taskservice.TaskService,
	rc *redisclient.RedisClient,
) *TaskHandler {
	return &TaskHandler{
		TaskService: ts,
		RedisClient: rc,
	}
}

// TaskResponse :
type TaskResponse struct {
	ID         uint   `json:"id"`
	Markdown   string `json:"markdown"`
	IsDone     bool   `json:"is_done"`
	ColumnUUID string `json:"column_id"`
	Date       string `json:"date"`
	Order      int    `json:"order"`
}

// ColumnResponse :
type ColumnResponse struct {
	UUID     string         `json:"id"`
	Name     string         `json:"name"`
	MetaText string         `json:"meta"`
	Tasks    []TaskResponse `json:"tasks"`
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
		Date:       th.TaskService.FormatDateToString(task.Date),
		Order:      task.Order,
	}
	_, err = th.RedisClient.SetUserLastUpdate(userID.(uint))
	if err != nil {
		glog.Error("Could not set the last updated to the cache", err)
	}
	c.JSON(http.StatusOK, taskResp)
}

// UpdateTask :
func (th *TaskHandler) UpdateTask(c *gin.Context) {
	userID, ok := c.Get("user_id")

	if !ok {
		glog.Error("Could not get the user id from context")
		c.JSON(
			http.StatusInternalServerError,
			gin.H{"error": "could not get the user id from context"},
		)
		return
	}
	tID := c.Param("task_id")
	taskID, err := strconv.ParseUint(tID, 10, 32)
	if err != nil {
		c.JSON(
			http.StatusBadRequest,
			gin.H{"error": "Invalid task ID"},
		)
		return
	}
	type RequestBody struct {
		Markdown string `json:"markdown"`
		IsDone   bool   `json:"is_done"`
	}
	var request RequestBody
	err = c.BindJSON(&request)
	if err != nil {
		c.JSON(
			http.StatusBadRequest,
			gin.H{"error": "Can not bind the request body"},
		)
		return
	}

	task, err := th.TaskService.UpdateTask(
		uint(taskID), request.Markdown, request.IsDone, userID.(uint),
	)
	if err != nil {
		c.JSON(
			http.StatusBadRequest,
			gin.H{"error": err.Error()},
		)
		return
	}
	_, err = th.RedisClient.SetUserLastUpdate(userID.(uint))
	if err != nil {
		glog.Error("Could not set the last updated to the cache", err)
	}
	col := ""
	if task.ColumnID != 0 {
		column, err := th.TaskService.GetColumnByID(task.ColumnID, userID.(uint))
		if err != nil {
			glog.Error("Error getting the column", err)
			c.JSON(
				http.StatusBadRequest,
				gin.H{"error": err.Error()},
			)
			return
		}
		col = column.UUID
	}
	taskResp := TaskResponse{
		ID:         task.ID,
		Markdown:   task.Markdown,
		IsDone:     task.Done,
		ColumnUUID: col,
		Date:       th.TaskService.FormatDateToString(task.Date),
		Order:      task.Order,
	}
	c.JSON(http.StatusOK, taskResp)
}

// GetTask :
func (th *TaskHandler) GetTask(c *gin.Context) {
	userID, ok := c.Get("user_id")

	if !ok {
		glog.Error("Could not get the user id from context")
		c.JSON(
			http.StatusInternalServerError,
			gin.H{"error": "could not get the user id from context"},
		)
		return
	}
	tID := c.Param("task_id")
	taskID, err := strconv.ParseUint(tID, 10, 32)
	if err != nil {
		c.JSON(
			http.StatusBadRequest,
			gin.H{"error": "Invalid task ID"},
		)
		return
	}

	task, err := th.TaskService.GetTaskByID(uint(taskID), userID.(uint))
	if err != nil {
		c.JSON(
			http.StatusBadRequest,
			gin.H{"error": err.Error()},
		)
		return
	}
	col := ""
	if task.ColumnID != 0 {
		column, err := th.TaskService.GetColumnByID(task.ColumnID, userID.(uint))
		if err != nil {
			glog.Error("Error getting the column", err)
			c.JSON(
				http.StatusBadRequest,
				gin.H{"error": err.Error()},
			)
			return
		}
		col = column.UUID
	}
	taskResp := TaskResponse{
		ID:         task.ID,
		Markdown:   task.Markdown,
		IsDone:     task.Done,
		ColumnUUID: col,
		Date:       th.TaskService.FormatDateToString(task.Date),
		Order:      task.Order,
	}
	c.JSON(http.StatusOK, taskResp)
}

// DeleteTask :
func (th *TaskHandler) DeleteTask(c *gin.Context) {
	userID, ok := c.Get("user_id")

	if !ok {
		glog.Error("Could not get the user id from context")
		c.JSON(
			http.StatusInternalServerError,
			gin.H{"error": "could not get the user id from context"},
		)
		return
	}
	tID := c.Param("task_id")
	taskID, err := strconv.ParseUint(tID, 10, 32)
	if err != nil {
		c.JSON(
			http.StatusBadRequest,
			gin.H{"error": "Invalid task ID"},
		)
		return
	}

	err = th.TaskService.DeleteTask(uint(taskID), userID.(uint))
	if err != nil {
		c.JSON(
			http.StatusBadRequest,
			gin.H{"error": err.Error()},
		)
		return
	}
	_, err = th.RedisClient.SetUserLastUpdate(userID.(uint))
	if err != nil {
		glog.Error("Could not set the last updated to the cache", err)
	}
	c.JSON(http.StatusOK, gin.H{"success": "ok"})
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
	_, err = th.RedisClient.SetUserLastUpdate(userID.(uint))
	if err != nil {
		glog.Error("Could not set the last updated to the cache", err)
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
	_, err = th.RedisClient.SetUserLastUpdate(userID.(uint))
	if err != nil {
		glog.Error("Could not set the last updated to the cache", err)
	}
	c.JSON(http.StatusOK, gin.H{"success": "ok"})
}

// DeleteColumn :
func (th *TaskHandler) DeleteColumn(c *gin.Context) {
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

	err := th.TaskService.DeleteColumn(colUUID, userID.(uint))
	if err != nil {
		c.JSON(
			http.StatusBadRequest,
			gin.H{"error": err.Error()},
		)
		return
	}
	_, err = th.RedisClient.SetUserLastUpdate(userID.(uint))
	if err != nil {
		glog.Error("Could not set the last updated to the cache", err)
	}
	c.JSON(http.StatusOK, gin.H{"success": "ok"})
}

// GetColumns :
func (th *TaskHandler) GetColumns(c *gin.Context) {
	userID, ok := c.Get("user_id")

	if !ok {
		glog.Error("Could not get the user id from context")
		c.JSON(
			http.StatusInternalServerError,
			gin.H{"error": "could not get the user id from context"},
		)
		return
	}
	cols, err := th.TaskService.GetColumns(userID.(uint))
	if err != nil {
		c.JSON(
			http.StatusInternalServerError,
			gin.H{"error": err.Error()},
		)
		return
	}
	colresp := []ColumnResponse{}

	for _, col := range cols {
		tasks, err := th.TaskService.GetTasksByColumnID(col.ID, userID.(uint))
		if err != nil {
			c.JSON(
				http.StatusInternalServerError,
				gin.H{"error": err.Error()},
			)
			return
		}
		taskresp := []TaskResponse{}
		for _, task := range tasks {
			taskresp = append(
				taskresp,
				TaskResponse{
					ID:         task.ID,
					Markdown:   task.Markdown,
					IsDone:     task.Done,
					ColumnUUID: col.UUID,
					Date:       th.TaskService.FormatDateToString(task.Date),
					Order:      task.Order,
				},
			)
		}
		colresp = append(
			colresp,
			ColumnResponse{
				UUID:  col.UUID,
				Name:  col.Name,
				Tasks: taskresp,
			},
		)
	}

	c.JSON(http.StatusOK, colresp)
}

// GetColumn :
func (th *TaskHandler) GetColumn(c *gin.Context) {
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

	col, err := th.TaskService.GetColumnByUUID(colUUID, userID.(uint))
	if err != nil {
		c.JSON(
			http.StatusInternalServerError,
			gin.H{"error": err.Error()},
		)
		return
	}
	tasks, err := th.TaskService.GetTasksByColumnID(col.ID, userID.(uint))
	if err != nil {
		c.JSON(
			http.StatusInternalServerError,
			gin.H{"error": err.Error()},
		)
		return
	}
	taskresp := []TaskResponse{}
	for _, task := range tasks {
		taskresp = append(
			taskresp,
			TaskResponse{
				ID:         task.ID,
				Markdown:   task.Markdown,
				IsDone:     task.Done,
				ColumnUUID: col.UUID,
				Date:       th.TaskService.FormatDateToString(task.Date),
				Order:      task.Order,
			},
		)
	}
	colresp := ColumnResponse{
		UUID:  col.UUID,
		Name:  col.Name,
		Tasks: taskresp,
	}

	c.JSON(http.StatusOK, colresp)
}

// GetTasks :
func (th *TaskHandler) GetTasks(c *gin.Context) {
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
		From string `form:"from"`
		To   string `form:"to"`
	}
	var request RequestBody
	err := c.Bind(&request)
	if err != nil {
		c.JSON(
			http.StatusBadRequest,
			gin.H{"error": "Invalid data"},
		)
		return
	}
	if request.From == "" {
		glog.Error("From date is missing")
		c.JSON(
			http.StatusBadRequest,
			gin.H{"error": "From date is missing"},
		)
		return
	}
	if request.To == "" {
		glog.Error("To date is missing")
		c.JSON(
			http.StatusBadRequest,
			gin.H{"error": "To date is missing"},
		)
		return
	}
	dates, err := th.TaskService.GetDateRange(request.From, request.To)
	if err != nil {
		glog.Error("Something  is wrong with date range", err)
		c.JSON(
			http.StatusBadRequest,
			gin.H{"error": "Can not find a date range"},
		)
		return
	}
	colResp := []ColumnResponse{}
	for _, date := range dates {
		tasks, err := th.TaskService.GetTasksByDate(date, userID.(uint))
		if err != nil {
			glog.Error("Could not fetch tasks for the date - ", err)
			c.JSON(
				http.StatusBadGateway,
				gin.H{"error": "Could not fetch tasks for the date"},
			)
			return
		}
		taskResp := []TaskResponse{}
		for _, task := range tasks {
			taskResp = append(
				taskResp,
				TaskResponse{
					ID:       task.ID,
					Markdown: task.Markdown,
					IsDone:   task.Done,
					Date:     th.TaskService.FormatDateToString(task.Date),
					Order:    task.Order,
				},
			)
		}
		colResp = append(
			colResp,
			ColumnResponse{
				Name:     date.Weekday().String(),
				MetaText: date.Format("2006-01-02"),
				UUID:     date.Format("2006-01-02"),
				Tasks:    taskResp,
			},
		)
	}
	c.JSON(http.StatusOK, colResp)
}

// ReposTask :
func (th *TaskHandler) ReposTask(c *gin.Context) {
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
		ColumnUUID string `json:"column_id"`
		Date       string `json:"date"`
		TaskIDs    []uint `json:"task_ids"`
	}

	var request RequestBody
	err := c.BindJSON(&request)
	if err != nil || (request.ColumnUUID == "" && request.Date == "") {
		glog.Error("column_id or date is needed", err)
		c.JSON(
			http.StatusBadRequest,
			gin.H{"error": "column_id or date is needed"},
		)
		return
	}

	if request.Date != "" {
		date, err := th.TaskService.FormatDate(request.Date)
		if err != nil {
			c.JSON(
				http.StatusBadRequest,
				gin.H{"error": "Date format is invalid"},
			)
			return
		}
		err = th.TaskService.ReposTaskDate(
			request.TaskIDs,
			date,
			userID.(uint),
		)
		if err != nil {
			c.JSON(
				http.StatusBadGateway,
				gin.H{"error": "Could not resposition the task"},
			)
			return
		}
	} else {
		err = th.TaskService.ReposTaskColumn(
			request.TaskIDs,
			request.ColumnUUID,
			userID.(uint),
		)
		if err != nil {
			c.JSON(
				http.StatusBadGateway,
				gin.H{"error": "Could not resposition the task"},
			)
			return
		}
	}
	_, err = th.RedisClient.SetUserLastUpdate(userID.(uint))
	if err != nil {
		glog.Error("Could not set the last updated to the cache", err)
	}
	c.JSON(http.StatusOK, gin.H{"success": "ok"})
}

// GetLastUpdated :
func (th *TaskHandler) GetLastUpdated(c *gin.Context) {
	userID, ok := c.Get("user_id")

	if !ok {
		glog.Error("Could not get the user id from context")
		c.JSON(
			http.StatusInternalServerError,
			gin.H{"error": "could not get the user id from context"},
		)
		return
	}
	val, err := th.RedisClient.GetUserLastUpdate(userID.(uint))
	if err != nil {
		glog.Error("Could not get the last updated from the cache", err)
		c.JSON(http.StatusBadGateway, gin.H{"last_updated": ""})
		return
	}
	c.JSON(http.StatusOK, gin.H{"last_updated": val})
}
