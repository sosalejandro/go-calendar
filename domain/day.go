package domain

import (
	"errors"
	"time"

	domain_errors "github.com/sosalejandro/go-calendar/domain/errors"
)

// Day represents a day within a month with associated tasks
type Day struct {
	day   int
	tasks []*Task
}

// NewDay creates a new day
func NewDay(day int) (*Day, error) {
	d := &Day{day: day}

	if err := d.Validate(); err != nil {
		return nil, err
	}

	return d, nil

}

// binarySearch searches for the key in the day's tasks slice.
func binarySearch(day *Day, key time.Time) int {
	low, high := 0, len(day.tasks)-1

	for low <= high {
		mid := (low + high) / 2
		if day.tasks[mid].time.Equal(key) {
			return mid // Found exact match
		} else if day.tasks[mid].time.Before(key) {
			low = mid + 1
		} else {
			high = mid - 1
		}
	}

	return low // Position for insertion
}

// addTask adds a task to the day.
//
// Inserts the task into the tasks slice in sorted order by time.
//
// If the task is nil, AddTask returns domain_errors.ErrTaskCannotBeNil.
func (d *Day) addTask(task *Task) error {
	// Validate the task
	if task == nil {
		return domain_errors.ErrTaskCannotBeNil
	}

	// Find the correct position for the task
	position := binarySearch(d, task.time)

	// Insert the task at the correct position
	d.tasks = append(d.tasks[:position], append([]*Task{task}, d.tasks[position:]...)...)

	return nil
}

// getTasks returns the tasks for the day
func (d *Day) getTasks() []*Task {
	return d.tasks
}

// sortTasks sorts the tasks for the day by time
func (d *Day) sortTasks() {
	// Insertion sort
	for i := 1; i < len(d.tasks); i++ {
		j := i
		for j > 0 && d.tasks[j-1].time.After(d.tasks[j].time) {
			d.tasks[j], d.tasks[j-1] = d.tasks[j-1], d.tasks[j]
			j--
		}
	}
}

func (d *Day) getDay() (int, time.Weekday) {

	var weekDay time.Weekday

	if len(d.tasks) > 0 {
		weekDay = d.tasks[0].time.Weekday()
	}

	return d.day, weekDay
}

// Validate validates the day
func (d *Day) Validate() (errc error) {
	if d.day == 0 {
		errc = errors.Join(domain_errors.ErrDayRequired, errc)
	}

	if (d.day >= 1 && d.day <= 31) == false {
		errc = errors.Join(domain_errors.ErrInvalidDay, errc)
	}

	return errc
}
