package service

import "time"

// ITaskDBService :
type ITaskDBService interface {
	FindOrCreateRecurringTaskStatus(taskID uint, date *time.Time) (RecurringTaskStatus, error)
	UpdateRecurringTaskStatus(rts RecurringTaskStatus) error
	GetRecurringTaskStatusByID(recurringID uint) (RecurringTaskStatus, error)
}
