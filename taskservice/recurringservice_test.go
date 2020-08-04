package service

import (
	"testing"
	"time"

	"github.com/golang/mock/gomock"
)

func TestIsRecurringTask(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	tests := map[string]struct {
		service    func() *RecurringTaskService
		input      string
		outputTask string
		outputType RecurringType
	}{
		"all ok test": {
			service: func() *RecurringTaskService {
				mockTaskDBService := NewMockITaskDBService(ctrl)
				rts := NewRecurringTaskService(mockTaskDBService)
				return rts
			},
			input:      "This is a task",
			outputTask: "This is a task",
			outputType: RecurringNone,
		},

		"every day": {
			service: func() *RecurringTaskService {
				mockTaskDBService := NewMockITaskDBService(ctrl)
				rts := NewRecurringTaskService(mockTaskDBService)
				return rts
			},
			input:      "This is a task every day",
			outputTask: "This is a task",
			outputType: RecurringEveryDay,
		},

		"every day with more": {
			service: func() *RecurringTaskService {
				mockTaskDBService := NewMockITaskDBService(ctrl)
				rts := NewRecurringTaskService(mockTaskDBService)
				return rts
			},
			input:      "This is a task every day once",
			outputTask: "This is a task every day once",
			outputType: RecurringNone,
		},

		"just every": {
			service: func() *RecurringTaskService {
				mockTaskDBService := NewMockITaskDBService(ctrl)
				rts := NewRecurringTaskService(mockTaskDBService)
				return rts
			},
			input:      "This is a task every",
			outputTask: "This is a task every",
			outputType: RecurringNone,
		},

		"every week": {
			service: func() *RecurringTaskService {
				mockTaskDBService := NewMockITaskDBService(ctrl)
				rts := NewRecurringTaskService(mockTaskDBService)
				return rts
			},
			input:      "This is a task every week",
			outputTask: "This is a task",
			outputType: RecurringEveryWeek,
		},

		"every month": {
			service: func() *RecurringTaskService {
				mockTaskDBService := NewMockITaskDBService(ctrl)
				rts := NewRecurringTaskService(mockTaskDBService)
				return rts
			},
			input:      "This is a task every month",
			outputTask: "This is a task",
			outputType: RecurringEveryMonth,
		},

		"every month together": {
			service: func() *RecurringTaskService {
				mockTaskDBService := NewMockITaskDBService(ctrl)
				rts := NewRecurringTaskService(mockTaskDBService)
				return rts
			},
			input:      "This is a task everymonth",
			outputTask: "This is a task",
			outputType: RecurringEveryMonth,
		},

		"every year": {
			service: func() *RecurringTaskService {
				mockTaskDBService := NewMockITaskDBService(ctrl)
				rts := NewRecurringTaskService(mockTaskDBService)
				return rts
			},
			input:      "This is a task every year",
			outputTask: "This is a task",
			outputType: RecurringEveryYear,
		},

		"every year together": {
			service: func() *RecurringTaskService {
				mockTaskDBService := NewMockITaskDBService(ctrl)
				rts := NewRecurringTaskService(mockTaskDBService)
				return rts
			},
			input:      "This is a task everyyear",
			outputTask: "This is a task",
			outputType: RecurringEveryYear,
		},

		"empty test": {
			service: func() *RecurringTaskService {
				mockTaskDBService := NewMockITaskDBService(ctrl)
				rts := NewRecurringTaskService(mockTaskDBService)
				return rts
			},
			input:      "",
			outputTask: "",
			outputType: RecurringNone,
		},

		"just keyword": {
			service: func() *RecurringTaskService {
				mockTaskDBService := NewMockITaskDBService(ctrl)
				rts := NewRecurringTaskService(mockTaskDBService)
				return rts
			},
			input:      "every day",
			outputTask: "every day",
			outputType: RecurringNone,
		},

		"with special char": {
			service: func() *RecurringTaskService {
				mockTaskDBService := NewMockITaskDBService(ctrl)
				rts := NewRecurringTaskService(mockTaskDBService)
				return rts
			},
			input:      "hello every month`",
			outputTask: "hello every month`",
			outputType: RecurringNone,
		},

		"multiple every - 1": {
			service: func() *RecurringTaskService {
				mockTaskDBService := NewMockITaskDBService(ctrl)
				rts := NewRecurringTaskService(mockTaskDBService)
				return rts
			},
			input:      "every day remember to rock every day",
			outputTask: "every day remember to rock",
			outputType: RecurringEveryDay,
		},

		"multiple every - 2": {
			service: func() *RecurringTaskService {
				mockTaskDBService := NewMockITaskDBService(ctrl)
				rts := NewRecurringTaskService(mockTaskDBService)
				return rts
			},
			input:      "every day remember every day every month",
			outputTask: "every day remember every day",
			outputType: RecurringEveryMonth,
		},
	}
	for name, param := range tests {
		t.Run(name, func(t *testing.T) {
			s := param.service()
			optask, optype := s.IsRecurringTask(param.input)
			if optask != param.outputTask {
				t.Errorf(
					"Expected task '%s' but got '%s'",
					param.outputTask, optask,
				)
			}
			if optype != param.outputType {
				t.Errorf(
					"Expected type '%s' but got '%s'",
					param.outputType, optype,
				)
			}
		})
	}
}

