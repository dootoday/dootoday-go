package service

import (
	"testing"

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
