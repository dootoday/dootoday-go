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