func TestDoesMatchRecurring(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	timeToday, _ := time.Parse("2006-01-02", time.Now().Format("2006-01-02"))
	timeYesterday, _ := time.Parse("2006-01-02", time.Now().AddDate(0, 0, -1).Format("2006-01-02"))
	timeTomorrow, _ := time.Parse("2006-01-02", time.Now().AddDate(0, 0, 1).Format("2006-01-02"))
	timeNextWeek, _ := time.Parse("2006-01-02", time.Now().AddDate(0, 0, 7).Format("2006-01-02"))
	timeLastWeek, _ := time.Parse("2006-01-02", time.Now().AddDate(0, 0, -7).Format("2006-01-02"))

	// Month check
	thirdAugust, _ := time.Parse("2006-01-02", "2020-08-03")
	thirdSeptember, _ := time.Parse("2006-01-02", "2020-09-03")
	fourthSeptember, _ := time.Parse("2006-01-02", "2020-09-04")
	thirtyFirstJan, _ := time.Parse("2006-01-02", "2020-01-31")
	twentyEighthFeb, _ := time.Parse("2006-01-02", "2020-02-28")
	twentyNinthFeb, _ := time.Parse("2006-01-02", "2020-02-29")
	thirtyApril, _ := time.Parse("2006-01-02", "2020-04-30")

	// Year check
	thirtyFirstJan2021, _ := time.Parse("2006-01-02", "2021-01-31")
	twentyEighthFeb2021, _ := time.Parse("2006-01-02", "2021-02-28")
	thirtyFirstMarch2021, _ := time.Parse("2006-01-02", "2021-03-31")
	thirtyApril2021, _ := time.Parse("2006-01-02", "2021-04-30")

	tests := map[string]struct {
		service            func() *RecurringTaskService
		inputStratTime     time.Time
		inputCheckTime     time.Time
		inputRecurringType RecurringType
		output             bool
	}{
		"same day as recurring": {
			service: func() *RecurringTaskService {
				mockTaskDBService := NewMockITaskDBService(ctrl)
				rts := NewRecurringTaskService(mockTaskDBService)
				return rts
			},
			inputStratTime:     timeToday,
			inputCheckTime:     timeToday,
			inputRecurringType: RecurringEveryDay,
			output:             true,
		},

		"prev day to recurring": {
			service: func() *RecurringTaskService {
				mockTaskDBService := NewMockITaskDBService(ctrl)
				rts := NewRecurringTaskService(mockTaskDBService)
				return rts
			},
			inputStratTime:     timeToday,
			inputCheckTime:     timeYesterday,
			inputRecurringType: RecurringEveryDay,
			output:             false,
		},

		"future day to recurring for everyday": {
			service: func() *RecurringTaskService {
				mockTaskDBService := NewMockITaskDBService(ctrl)
				rts := NewRecurringTaskService(mockTaskDBService)
				return rts
			},
			inputStratTime:     timeToday,
			inputCheckTime:     timeTomorrow,
			inputRecurringType: RecurringEveryDay,
			output:             true,
		},

		"future day to recurring for everyweek match": {
			service: func() *RecurringTaskService {
				mockTaskDBService := NewMockITaskDBService(ctrl)
				rts := NewRecurringTaskService(mockTaskDBService)
				return rts
			},
			inputStratTime:     timeToday,
			inputCheckTime:     timeNextWeek,
			inputRecurringType: RecurringEveryWeek,
			output:             true,
		},

		"future day to recurring for everyweek no match": {
			service: func() *RecurringTaskService {
				mockTaskDBService := NewMockITaskDBService(ctrl)
				rts := NewRecurringTaskService(mockTaskDBService)
				return rts
			},
			inputStratTime:     timeToday,
			inputCheckTime:     timeTomorrow,
			inputRecurringType: RecurringEveryWeek,
			output:             false,
		},

		"prev week to recurring for everyweek no match": {
			service: func() *RecurringTaskService {
				mockTaskDBService := NewMockITaskDBService(ctrl)
				rts := NewRecurringTaskService(mockTaskDBService)
				return rts
			},
			inputStratTime:     timeToday,
			inputCheckTime:     timeLastWeek,
			inputRecurringType: RecurringEveryWeek,
			output:             false,
		},

		"same month for everymonth match": {
			service: func() *RecurringTaskService {
				mockTaskDBService := NewMockITaskDBService(ctrl)
				rts := NewRecurringTaskService(mockTaskDBService)
				return rts
			},
			inputStratTime:     thirdAugust,
			inputCheckTime:     thirdAugust,
			inputRecurringType: RecurringEveryMonth,
			output:             true,
		},

		"next month for everymonth match": {
			service: func() *RecurringTaskService {
				mockTaskDBService := NewMockITaskDBService(ctrl)
				rts := NewRecurringTaskService(mockTaskDBService)
				return rts
			},
			inputStratTime:     thirdAugust,
			inputCheckTime:     thirdSeptember,
			inputRecurringType: RecurringEveryMonth,
			output:             true,
		},

		"next month for everymonth no match": {
			service: func() *RecurringTaskService {
				mockTaskDBService := NewMockITaskDBService(ctrl)
				rts := NewRecurringTaskService(mockTaskDBService)
				return rts
			},
			inputStratTime:     thirdAugust,
			inputCheckTime:     fourthSeptember,
			inputRecurringType: RecurringEveryMonth,
			output:             false,
		},

		"next month for everymonth lastday match - 1": {
			service: func() *RecurringTaskService {
				mockTaskDBService := NewMockITaskDBService(ctrl)
				rts := NewRecurringTaskService(mockTaskDBService)
				return rts
			},
			inputStratTime:     thirtyFirstJan,
			inputCheckTime:     twentyNinthFeb,
			inputRecurringType: RecurringEveryMonth,
			output:             true,
		},

		"next month for everymonth lastday no match - 2": {
			service: func() *RecurringTaskService {
				mockTaskDBService := NewMockITaskDBService(ctrl)
				rts := NewRecurringTaskService(mockTaskDBService)
				return rts
			},
			inputStratTime:     thirtyFirstJan,
			inputCheckTime:     twentyEighthFeb,
			inputRecurringType: RecurringEveryMonth,
			output:             false,
		},

		"next month for everymonth lastday match - 3": {
			service: func() *RecurringTaskService {
				mockTaskDBService := NewMockITaskDBService(ctrl)
				rts := NewRecurringTaskService(mockTaskDBService)
				return rts
			},
			inputStratTime:     thirtyFirstJan,
			inputCheckTime:     thirtyApril,
			inputRecurringType: RecurringEveryMonth,
			output:             true,
		},

		"next month for everymonth lastday match - 4": {
			service: func() *RecurringTaskService {
				mockTaskDBService := NewMockITaskDBService(ctrl)
				rts := NewRecurringTaskService(mockTaskDBService)
				return rts
			},
			inputStratTime:     twentyNinthFeb,
			inputCheckTime:     twentyEighthFeb2021,
			inputRecurringType: RecurringEveryMonth,
			output:             true,
		},

		"prev month for everymonth lastday no match": {
			service: func() *RecurringTaskService {
				mockTaskDBService := NewMockITaskDBService(ctrl)
				rts := NewRecurringTaskService(mockTaskDBService)
				return rts
			},
			inputStratTime:     thirtyApril,
			inputCheckTime:     thirtyFirstJan,
			inputRecurringType: RecurringEveryMonth,
			output:             false,
		},

		"year for every year match - 1": {
			service: func() *RecurringTaskService {
				mockTaskDBService := NewMockITaskDBService(ctrl)
				rts := NewRecurringTaskService(mockTaskDBService)
				return rts
			},
			inputStratTime:     thirtyFirstJan,
			inputCheckTime:     thirtyFirstJan2021,
			inputRecurringType: RecurringEveryYear,
			output:             true,
		},

		"year for every year match - 2": {
			service: func() *RecurringTaskService {
				mockTaskDBService := NewMockITaskDBService(ctrl)
				rts := NewRecurringTaskService(mockTaskDBService)
				return rts
			},
			inputStratTime:     twentyNinthFeb,
			inputCheckTime:     twentyEighthFeb2021,
			inputRecurringType: RecurringEveryYear,
			output:             true,
		},

		"year for every year match - 3": {
			service: func() *RecurringTaskService {
				mockTaskDBService := NewMockITaskDBService(ctrl)
				rts := NewRecurringTaskService(mockTaskDBService)
				return rts
			},
			inputStratTime:     thirtyApril,
			inputCheckTime:     thirtyApril2021,
			inputRecurringType: RecurringEveryYear,
			output:             true,
		},

		"year for every year no match - 1": {
			service: func() *RecurringTaskService {
				mockTaskDBService := NewMockITaskDBService(ctrl)
				rts := NewRecurringTaskService(mockTaskDBService)
				return rts
			},
			inputStratTime:     thirtyFirstJan,
			inputCheckTime:     thirtyFirstMarch2021,
			inputRecurringType: RecurringEveryYear,
			output:             false,
		},

		"year for every year no match - 2": {
			service: func() *RecurringTaskService {
				mockTaskDBService := NewMockITaskDBService(ctrl)
				rts := NewRecurringTaskService(mockTaskDBService)
				return rts
			},
			inputStratTime:     thirtyFirstJan,
			inputCheckTime:     thirtyApril2021,
			inputRecurringType: RecurringEveryYear,
			output:             false,
		},
	}

	for name, param := range tests {
		t.Run(name, func(t *testing.T) {
			s := param.service()
			output := s.DoesMatchRecurring(
				param.inputStratTime,
				param.inputCheckTime,
				param.inputRecurringType,
			)
			if output != param.output {
				t.Errorf(
					"Expected task '%v' but got '%v'",
					param.output, output,
				)
			}
		})
	}
}
