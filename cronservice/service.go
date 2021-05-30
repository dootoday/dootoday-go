package service

import (
	"fmt"
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

// MoveTasksToTodayCron :
func (cs *CronService) MoveTasksToTodayCron() error {
	utcNow := time.Now().UTC()
	utcNowInMins := (utcNow.Hour() * 60) + utcNow.Minute()
	fmt.Println(utcNow.Hour())
	fmt.Println(utcNowInMins)
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
