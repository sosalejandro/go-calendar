package domain

import (
	"errors"
	"time"

	domain_errors "github.com/sosalejandro/go-calendar/domain/errors"
)

// Month represents a calendar month
type Month struct {
	month time.Month
	year  int
	days  map[int]*Day
}

// Validate validates the month
func (m *Month) Validate() error {
	if m.month < 1 || m.month > 12 {
		return domain_errors.ErrInvalidMonth
	}

	if m.year < 1 {
		return domain_errors.ErrInvalidYear
	}

	return nil
}

// NewMonth creates a new month
func NewMonth(month time.Month, year int) (*Month, error) {
	m := &Month{
		month: month,
		year:  year,
		days:  make(map[int]*Day, 31),
	}

	if err := m.Validate(); err != nil {
		return nil, err
	}

	return m, nil
}

// addDay adds a day to the month
//
// If the day already exists, addDay returns the existing day.
// If the day is invalid, addDay returns domain_errors.ErrAddDay.
// If the day cannot be added, addDay returns domain_errors.ErrAddDay.
// If the day is added successfully, addDay returns the day.
func (m *Month) addDay(day int) (*Day, error) {
	var d *Day

	// Validate the day exists
	d, exists := m.days[day]
	if exists {
		return d, nil
	}

	// Create a new day
	d, err := NewDay(day)
	if err != nil {
		return nil, errors.Join(domain_errors.ErrAddDay, err)
	}

	// int day
	var id int

	id, _ = d.getDay()
	// Add the day to the month's day
	m.days[id] = d

	return d, nil
}

// addTaskToDay adds a task to a specific day of the month
// If the day does not exist, addTaskToDay returns domain_errors.ErrDayNotFound.
func (m *Month) addTaskToDay(day int, task *Task) error {
	d, exists := m.days[day]
	if !exists {
		return domain_errors.ErrDayNotFound
	}

	return d.addTask(task)
}

// getDay returns the day for the month
func (m *Month) getDay(day int) (*Day, error) {
	d, exists := m.days[day]
	if !exists {
		return nil, domain_errors.ErrDayNotFound
	}

	return d, nil
}
