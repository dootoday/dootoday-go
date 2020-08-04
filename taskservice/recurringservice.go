package service

import (
	"strings"
	"time"
)

// RecurringTaskService :
type RecurringTaskService struct {
	TDS ITaskDBService
}

// NewRecurringTaskService  :
func NewRecurringTaskService(tds ITaskDBService) *RecurringTaskService {
	return &RecurringTaskService{
		TDS: tds,
	}
}

// IsRecurringTask :
// Returns Actual Task and Recurring Task Type
func (ts *RecurringTaskService) IsRecurringTask(task string) (string, RecurringType) {
	splitby := "every"
	splitted := strings.Split(task, splitby)
	if len(splitted) > 1 {
		last := splitted[len(splitted)-1]
		trimmed := strings.TrimSpace(last)
		actualTask := strings.TrimSpace(strings.Join(splitted[:len(splitted)-1], splitby))
		if actualTask == "" {
			return task, RecurringNone
		}
		if string(trimmed) == string(RecurringEveryDay) {
			return actualTask, RecurringEveryDay
		}
		if string(trimmed) == string(RecurringEveryMonth) {
			return actualTask, RecurringEveryMonth
		}
		if string(trimmed) == string(RecurringEveryWeek) {
			return actualTask, RecurringEveryWeek
		}
		if string(trimmed) == string(RecurringEveryYear) {
			return actualTask, RecurringEveryYear
		}
	}
	return task, RecurringNone
}

// DoesMatchRecurring :
func (ts *RecurringTaskService) DoesMatchRecurring(
	startTime time.Time,
	checkTime time.Time,
	recurringType RecurringType,
) bool {
	// If the recurring task created before current date then
	// current date should not match with the task
	if checkTime.Before(startTime) {
		return false
	}
	// If the recurring type is every day then any current date
	// should match with the task except the condition above
	if recurringType == RecurringEveryDay {
		return true
	}
	// if the recurring type is every week then any day with same
	// week dayshould match
	if recurringType == RecurringEveryWeek {
		if checkTime.Weekday() == startTime.Weekday() {
			return true
		}
	}
	// if the recurring type is every month then any day with same
	// day day should match
	// If it's the last day of the month
	if recurringType == RecurringEveryMonth {
		nextDayOfStartTime := startTime.AddDate(0, 0, 1)
		isStartLastDay := nextDayOfStartTime.Day() == 1
		nextDayOfCheckTime := checkTime.AddDate(0, 0, 1)
		isCheckLastDay := nextDayOfCheckTime.Day() == 1
		if isCheckLastDay && isStartLastDay {
			return true
		}
		if checkTime.Day() == startTime.Day() {
			return true
		}
	}

	// if the recurring type is every week then any day with same
	// week dayshould match
	if recurringType == RecurringEveryYear {
		if checkTime.Month() == startTime.Month() {
			if startTime.Month().String() == "February" &&
				startTime.Day() == 29 {
				nextDayOfCheckTime := checkTime.AddDate(0, 0, 1)
				return nextDayOfCheckTime.Day() == 1
			}
			if checkTime.Day() == startTime.Day() {
				return true
			}
		}
	}
	return false
}
