package service

import (
	"errors"
	"time"
)

var (
	// ErrColumnNotFound : Error when column is not found
	ErrColumnNotFound = errors.New("Column not found")

	// ErrColumnForbidden : Error when column is forbidden
	ErrColumnForbidden = errors.New("Can not update the column")

	// ErrInvalidDateFormat : Error when date format is invalid
	ErrInvalidDateFormat = errors.New("Date format is invalid")
)

// TaskService :
type TaskService struct {
	TDS *TaskDBService
}

// NewTaskService  :
func NewTaskService(tds *TaskDBService) *TaskService {
	return &TaskService{
		TDS: tds,
	}
}

// CreateTask :
func (ts *TaskService) CreateTask(
	markdown string, isDone bool,
	userID uint, columnUUID string,
	date string,
) (Task, error) {
	// If column id is provided
	// Check if the column belongs to the same user
	if columnUUID != "" {
		column, err := ts.TDS.GetColumnByUUID(columnUUID)
		if err != nil {
			return Task{}, ErrColumnNotFound
		}
		if column.UserID != userID {
			return Task{}, ErrColumnForbidden
		}
		return ts.TDS.CreateTaskOnColumn(markdown, isDone, userID, column.ID)
	}
	// Format date
	var formattedDate time.Time
	if date != "" {
		var err error
		formattedDate, err = time.Parse("2006-01-02", date)
		if err != nil {
			return Task{}, ErrInvalidDateFormat
		}
	} else {
		formattedDate, _ = time.Parse("2006-01-02", time.Now().Format("2006-01-02"))
	}
	return ts.TDS.CreateTaskOnDate(markdown, isDone, userID, formattedDate)
}

// CreateColumn :
func (ts *TaskService) CreateColumn(userID uint, name string) (Column, error) {
	return ts.TDS.CreateColumn(userID, name)
}

// UpdateColumn :
func (ts *TaskService) UpdateColumn(colUUID string, name string, userID uint) error {
	col, err := ts.TDS.GetColumnByUUID(colUUID)
	if err != nil {
		return ErrColumnNotFound
	}
	if col.UserID != userID {
		return ErrColumnNotFound
	}
	return ts.TDS.UpdateColumn(col.ID, name)
}

// DeleteColumn :
func (ts *TaskService) DeleteColumn(colUUID string, userID uint) error {
	col, err := ts.TDS.GetColumnByUUID(colUUID)
	if err != nil {
		return ErrColumnNotFound
	}
	if col.UserID != userID {
		return ErrColumnNotFound
	}
	return ts.TDS.DeleteColumn(col.ID)
}

// GetColumns :
func (ts *TaskService) GetColumns(userID uint) ([]Column, error) {
	return ts.TDS.GetColumnsByUserID(userID)
}

// GetColumnByUUID :
func (ts *TaskService) GetColumnByUUID(uuid string, userID uint) (Column, error) {
	col, err := ts.TDS.GetColumnByUUID(uuid)
	if err != nil {
		return col, err
	}
	if col.UserID != userID {
		return col, ErrColumnNotFound
	}
	return col, err
}

// GetTasksByColumnID :
func (ts *TaskService) GetTasksByColumnID(colID uint) ([]Task, error) {
	return ts.TDS.GetTasksByColumn(colID)
}
