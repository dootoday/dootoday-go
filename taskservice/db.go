package service

import (
	"errors"
	"time"

	"github.com/golang/glog"
	"github.com/google/uuid"
	"github.com/jinzhu/gorm"
)

// Task :
type Task struct {
	gorm.Model
	UserID   uint `gorm:"index:usertask"`
	ColumnID uint `gorm:"index:columntask"`
	Markdown string
	Order    int
	Done     bool
	Date     *time.Time `gorm:"default:NULL"`
}

// Column :
type Column struct {
	gorm.Model
	UUID   string `gorm:"index:columnuuid"`
	UserID uint   `gorm:"index:usercolumn"`
	Name   string
}

// TaskDBService :
type TaskDBService struct {
	DB *gorm.DB
}

// NewTaskDBService  :
func NewTaskDBService(db *gorm.DB) *TaskDBService {
	return &TaskDBService{
		DB: db,
	}
}

// Migrate :
func (ts *TaskDBService) Migrate() error {
	glog.Info("Creating tasks table")
	err := ts.DB.AutoMigrate(&Task{}).Error
	if err != nil {
		glog.Error(err)
	}
	glog.Info("Creating columns table")
	err = ts.DB.AutoMigrate(&Column{}).Error
	if err != nil {
		glog.Error(err)
	}
	return nil
}

// CreateTaskOnColumn :
func (ts *TaskDBService) CreateTaskOnColumn(
	markdown string, isDone bool, userID uint, columnID uint,
) (Task, error) {
	order := 1
	tasks, err := ts.GetTasksByColumn(columnID, userID)
	if err != nil {
		return Task{}, err
	}
	order = len(tasks) + 1
	newTask := Task{
		UserID:   userID,
		ColumnID: columnID,
		Markdown: markdown,
		Done:     isDone,
		Order:    order,
	}
	err = ts.DB.Create(&newTask).Error
	return newTask, err
}

// CreateTaskOnDate :
func (ts *TaskDBService) CreateTaskOnDate(
	markdown string, isDone bool, userID uint, date time.Time,
) (Task, error) {
	order := 1
	tasks, err := ts.GetTasksByDate(date, userID)
	if err != nil {
		return Task{}, err
	}
	order = len(tasks) + 1
	newTask := Task{
		UserID:   userID,
		Markdown: markdown,
		Done:     isDone,
		Date:     &date,
		Order:    order,
	}
	err = ts.DB.Create(&newTask).Error
	return newTask, err
}

// GetTasksByColumn :
func (ts *TaskDBService) GetTasksByColumn(columnID uint, userID uint) ([]Task, error) {
	var tasks []Task
	err := ts.DB.Where("column_id=? AND user_id=?", columnID, userID).Order("order").Find(&tasks).Error
	if err == gorm.ErrRecordNotFound {
		err = nil
	}
	return tasks, err
}

// GetTasksByDate :
func (ts *TaskDBService) GetTasksByDate(date time.Time, userID uint) ([]Task, error) {
	var tasks []Task
	err := ts.DB.Where("date=? AND user_id=?", date, userID).Order("order").Find(&tasks).Error
	if err == gorm.ErrRecordNotFound {
		err = nil
	}
	return tasks, err
}

// GetTasksByIDs :
func (ts *TaskDBService) GetTasksByIDs(ids []uint) ([]Task, error) {
	var tasks []Task
	err := ts.DB.Where("id IN (?)", ids).Find(&tasks).Error
	if err == gorm.ErrRecordNotFound {
		err = nil
	}
	return tasks, err
}

// GetTaskByID :
func (ts *TaskDBService) GetTaskByID(id uint) (Task, error) {
	var task Task
	err := ts.DB.Where("id=?", id).Find(&task).Error
	return task, err
}

// UpdateTaskValue :
func (ts *TaskDBService) UpdateTaskValue(
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
func (ts *TaskDBService) UpdateTaskStatus(
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
func (ts *TaskDBService) DeleteTaskByID(taskID uint) error {
	return ts.DB.Where("id=?", taskID).Delete(&Task{}).Error
}

// VerifyTaskUser :
func (ts *TaskDBService) VerifyTaskUser(taskIDs []uint, userID uint) error {
	tasks := []Task{}
	err := ts.DB.Where("id IN (?) AND user_id=?", taskIDs, userID).Find(&tasks).Error
	if err != nil {
		return err
	}
	if len(tasks) != len(taskIDs) {
		return errors.New("Forbidden task ID")
	}
	return nil
}

// ReposTaskDate :
func (ts *TaskDBService) ReposTaskDate(
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
func (ts *TaskDBService) ReposTaskColumn(
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
func (ts *TaskDBService) CreateColumn(userID uint, name string) (Column, error) {
	newuuid := uuid.New().String()
	newcol := Column{
		UUID:   newuuid,
		Name:   name,
		UserID: userID,
	}
	err := ts.DB.Create(&newcol).Error
	if err != nil {
		return Column{}, err
	}
	return newcol, err
}

// GetColumnsByUserID :
func (ts *TaskDBService) GetColumnsByUserID(userID uint) ([]Column, error) {
	var columns []Column
	err := ts.DB.Where("user_id=?", userID).Find(&columns).Error
	if err == gorm.ErrRecordNotFound {
		err = nil
	}
	return columns, err
}

// GetColumnByID :
func (ts *TaskDBService) GetColumnByID(colID uint) (Column, error) {
	var column Column
	err := ts.DB.Where("id=?", colID).Find(&column).Error
	return column, err
}

// GetColumnByUUID :
func (ts *TaskDBService) GetColumnByUUID(colID string) (Column, error) {
	var column Column
	err := ts.DB.Where("uuid=?", colID).Find(&column).Error
	return column, err
}

// UpdateColumn :
func (ts *TaskDBService) UpdateColumn(columnID uint, name string) error {
	column, err := ts.GetColumnByID(columnID)
	if err != nil {
		return err
	}
	err = ts.DB.Model(&column).Update(map[string]interface{}{"name": name}).Error
	return err
}

// DeleteColumn :
func (ts *TaskDBService) DeleteColumn(colID uint) error {
	return ts.DB.Where("id=?", colID).Delete(&Column{}).Error
}
