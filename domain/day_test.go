package domain

import (
	"testing"
	"time"

	domain_errors "github.com/sosalejandro/go-calendar/domain/errors"
	"github.com/stretchr/testify/assert"
)

func TestNewDay(t *testing.T) {
	// Define test cases
	testCases := []struct {
		name     string
		day      int
		expected error
	}{
		{
			name:     "Invalid day",
			day:      0,
			expected: domain_errors.ErrDayRequired,
		},
		{
			name:     "Valid day",
			day:      1,
			expected: nil,
		},
		// Add more test cases as needed
	}

	// Run test cases
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			_, err := NewDay(tc.day)
			assert.ErrorIs(t, err, tc.expected)
		})
	}
}

func TestDay_addTask(t *testing.T) {
	// Define test cases
	testCases := []struct {
		name        string
		day         *Day
		task        *Task
		expected    int
		expectedErr error
	}{
		{
			name: "Add task to day after existing task",
			day: &Day{
				day: 1,
				tasks: []*Task{
					{
						time: time.Date(2022, 1, 1, 0, 0, 0, 0, time.UTC),
					},
				},
			},
			task: &Task{
				time: time.Date(2022, 1, 1, 1, 0, 0, 0, time.UTC),
			},
			expected: 2,
		},
		{
			name: "Add task to day before existing task",
			day: &Day{
				day: 1,
				tasks: []*Task{
					{
						time: time.Date(2022, 1, 1, 2, 0, 0, 0, time.UTC),
					},
				},
			},
			task: &Task{
				time: time.Date(2022, 1, 1, 1, 0, 0, 0, time.UTC),
			},
			expected: 2,
		},
		{
			name: "Add task to day with exact same time as existing task",
			day: &Day{
				day: 1,
				tasks: []*Task{
					{
						time: time.Date(2022, 1, 1, 0, 1, 1, 1, time.UTC),
					},
				},
			},
			task: &Task{
				time: time.Date(2022, 1, 1, 0, 1, 1, 1, time.UTC),
			},
			expected: 2,
		},
		{
			name: "Add invalid task to day",
			day: &Day{
				day:   1,
				tasks: []*Task{},
			},
			task:        nil,
			expected:    0,
			expectedErr: domain_errors.ErrTaskCannotBeNil,
		},
		// Add more test cases as needed
	}

	// Run test cases
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := tc.day.addTask(tc.task)
			assert.ErrorIs(t, tc.expectedErr, err)

			assert.Len(t, tc.day.getTasks(), tc.expected)
		})
	}
}

func TestDay_getDay(t *testing.T) {
	// Define test cases
	testCases := []struct {
		name            string
		day             *Day
		expectedDay     int
		expectedWeekday time.Weekday
	}{
		{
			name: "Day with tasks",
			day: &Day{
				day: 1,
				tasks: []*Task{
					{
						time: time.Date(2022, 1, 1, 0, 0, 0, 0, time.UTC),
					},
				},
			},
			expectedDay:     1,
			expectedWeekday: time.Saturday,
		},
		// Add more test cases as needed
	}

	// Run test cases
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			day, weekday := tc.day.getDay()
			assert.Equal(t, tc.expectedDay, day)
			assert.Equal(t, tc.expectedWeekday, weekday)
		})
	}
}

func TestDay_sortTasks(t *testing.T) {
	// Define test cases
	testCases := []struct {
		name     string
		day      *Day
		expected []*Task
	}{
		{
			name: "Sort tasks",
			day: &Day{
				day: 1,
				tasks: []*Task{
					{
						time: time.Date(2022, 1, 1, 12, 0, 0, 0, time.UTC),
					},
					{
						time: time.Date(2022, 1, 1, 10, 0, 0, 0, time.UTC),
					},
				},
			},
			expected: []*Task{
				{
					time: time.Date(2022, 1, 1, 10, 0, 0, 0, time.UTC),
				},
				{
					time: time.Date(2022, 1, 1, 12, 0, 0, 0, time.UTC),
				},
			},
		},
		// Add more test cases as needed
	}

	// Run test cases
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tc.day.sortTasks()
			assert.Equal(t, tc.expected, tc.day.getTasks())
		})
	}
}

func TestDay_Validate(t *testing.T) {
	// Define test cases
	testCases := []struct {
		name     string
		day      *Day
		expected error
	}{
		{
			name: "Invalid day",
			day: &Day{
				day: 0,
			},
			expected: domain_errors.ErrDayRequired,
		},
		{
			name: "Valid day",
			day: &Day{
				day: 1,
			},
			expected: nil,
		},
		// Add more test cases as needed
	}

	// Run test cases
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := tc.day.Validate()
			assert.ErrorIs(t, err, tc.expected)
		})
	}
}
