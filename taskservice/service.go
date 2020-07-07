package service

import (
	"time"

	"github.com/google/uuid"
	"github.com/jinzhu/gorm"
)

// TaskService :
type TaskService struct {
	DB *gorm.DB
}

// NewTaskService  :
func NewTaskService(db *gorm.DB) *TaskService {
	return &TaskService{
		DB: db,
	}
}

// CreateTask :
func (ts *TaskService) CreateTask(
	markdown string, userID uint, columnID uint, date time.Time,
) (uint, error) {
	order := 1
	if columnID > 0 {
		tasks, err := ts.GetTasksByColumn(columnID)
		if err != nil {
			return uint(0), err
		}
		order = len(tasks) + 1
	} else {
		tasks, err := ts.GetTasksByDate(date)
		if err != nil {
			return uint(0), err
		}
		order = len(tasks) + 1
	}
	newTask := Task{
		UserID:   userID,
		ColumnID: columnID,
		Markdown: markdown,
		Done:     true,
		Date:     date,
		Order:    order,
	}
	err := ts.DB.Create(&newTask).Error
	return newTask.ID, err
}

// GetTasksByColumn :
func (ts *TaskService) GetTasksByColumn(columnID uint) ([]Task, error) {
	var tasks []Task
	err := ts.DB.Where("column_id=?", columnID).Order("order").Find(&tasks).Error
	if err == gorm.ErrRecordNotFound {
		err = nil
	}
	return tasks, err
}

// GetTasksByDate :
func (ts *TaskService) GetTasksByDate(date time.Time) ([]Task, error) {
	var tasks []Task
	err := ts.DB.Where("date=?", date).Order("order").Find(&tasks).Error
	if err == gorm.ErrRecordNotFound {
		err = nil
	}
	return tasks, err
}

// GetTasksByIDs :
func (ts *TaskService) GetTasksByIDs(ids []uint) ([]Task, error) {
	var tasks []Task
	err := ts.DB.Where("id IN (?)", ids).Find(&tasks).Error
	if err == gorm.ErrRecordNotFound {
		err = nil
	}
	return tasks, err
}

// GetTaskByID :
func (ts *TaskService) GetTaskByID(id uint) (Task, error) {
	var task Task
	err := ts.DB.Where("id=?", id).Find(&task).Error
	return task, err
}

// UpdateTaskValue :
func (ts *TaskService) UpdateTaskValue(
	taskID uint, markdown string,
) error {
	task, err := ts.GetTaskByID(taskID)
	if err != nil {
		return err
	}
	err = ts.DB.Model(&task).Update(map[string]interface{}{"markdown": markdown}).Error
	return err
}

// UpdateTaskStatus :
func (ts *TaskService) UpdateTaskStatus(
	taskID uint, done bool,
) error {
	task, err := ts.GetTaskByID(taskID)
	if err != nil {
		return err
	}
	err = ts.DB.Model(&task).Update(map[string]interface{}{"done": done}).Error
	return err
}

// DeleteTaskByID :
func (ts *TaskService) DeleteTaskByID(taskID uint) error {
	return ts.DB.Where("id=?", taskID).Delete(&Task{}).Error
}

// ReposTaskDate :
func (ts *TaskService) ReposTaskDate(
	taskIDs []uint, date time.Time,
) error {
	tx := ts.DB.Begin()
	for idx, taskID := range taskIDs {
		task, err := ts.GetTaskByID(taskID)
		if err != nil {
			tx.Rollback()
			return err
		}
		tx.Model(&task).Update(map[string]interface{}{
			"column_id": nil,
			"order":     idx + 1,
			"date":      date,
		})
	}
	return tx.Commit().Error
}

// ReposTaskColumn :
func (ts *TaskService) ReposTaskColumn(
	taskIDs []uint, columnID uint,
) error {
	tx := ts.DB.Begin()
	for idx, taskID := range taskIDs {
		task, err := ts.GetTaskByID(taskID)
		if err != nil {
			tx.Rollback()
			return err
		}
		tx.Model(&task).Update(map[string]interface{}{
			"column_id": columnID,
			"order":     idx + 1,
			"date":      nil,
		})
	}
	return tx.Commit().Error
}

// CreateColumn :
func (ts *TaskService) CreateColumn(userID uint, name string) (uint, error) {
	newuuid := uuid.New().String()
	newcol := Column{
		UUID: newuuid,
		Name: name,
	}
	err := ts.DB.Create(&newcol).Error
	if err != nil {
		return uint(0), err
	}
	return newcol.ID, err
}

// GetColumnsByUserID :
func (ts *TaskService) GetColumnsByUserID(userID uint) ([]Column, error) {
	var columns []Column
	err := ts.DB.Where("user_id=?", userID).Find(&columns).Error
	return columns, err
}

// GetColumnByID :
func (ts *TaskService) GetColumnByID(colID uint) (Column, error) {
	var column Column
	err := ts.DB.Where("id=?", colID).Find(&column).Error
	return column, err
}

// UpdateColumn :
func (ts *TaskService) UpdateColumn(columnID uint, name string) error {
	column, err := ts.GetColumnByID(columnID)
	if err != nil {
		return err
	}
	err = ts.DB.Model(&column).Update(map[string]interface{}{"name": name}).Error
	return err
}

// DeleteColumn :
func (ts *TaskService) DeleteColumn(colID uint) error {
	return ts.DB.Where("id=?", colID).Delete(&Column{}).Error
}
