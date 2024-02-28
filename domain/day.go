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
	d := &Day{day: day, tasks: make([]*Task, 0)}

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
		if day.tasks[mid].GetTime().Equal(key) {
			return mid // Found exact match
		} else if day.tasks[mid].GetTime().Before(key) {
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
	position := binarySearch(d, task.GetTime())

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
		for j > 0 && d.tasks[j-1].GetTime().After(d.tasks[j].GetTime()) {
			d.tasks[j], d.tasks[j-1] = d.tasks[j-1], d.tasks[j]
			j--
		}
	}
}

// getDay returns the day and the day of the week
//
// If there are no tasks, the day of the week is not returned.
// If there are tasks, the day of the week is returned.
// The day of the week is the day of the week of the first task.
func (d *Day) getDay() (int, time.Weekday) {

	var weekDay time.Weekday

	if len(d.tasks) > 0 {
		weekDay = d.tasks[0].GetDayOfWeek()
	}

	return d.day, weekDay
}

// Validate validates the day
func (d *Day) Validate() (errc error) {
	if d.day == 0 {
		errc = errors.Join(domain_errors.ErrDayRequired, errc)
	}

	if !(d.day >= 1 && d.day <= 31) {
		errc = errors.Join(domain_errors.ErrInvalidDay, errc)
	}

	return errc
}

// deleteTask deletes a task from the day
func (d *Day) deleteTask(position int, task *Task) error {
	if position >= len(d.tasks) || !d.tasks[position].time.Equal(task.time) {
		return domain_errors.ErrTaskNotFound
	}

	// Delete the task
	d.tasks = append(d.tasks[:position], d.tasks[position+1:]...)

	return nil
}

// updateTask updates a task in the day
func (d *Day) updateTask(position int, task *Task) error {
	// if position >= len(d.tasks) || !d.tasks[position].GetTime().Equal(task.time) {
	if position >= len(d.tasks) {
		return domain_errors.ErrTaskNotFound
	}

	// Update the task
	d.tasks[position] = task

	return nil
}

type findTaskFunc func(*Day, string, time.Time) int

type findType string

var (
	findByTitleAndTime findType = "titleAndTime"
)

func (d *Day) findTask(findTaskFunc findTaskFunc,
	title string,
	time time.Time) (int, *Task, error) {
	position := findTaskFunc(d, title, time)
	if position == -1 {
		return 0, nil, domain_errors.ErrTaskNotFound
	}

	return position, d.tasks[position], nil
}

// func (d *Day) getTaskByTitleAndTime(title string, time time.Time) (int, *Task, error) {
// 	position := binarySearchByTitleAndTime(d, title, time)
// 	if position == -1 {
// 		return 0, nil, domain_errors.ErrTaskNotFound
// 	}

// 	return position, d.tasks[position], nil
// }

// binarySearchTask searches for the task in the day's tasks slice.
func binarySearchByTitleAndTime(day *Day, title string, time time.Time) int {
	return binarySearchByCondition(day, time, func(t *Task) bool {
		return t.GetTitle() == title && t.GetTime().Equal(time)
	})
}

// binarySearchByCondition searches for the task in the day's tasks slice.
func binarySearchByCondition(day *Day, time time.Time, condition func(*Task) bool) int {
	low, high := 0, len(day.tasks)-1

	for low <= high {
		mid := (low + high) / 2
		if condition(day.tasks[mid]) {
			return mid // Found exact match
		} else if day.tasks[mid].GetTime().Before(time) {
			low = mid + 1
		} else {
			high = mid - 1
		}
	}

	return -1 // Not found
}
