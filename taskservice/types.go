package service

// RecurringType :
type RecurringType string

const (
	// RecurringEveryDay :
	RecurringEveryDay = RecurringType("day")

	// RecurringEveryWeek :
	RecurringEveryWeek = RecurringType("week")

	// RecurringEveryMonth :
	RecurringEveryMonth = RecurringType("month")

	// RecurringEveryYear :
	RecurringEveryYear = RecurringType("year")

	// RecurringNone :
	RecurringNone = RecurringType("none")
)
