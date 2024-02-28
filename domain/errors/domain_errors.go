package domain_errors

import "errors"

var (
	// ErrTitleRequired is returned when a task title is required
	ErrTitleRequired = errors.New("title is required")
	// ErrTimeRequired is returned when a task time is required
	ErrTimeRequired = errors.New("time is required")
	// ErrDescriptionRequired is returned when a task description is required
	ErrDescriptionRequired = errors.New("description is required")
	// ErrDayOfWeekRequired is returned when a task day of the week is required
	ErrDayOfWeekRequired = errors.New("day of the week is required")
	// ErrDayRequired is returned when a day is required
	ErrDayRequired = errors.New("day is required")
)

var (
	// ErrInvalidTask is returned when a task is invalid
	ErrInvalidTask = errors.New("invalid task")
	// ErrInvalidDay is returned when a day is invalid
	ErrInvalidDay = errors.New("invalid day")
	// ErrInvalidMonth is returned when a month is invalid
	ErrInvalidMonth = errors.New("invalid month")
	// ErrInvalidYear is returned when a year is invalid
	ErrInvalidYear = errors.New("invalid year")
	// ErrInvalidTaskID is returned when a task ID is invalid
	ErrInvalidTaskID = errors.New("invalid task ID")
)

var (
	// ErrTaskCannotBeNil is returned when a task is nil
	ErrTaskCannotBeNil = errors.New("task is nil")
)

var (
	// ErrDayNotFound is returned when a day is not found
	ErrDayNotFound = errors.New("day not found")
	// ErrTaskNotFound is returned when a task is not found
	ErrTaskNotFound = errors.New("task not found")
)

var (
	// ErrAddDay is returned when a day cannot be added
	ErrAddDay = errors.New("cannot add day")
	// ErrAddMonth is returned when a month cannot be added
	ErrAddMonth = errors.New("cannot add month")
)
