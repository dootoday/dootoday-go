package service

import (
	"time"

	emailservice "apidootoday/emailservice"
	taskservice "apidootoday/taskservice"
	userservice "apidootoday/user"
)

// CronService :
type CronService struct {
	us *userservice.UserService
	ts *taskservice.TaskService
	es *emailservice.EmailService
}

// NewCronService :
func NewCronService(
	us *userservice.UserService,
	ts *taskservice.TaskService,
	es *emailservice.EmailService,
) *CronService {
	return &CronService{
		us: us,
		ts: ts,
		es: es,
	}
}

// CalculateOffset :
// Time in minutes
func (cs *CronService) CalculateOffset(mins int, utcNowInMin int) int {
	return (-1) * (mins - utcNowInMin)
}

// MoveTasksToTodayCron :
func (cs *CronService) MoveTasksToTodayCron() error {
	utcNow := time.Now().UTC()
	utcNowInMins := (utcNow.Hour() * 60) + utcNow.Minute()
	offset := 0
	if utcNowInMins < 720 {
		offset = utcNowInMins
	} else {
		offset = utcNowInMins - 1440
	}
	users, err := cs.us.GetUsersByTimeZoneOffset(offset)
	if err != nil {
		return err
	}
	newDateForTasks := utcNow
	if offset < 0 {
		// If offset is negative then the place is
		// ahead of UTC so the date there is one day later
		newDateForTasks = utcNow.Add(time.Hour * 24)
	}
	for _, user := range users {
		if user.AllowAutoUpdate {
			tasks, err := cs.ts.UpdateNonRecurringTaskDatesByUserID(user.ID, newDateForTasks)
			if err != nil {
				return err
			}

			// Send email here to notify the user that the tasks are moved
			if len(tasks) > 0 && user.AllowDailyEmailUpdate {
				cs.es.SendTaskMoveEmail(user.Email, user.FirstName+" "+user.LastName, user.FirstName, tasks)
			}
		}
	}
	return nil
}

// DailyMorningEmailCron :
// For 07:00 hours
func (cs *CronService) DailyMorningEmailCron() error {
	utcNow := time.Now().UTC()
	utcNowInMins := (utcNow.Hour() * 60) + utcNow.Minute()

	// Offset for 7 am in the morning
	offset := cs.CalculateOffset(7*60, utcNowInMins)

	users, err := cs.us.GetUsersByTimeZoneOffset(offset)
	if err != nil {
		return err
	}

	for _, user := range users {
		if user.ID == 2 && user.Email == "sanborn.sen@gmail.com" && user.AllowDailyEmailUpdate {
			// Send email here to notify the user that the tasks are moved
			cs.es.SendWelcomeEmail(user.Email, user.FirstName+" "+user.LastName, user.FirstName)
		}
	}
	return nil
}
