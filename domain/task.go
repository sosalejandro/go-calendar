package domain

import (
	"context"
	"errors"
	"time"

	domain_errors "github.com/sosalejandro/go-calendar/domain/errors"
)

// Task represents a calendar task
type Task struct {
	id                *TaskID
	title             string
	repeating         bool
	repeatingInterval time.Duration
	description       string
	completed         bool
	dayOfWeek         time.Weekday
	time              time.Time
}

// GetID returns the task ID
func (t *Task) GetID() TaskID {
	return *t.id
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
func (t *Task) IsRepeating() (bool, time.Duration) {
	if t.repeating {
		return t.repeating, t.repeatingInterval
	}

	return false, 0
}

// NewTask creates a new task
func NewTask(
	taskId *TaskID,
	title, description string,
	repeating bool,
	repeatingInterval time.Duration,
	dayOfWeek time.Weekday,
	time time.Time) (*Task, error) {
	task := &Task{
		id:                taskId,
		title:             title,
		repeating:         repeating,
		repeatingInterval: repeatingInterval,
		description:       description,
		dayOfWeek:         dayOfWeek,
		time:              time,
	}

	if err := task.Validate(); err != nil {
		return nil, errors.Join(domain_errors.ErrInvalidTask, err)
	}

	return task, nil
}

// Validate validates the task
func (t *Task) Validate() (errc error) {
	isOriginal := t.id.original && t.id.secondaryId == t.id.primaryId
	isCopy := !t.id.original && t.id.secondaryId != t.id.primaryId

	if !isCopy && !isOriginal {
		errc = errors.Join(domain_errors.ErrInvalidTaskID, errc)
	}

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

// CalculateNextRepeatingTime calculates the next repeating time
func (t *Task) CalculateNextRepeatingTime() time.Time {
	return t.time.Add(t.repeatingInterval)
}

// searchRepetition returns a channel with the next repeating times
//
// If the task is not repeating, the channel will contain only the task time
// If the task is repeating, the channel will contain the next repeating times
// The channel will be closed when there are no more repeating times
func (t *Task) searchRepetition(ctx context.Context) <-chan time.Time {
	// var times []time.Time
	// quit := make(chan struct{})
	isRepeating, interval := t.IsRepeating()

	timeChan := make(chan time.Time)

	if (!isRepeating) || (interval == 0) {
		go func(ch chan<- time.Time) {
			defer close(ch)
			ch <- t.time
		}(timeChan)
		return timeChan
	}

	// next month in time
	lastDay := time.Date(t.time.Year(), t.time.Month()+1, 0, 23, 59, 59, 0, t.time.Location())

	go func() {
		defer close(timeChan)
		currentTime := t.time
		for {
			select {
			case <-ctx.Done():
				// Context was cancelled, exit the goroutine
				return
			default:
				if currentTime.After(lastDay) {
					return
				}

				timeChan <- currentTime
				currentTime = currentTime.Add(interval)
			}
		}
	}()

	// for {
	// 	select {
	// 	case time := <-timeChan:
	// 		times = append(times, time)
	// 	case <-quit:
	// 		close(timeChan)
	// 		close(quit)
	// 		return times
	// 	}
	// }

	return timeChan
}

func (t *Task) createTasksFromDates(ctx context.Context,
	datesChan <-chan time.Time,
	taskIDFactory TaskIDFactory) <-chan *Task {
	taskChan := make(chan *Task)

	go func() {
		defer close(taskChan)

		for date := range datesChan {
			select {
			case <-ctx.Done():
				// Context was cancelled, exit the goroutine
				return
			default:
				taskID := taskIDFactory.CreateTaskID(t.id.primaryId.String())

				task, err := NewTask(
					taskID,
					t.title, t.description,
					t.repeating, t.repeatingInterval,
					t.dayOfWeek, date)
				if err != nil {
					// TODO: handle error
					continue
				}

				taskChan <- task
			}
		}
	}()

	return taskChan
}
