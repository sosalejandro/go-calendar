package domain

import (
	"testing"
	"time"

	domain_errors "github.com/sosalejandro/go-calendar/domain/errors"
	"github.com/stretchr/testify/assert"
)

func TestTask_Validate(t *testing.T) {
	// Define test cases
	testCases := []struct {
		name  string
		task  *Task
		error error
	}{
		{
			name: "Empty title",
			task: &Task{
				title:       "",
				description: "description",
				dayOfWeek:   1,
				time:        time.Now(),
			},
			error: domain_errors.ErrTitleRequired,
		},
		{
			name: "Empty description",
			task: &Task{
				title:       "title",
				description: "",
				dayOfWeek:   1,
				time:        time.Now(),
			},
			error: domain_errors.ErrDescriptionRequired,
		},
		{
			name: "Zero dayOfWeek",
			task: &Task{
				title:       "title",
				description: "description",
				dayOfWeek:   0,
				time:        time.Now(),
			},
			error: domain_errors.ErrDayOfWeekRequired,
		},
		{
			name: "Zero time",
			task: &Task{
				title:       "title",
				description: "description",
				dayOfWeek:   1,
				time:        time.Time{},
			},
			error: domain_errors.ErrTimeRequired,
		},
		{
			name: "Valid task",
			task: &Task{
				title:       "title",
				description: "description",
				dayOfWeek:   1,
				time:        time.Now(),
			},
			error: nil,
		},
	}

	// Run test cases
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := tc.task.Validate()
			assert.ErrorIs(t, err, tc.error)
		})
	}
}

func TestNewTask(t *testing.T) {
	// Define test cases
	testCases := []struct {
		name        string
		title       string
		description string
		repeating   bool
		dayOfWeek   time.Weekday
		time        time.Time
		error       error
	}{
		{
			name:        "Empty title",
			title:       "",
			description: "description",
			repeating:   false,
			dayOfWeek:   1,
			time:        time.Now(),
			error:       domain_errors.ErrTitleRequired,
		},
		{
			name:        "Empty description",
			title:       "title",
			description: "",
			repeating:   false,
			dayOfWeek:   1,
			time:        time.Now(),
			error:       domain_errors.ErrDescriptionRequired,
		},
		{
			name:        "Zero dayOfWeek",
			title:       "title",
			description: "description",
			repeating:   false,
			dayOfWeek:   0,
			time:        time.Now(),
			error:       domain_errors.ErrDayOfWeekRequired,
		},
		{
			name:        "Zero time",
			title:       "title",
			description: "description",
			repeating:   false,
			dayOfWeek:   1,
			time:        time.Time{},
			error:       domain_errors.ErrTimeRequired,
		},
		{
			name:        "Valid task",
			title:       "title",
			description: "description",
			repeating:   false,
			dayOfWeek:   1,
			time:        time.Now(),
			error:       nil,
		},
	}

	// Run test cases
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			task, err := NewTask(tc.title, tc.description, tc.repeating, tc.dayOfWeek, tc.time)
			assert.ErrorIs(t, err, tc.error)

			if err != nil {
				assert.ErrorIs(t, err, domain_errors.ErrInvalidTask)
			}

			if err == nil {
				assert.Equal(t, tc.title, task.GetTitle())
				assert.Equal(t, tc.description, task.GetDescription())
				assert.Equal(t, tc.repeating, task.IsRepeating())
				assert.Equal(t, tc.dayOfWeek, task.GetDayOfWeek())
				assert.Equal(t, tc.time, task.GetTime())
			}
		})
	}
}

func TestTask_Complete(t *testing.T) {
	// Define test cases
	testCases := []struct {
		name     string
		task     *Task
		expected bool
	}{
		{
			name: "Task not completed",
			task: &Task{
				title:       "title",
				description: "description",
				completed:   false,
			},
			expected: true,
		},
		{
			name: "Task already completed",
			task: &Task{
				title:       "title",
				description: "description",
				completed:   true,
			},
			expected: true,
		},
	}

	// Run test cases
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tc.task.Complete()
			assert.Equal(t, tc.expected, tc.task.IsCompleted())
		})
	}
}
