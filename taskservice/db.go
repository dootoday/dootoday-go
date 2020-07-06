package service

import (
	"time"

	"github.com/golang/glog"
	"github.com/jinzhu/gorm"
)

// Task :
type Task struct {
	gorm.Model
	UserID   uint `gorm:"index:usertask"`
	ColumnID uint `gorm:"index:columntask"`
	Markdown string
	Order    int
	Done     bool      `gorm:"default:NULL"`
	Date     time.Time `gorm:"default:NULL"`
}

// Column :
type Column struct {
	gorm.Model
	UUID   string `gorm:"index:columnuuid"`
	UserID uint   `gorm:"index:usercolumn"`
	Name   string
}

// Migrate :
func (ts *TaskService) Migrate() error {
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
