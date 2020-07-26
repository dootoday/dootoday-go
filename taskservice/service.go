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

	// ErrTaskNotFound : Error when task is not found
	ErrTaskNotFound = errors.New("Task not found")

	// ErrTaskRepos : Error when task is repositioned
	ErrTaskRepos = errors.New("Task not repositioned")
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

// GetTaskByID :
func (ts *TaskService) GetTaskByID(taskID uint, userID uint) (Task, error) {
	task, err := ts.TDS.GetTaskByID(taskID)
	if err != nil || task.UserID != userID {
		return task, ErrTaskNotFound
	}
	return task, err
}

// DeleteTask :
func (ts *TaskService) DeleteTask(taskID uint, userID uint) error {
	task, err := ts.GetTaskByID(taskID, userID)
	if err != nil {
		return err
	}
	err = ts.TDS.DeleteTaskByID(task.ID)
	if err != nil {
		return err
	}
	return nil
}

// UpdateTask :
func (ts *TaskService) UpdateTask(
	taskID uint, markdown string, isDone bool, userID uint,
) (Task, error) {
	task, err := ts.GetTaskByID(taskID, userID)
	if err != nil {
		return task, err
	}
	if task.Markdown != markdown && markdown != "" {
		err = ts.TDS.UpdateTaskValue(task.ID, markdown)
		if err != nil {
			return task, err
		}
	}
	if task.Done != isDone {
		err = ts.TDS.UpdateTaskStatus(task.ID, isDone)
		if err != nil {
			return task, err
		}
	}
	// ToDo : Need a better approach here
	return ts.GetTaskByID(taskID, userID)
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

// GetColumnByID :
func (ts *TaskService) GetColumnByID(id uint, userID uint) (Column, error) {
	col, err := ts.TDS.GetColumnByID(id)
	if err != nil {
		return col, err
	}
	if col.UserID != userID {
		return col, ErrColumnNotFound
	}
	return col, err
}

// GetTasksByColumnID :
func (ts *TaskService) GetTasksByColumnID(colID uint, userID uint) ([]Task, error) {
	return ts.TDS.GetTasksByColumn(colID, userID)
}

// GetDateRange  :
func (ts *TaskService) GetDateRange(fromDate string, toDate string) ([]time.Time, error) {
	var output []time.Time
	formattedFromDate, err := time.Parse("2006-01-02", fromDate)
	if err != nil {
		return output, ErrInvalidDateFormat
	}
	formattedToDate, err := time.Parse("2006-01-02", toDate)
	if err != nil {
		return output, ErrInvalidDateFormat
	}
	for i := formattedFromDate; i.Before(formattedToDate.Add(time.Hour * 24)); i = i.Add(time.Hour * 24) {
		output = append(output, i)
	}
	return output, nil
}

// FormatDate :
func (ts *TaskService) FormatDate(date string) (time.Time, error) {
	return time.Parse("2006-01-02", date)
}

// FormatDateToString :
func (ts *TaskService) FormatDateToString(date *time.Time) string {
	if date != nil {
		return date.Format("2006-01-02")
	}
	return ""
}

// GetTasksByDate :
func (ts *TaskService) GetTasksByDate(date time.Time, userID uint) ([]Task, error) {
	return ts.TDS.GetTasksByDate(date, userID)
}

// ReposTaskDate :
func (ts *TaskService) ReposTaskDate(taskIDs []uint, date time.Time, userID uint) error {
	err := ts.TDS.VerifyTaskUser(taskIDs, userID)
	if err != nil {
		return ErrTaskNotFound
	}
	err = ts.TDS.ReposTaskDate(taskIDs, date)
	if err != nil {
		return ErrTaskRepos
	}
	return nil
}

// ReposTaskColumn :
func (ts *TaskService) ReposTaskColumn(taskIDs []uint, colUUID string, userID uint) error {
	err := ts.TDS.VerifyTaskUser(taskIDs, userID)
	if err != nil {
		return ErrTaskNotFound
	}
	col, err := ts.GetColumnByUUID(colUUID, userID)
	if err != nil {
		return err
	}
	err = ts.TDS.ReposTaskColumn(taskIDs, col.ID)
	if err != nil {
		return ErrTaskRepos
	}
	return nil
}
