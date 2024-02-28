package domain

import (
	"testing"
	"time"

	"github.com/google/uuid"
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
						time:      time.Date(2022, 1, 1, 0, 0, 0, 0, time.UTC),
						dayOfWeek: time.Saturday,
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

func TestFindTask(t *testing.T) {
	// Setup
	day := &Day{
		tasks: []*Task{
			{title: "Task1", time: time.Date(2022, 1, 1, 10, 0, 0, 0, time.UTC)},
			{title: "Task2", time: time.Date(2022, 1, 1, 11, 0, 0, 0, time.UTC)},
			{title: "Task3", time: time.Date(2022, 1, 1, 12, 0, 0, 0, time.UTC)},
		},
	}

	// findTaskFunc := func(d *Day, title string, time time.Time) int {
	// 	return binarySearchByTitleAndTime(d, title, time)
	// }

	tests := []struct {
		name      string
		title     string
		time      time.Time
		findFunc  findTaskFunc
		wantPos   int
		wantTitle string
		wantErr   error
	}{
		{
			name:      "binarySearchByTitleAndTime: Existing task",
			title:     "Task2",
			time:      time.Date(2022, 1, 1, 11, 0, 0, 0, time.UTC),
			findFunc:  binarySearchByTitleAndTime,
			wantPos:   1,
			wantTitle: "Task2",
			wantErr:   nil,
		},
		{
			name:     "binarySearchByTitleAndTime: Non-existing task",
			title:    "Task4",
			time:     time.Date(2022, 1, 1, 13, 0, 0, 0, time.UTC),
			findFunc: binarySearchByTitleAndTime,
			wantErr:  domain_errors.ErrTaskNotFound,
		},
		{
			name:     "binarySearchByTitleAndTime: Non-existing task",
			title:    "Task0",
			time:     time.Date(2022, 1, 1, 9, 0, 0, 0, time.UTC),
			findFunc: binarySearchByTitleAndTime,
			wantErr:  domain_errors.ErrTaskNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pos, task, err := day.findTask(tt.findFunc, tt.title, tt.time)
			if err != tt.wantErr {
				t.Fatalf("Expected error %v, got %v", tt.wantErr, err)
			}
			if tt.wantErr == nil {
				if pos != tt.wantPos {
					t.Errorf("Expected position %d, got %d", tt.wantPos, pos)
				}
				if task.title != tt.wantTitle {
					t.Errorf("Expected task title '%s', got '%s'", tt.wantTitle, task.title)
				}
			}
		})
	}
}

func TestDayUpdateTask(t *testing.T) {
	originalTask := uuid.New()
	originalTime := time.Now()
	task1 := &Task{id: &TaskID{
		primaryId:   originalTask,
		secondaryId: originalTask,
	}, time: originalTime}
	task2 := &Task{id: &TaskID{
		primaryId:   uuid.New(),
		secondaryId: originalTask,
	}, time: originalTime.Add(24 * time.Hour)}
	task3 := &Task{id: &TaskID{
		primaryId:   uuid.New(),
		secondaryId: originalTask,
	}, time: originalTime.Add(48 * time.Hour)}

	tests := []struct {
		name     string
		day      *Day
		position int
		newTask  *Task
		wantErr  error
	}{
		{
			name:     "Update existing task",
			day:      &Day{tasks: []*Task{task1, task2}},
			position: 1,
			newTask:  task3,
			wantErr:  nil,
		},
		{
			name:     "Update non-existing task",
			day:      &Day{tasks: []*Task{task1, task2}},
			position: 2,
			newTask:  task3,
			wantErr:  domain_errors.ErrTaskNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.day.updateTask(tt.position, tt.newTask)
			assert.ErrorIs(t, err, tt.wantErr)
			if err == nil {
				assert.Equal(t, tt.newTask, tt.day.tasks[tt.position])
			}
		})
	}
}

func TestDayDeleteTask(t *testing.T) {
	originalTask := uuid.New()
	originalTime := time.Now()
	task1 := &Task{id: &TaskID{
		primaryId:   originalTask,
		secondaryId: originalTask,
	}, time: originalTime}
	task2 := &Task{id: &TaskID{
		primaryId:   uuid.New(),
		secondaryId: originalTask,
	}, time: originalTime.Add(24 * time.Hour)}

	tests := []struct {
		name     string
		day      *Day
		position int
		task     *Task
		wantErr  error
	}{
		{
			name:     "Delete existing task",
			day:      &Day{tasks: []*Task{task1, task2}},
			position: 1,
			task:     task2,
			wantErr:  nil,
		},
		{
			name:     "Delete non-existing task",
			day:      &Day{tasks: []*Task{task1, task2}},
			position: 2,
			task:     task2,
			wantErr:  domain_errors.ErrTaskNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.day.deleteTask(tt.position, tt.task)
			assert.ErrorIs(t, err, tt.wantErr)
			if err == nil {
				assert.NotContains(t, tt.day.tasks, tt.task)
			}
		})
	}
}
