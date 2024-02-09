package domain

import (
	"errors"
	"time"

	domain_errors "github.com/sosalejandro/go-calendar/domain/errors"
)

// Task represents a calendar task
type Task struct {
	title       string
	repeating   bool
	description string
	completed   bool
	dayOfWeek   time.Weekday
	time        time.Time
}

// GetTitle returns the task title
func (t *Task) GetTitle() string {
	return t.title
}

// GetDescription returns the task description
func (t *Task) GetDescription() string {
	return t.description
}

// IsCompleted returns true if the task is completed
func (t *Task) IsCompleted() bool {
	return t.completed
}

// GetDayOfWeek returns the day of the week of the task
func (t *Task) GetDayOfWeek() time.Weekday {
	return t.dayOfWeek
}

// GetTime returns the time of the task
func (t *Task) GetTime() time.Time {
	return t.time
}

// IsRepeating returns true if the task is repeating
func (t *Task) IsRepeating() bool {
	return t.repeating
}

// NewTask creates a new task
func NewTask(title, description string, repeating bool, dayOfWeek time.Weekday, time time.Time) (*Task, error) {
	task := &Task{
		title:       title,
		repeating:   repeating,
		description: description,
		dayOfWeek:   dayOfWeek,
		time:        time,
	}

	if err := task.Validate(); err != nil {
		return nil, errors.Join(domain_errors.ErrInvalidTask, err)
	}

	return task, nil
}

// Validate validates the task
func (t *Task) Validate() (errc error) {
	if t.title == "" {
		errc = errors.Join(domain_errors.ErrTitleRequired, errc)
	}

	if t.description == "" {
		errc = errors.Join(domain_errors.ErrDescriptionRequired, errc)

	}

	if t.dayOfWeek == 0 {
		errc = errors.Join(domain_errors.ErrDayOfWeekRequired, errc)
	}

	if t.time.IsZero() {
		errc = errors.Join(domain_errors.ErrTimeRequired, errc)
	}

	return errc
}

// Complete marks the task as completed
func (t *Task) Complete() {
	t.completed = true
}
