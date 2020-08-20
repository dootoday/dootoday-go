package service

import (
	"errors"
	"strings"
	"time"

	"github.com/golang/glog"
	"github.com/google/uuid"
	"github.com/jinzhu/gorm"
)

// Task :
type Task struct {
	gorm.Model
	UserID        uint `gorm:"index:usertask"`
	ColumnID      uint `gorm:"index:columntask"`
	Markdown      string
	Order         int
	Done          bool
	RecurringType RecurringType `gorm:"default:'none'"`
	Date          string        `gorm:"type:date;default:NULL"`
}

// Column :
type Column struct {
	gorm.Model
	UUID   string `gorm:"index:columnuuid"`
	UserID uint   `gorm:"index:usercolumn"`
	Name   string
}

// RecurringTaskStatus :
type RecurringTaskStatus struct {
	gorm.Model
	Date   string `gorm:"type:date;index:recurring_taskstatus_date"`
	TaskID uint   `gorm:"index:recurring_taskstatus_task_id"`
	Done   bool
	Order  int
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
	glog.Info("Creating recurring task status table")
	err = ts.DB.AutoMigrate(&RecurringTaskStatus{}).Error
	if err != nil {
		glog.Error(err)
	}
	// This is temporary
	// Change the date column type
	err = ts.DB.Exec(`ALTER TABLE tasks MODIFY COLUMN date date;`).Error
	if err != nil {
		glog.Info(err)
	}
	err = ts.DB.Exec(`ALTER TABLE recurring_task_statuses MODIFY COLUMN date date;`).Error
	if err != nil {
		glog.Info(err)
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
	markdown string, isDone bool, userID uint, date time.Time, recurringType RecurringType,
) (Task, error) {
	tasks, err := ts.GetTasksByDate(date, userID)
	if err != nil {
		return Task{}, err
	}
	recTaskCount := ts.GetRecurringTaskCountByDate(date, userID)
	order := len(tasks) + recTaskCount + 1
	newTask := Task{
		UserID:        userID,
		Markdown:      markdown,
		Done:          isDone,
		Date:          date.Format("2006-01-02"),
		Order:         order,
		RecurringType: recurringType,
	}
	err = ts.DB.Create(&newTask).Error
	if err != nil {
		return newTask, err
	}

	if recurringType != RecurringNone {
		// Create a status entry for the date
		ts.CreateRecurringTaskStatus(newTask.ID, newTask.Date, order, isDone)
	}
	// Format task date
	newTask.Date = ts.FormatDateString(newTask.Date)
	return newTask, err
}

// GetTasksByColumn :
func (ts *TaskDBService) GetTasksByColumn(columnID uint, userID uint) ([]Task, error) {
	var tasks []Task
	err := ts.DB.Where("column_id=? AND user_id=? AND recurring_type='none'", columnID, userID).Order("order").Find(&tasks).Error
	if err == gorm.ErrRecordNotFound {
		err = nil
	}
	return tasks, err
}

// GetTasksByDate :
func (ts *TaskDBService) GetTasksByDate(date time.Time, userID uint) ([]Task, error) {
	var tasks []Task
	err := ts.DB.Where(
		"date=? AND user_id=? AND recurring_type='none'",
		date.Format("2006-01-02"), userID).Order("order").Find(&tasks).Error
	if err == gorm.ErrRecordNotFound {
		err = nil
	}
	if err == nil {
		// Format all task dates
		for idx := range tasks {
			tasks[idx].Date = ts.FormatDateString(tasks[idx].Date)
		}
	}
	return tasks, err
}

// GetRecurringTasks :
func (ts *TaskDBService) GetRecurringTasks(userID uint) ([]Task, error) {
	var tasks []Task
	err := ts.DB.Where("user_id=? AND recurring_type!='none'", userID).Order("order").Find(&tasks).Error
	if err == gorm.ErrRecordNotFound {
		err = nil
	}
	if err == nil {
		// Format all task dates
		for idx := range tasks {
			tasks[idx].Date = ts.FormatDateString(tasks[idx].Date)
		}
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
	if err == nil {
		// Format all task dates
		for idx := range tasks {
			tasks[idx].Date = ts.FormatDateString(tasks[idx].Date)
		}
	}
	return tasks, err
}

// GetTaskByID :
func (ts *TaskDBService) GetTaskByID(id uint) (Task, error) {
	var task Task
	err := ts.DB.Where("id=?", id).Find(&task).Error
	if err == nil {
		// Format task date
		task.Date = ts.FormatDateString(task.Date)
	}
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
	if err == nil {
		// Format task date
		task.Date = ts.FormatDateString(task.Date)
	}
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
	if err == nil {
		// Format task date
		task.Date = ts.FormatDateString(task.Date)
	}
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
	idx := 0
	for _, taskID := range taskIDs {
		task, err := ts.GetTaskByID(taskID)
		if err != nil {
			tx.Rollback()
			return err
		}
		if task.RecurringType == RecurringNone {
			tx.Model(&task).Update(map[string]interface{}{
				"column_id": nil,
				"order":     idx,
				"date":      date.Format("2006-01-02"),
			})
			idx = idx + 1
		} else {
			// for recurring task
			// if the task belong to the same date
			// then only make some changes
			rts, err := ts.FindOrCreateRecurringTaskStatus(
				task.ID,
				date.Format("2006-01-02"),
			)
			if err != nil {
				tx.Rollback()
				return err
			}
			// The date from DB looks like 2020-08-20T00:00:00-05:00
			// We need this in 2020-08-20 format
			// If the recurring task is of the same day
			// Only then update the order
			if date.Format("2006-01-02") == rts.Date {
				rts.Order = idx
				err := tx.Save(&rts).Error
				if err != nil {
					tx.Rollback()
					return err
				}
				idx = idx + 1
			}
		}
	}
	return tx.Commit().Error
}

// FormatDateString :
// The date from DB looks like 2020-08-20T00:00:00-05:00
// We need this in 2020-08-20 format
func (ts *TaskDBService) FormatDateString(date string) string {
	splitted := strings.Split(date, "T")
	if len(splitted) == 2 {
		return splitted[0]
	}
	return date
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

// FindOrCreateRecurringTaskStatus :
func (ts *TaskDBService) FindOrCreateRecurringTaskStatus(
	taskID uint,
	date string,
) (RecurringTaskStatus, error) {
	rts := RecurringTaskStatus{
		TaskID: taskID,
		Date:   date,
	}
	err := ts.DB.Where("task_id=? AND date=?", taskID, date).FirstOrCreate(&rts).Error
	if err == nil {
		rts.Date = ts.FormatDateString(rts.Date)
	}
	return rts, err
}

// CreateRecurringTaskStatus :
func (ts *TaskDBService) CreateRecurringTaskStatus(
	taskID uint,
	date string,
	order int,
	isDone bool,
) (RecurringTaskStatus, error) {
	rts := RecurringTaskStatus{
		TaskID: taskID,
		Date:   date,
		Order:  order,
		Done:   isDone,
	}
	err := ts.DB.Create(&rts).Error
	if err == nil {
		// Format task date
		rts.Date = ts.FormatDateString(rts.Date)
	}
	return rts, err
}

// GetRecurringTaskStatusByID :
func (ts *TaskDBService) GetRecurringTaskStatusByID(recurringID uint) (RecurringTaskStatus, error) {
	rts := RecurringTaskStatus{}
	err := ts.DB.Where("id=?", recurringID).First(&rts).Error
	if err == nil {
		// Format task date
		rts.Date = ts.FormatDateString(rts.Date)
	}
	return rts, err
}

// UpdateRecurringTaskStatus :
func (ts *TaskDBService) UpdateRecurringTaskStatus(
	rts RecurringTaskStatus,
) error {
	err := ts.DB.Save(&rts).Error
	if err == nil {
		// Format task date
		rts.Date = ts.FormatDateString(rts.Date)
	}
	return err
}

// GetRecurringTaskCountByDate : This function totally depend on the data
// In RecurringTaskStatus table. This function should be called for a date,
// Only when that table is populated for the date
func (ts *TaskDBService) GetRecurringTaskCountByDate(
	date time.Time, userID uint,
) int {
	recurringTasks, err := ts.GetRecurringTasks(userID)
	if err != nil || len(recurringTasks) == 0 {
		return 0
	}
	taskIDs := []uint{}
	for _, task := range recurringTasks {
		taskIDs = append(taskIDs, task.ID)
	}
	rts := []RecurringTaskStatus{}
	err = ts.DB.Where("task_id IN (?) AND date=?",
		taskIDs, date.Format("2006-01-02")).Find(&rts).Error
	if err != nil {
		return 0
	}
	return len(rts)
}
