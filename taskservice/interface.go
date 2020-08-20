package service

// ITaskDBService :
type ITaskDBService interface {
	FindOrCreateRecurringTaskStatus(taskID uint, date string) (RecurringTaskStatus, error)
	UpdateRecurringTaskStatus(rts RecurringTaskStatus) error
	GetRecurringTaskStatusByID(recurringID uint) (RecurringTaskStatus, error)
}
