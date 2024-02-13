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
)

var (
	// ErrTaskCannotBeNil is returned when a task is nil
	ErrTaskCannotBeNil = errors.New("task is nil")
)
