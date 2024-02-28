package domain

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
	domain_errors "github.com/sosalejandro/go-calendar/domain/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestTask_Validate(t *testing.T) {
	originalTaskIDGUID := uuid.New()
	taskId := &TaskID{
		secondaryId: originalTaskIDGUID,
		primaryId:   originalTaskIDGUID,
		original:    true,
	}

	// Define test cases
	testCases := []struct {
		name  string
		task  *Task
		error error
	}{
		{
			name: "Empty title",
			task: &Task{
				id:          taskId,
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
				id:          taskId,
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
				id:          taskId,
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
				id:          taskId,
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
				id:          taskId,
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
	originalTaskIDGUID := uuid.New()
	taskId := &TaskID{
		secondaryId: originalTaskIDGUID,
		primaryId:   originalTaskIDGUID,
		original:    true,
	}
	// Define test cases
	testCases := []struct {
		id                *TaskID
		name              string
		title             string
		description       string
		repeating         bool
		repeatingInterval time.Duration
		dayOfWeek         time.Weekday
		time              time.Time
		error             error
	}{
		{
			id:                taskId,
			name:              "Empty title",
			title:             "",
			description:       "description",
			repeating:         false,
			repeatingInterval: time.Duration(0),
			dayOfWeek:         1,
			time:              time.Now(),
			error:             domain_errors.ErrTitleRequired,
		},
		{
			id:                taskId,
			name:              "Empty description",
			title:             "title",
			description:       "",
			repeating:         false,
			repeatingInterval: time.Duration(0),
			dayOfWeek:         1,
			time:              time.Now(),
			error:             domain_errors.ErrDescriptionRequired,
		},
		{
			id:                taskId,
			name:              "Zero dayOfWeek",
			title:             "title",
			description:       "description",
			repeating:         false,
			repeatingInterval: time.Duration(0),
			dayOfWeek:         0,
			time:              time.Now(),
			error:             domain_errors.ErrDayOfWeekRequired,
		},
		{
			id:                taskId,
			name:              "Zero time",
			title:             "title",
			description:       "description",
			repeating:         false,
			repeatingInterval: time.Duration(0),
			dayOfWeek:         1,
			time:              time.Time{},
			error:             domain_errors.ErrTimeRequired,
		},
		{
			id:                taskId,
			name:              "Valid task",
			title:             "title",
			description:       "description",
			repeating:         true,
			repeatingInterval: time.Duration(0),
			dayOfWeek:         1,
			time:              time.Now(),
			error:             nil,
		},
		{
			id:                taskId,
			name:              "Valid non-repeating task",
			title:             "title",
			description:       "description",
			repeating:         false,
			repeatingInterval: time.Duration(0),
			dayOfWeek:         1,
			time:              time.Now(),
			error:             nil,
		},
		{
			id: &TaskID{
				primaryId:   originalTaskIDGUID,
				secondaryId: originalTaskIDGUID,
				original:    false,
			},
			name:              "Invalid copy task ID",
			title:             "title",
			description:       "description",
			repeating:         false,
			repeatingInterval: time.Duration(0),
			dayOfWeek:         1,
			time:              time.Now(),
			error:             domain_errors.ErrInvalidTaskID,
		},
		{
			id: &TaskID{
				primaryId:   originalTaskIDGUID,
				secondaryId: originalTaskIDGUID,
				original:    false,
			},
			name:              "Invalid original task ID",
			title:             "title",
			description:       "description",
			repeating:         false,
			repeatingInterval: time.Duration(0),
			dayOfWeek:         1,
			time:              time.Now(),
			error:             domain_errors.ErrInvalidTaskID,
		},
	}

	// Run test cases
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			task, err := NewTask(
				tc.id,
				tc.title, tc.description,
				tc.repeating, tc.repeatingInterval,
				tc.dayOfWeek, tc.time)
			assert.ErrorIs(t, err, tc.error)

			if err != nil {
				assert.ErrorIs(t, err, domain_errors.ErrInvalidTask)
			}

			if err == nil {
				assert.Equal(t, *tc.id, task.GetID())
				assert.Equal(t, tc.title, task.GetTitle())
				assert.Equal(t, tc.description, task.GetDescription())
				isRepeating, repeatingInterval := task.IsRepeating()
				assert.Equal(t, tc.repeating, isRepeating)
				assert.Equal(t, tc.repeatingInterval, repeatingInterval)
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

func TestSearchRepetition(t *testing.T) {
	tests := []struct {
		name       string
		repeating  bool
		time       time.Time
		interval   time.Duration
		cancel     bool
		wantLength int
	}{
		{
			name:       "Repeating task",
			repeating:  true,
			time:       time.Date(2024, time.January, 1, 0, 0, 0, 0, time.UTC),
			interval:   24 * time.Hour,
			cancel:     false,
			wantLength: 31,
		},
		{
			name:       "Non-repeating task",
			time:       time.Date(2024, time.January, 1, 0, 0, 0, 0, time.UTC),
			repeating:  false,
			interval:   24 * time.Hour,
			cancel:     false,
			wantLength: 1,
		},
		{
			name:       "Cancelled context",
			time:       time.Date(2024, time.January, 1, 0, 0, 0, 0, time.UTC),
			repeating:  true,
			interval:   24 * time.Hour,
			cancel:     true,
			wantLength: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			task := &Task{
				time:              tt.time,
				repeating:         tt.repeating,
				repeatingInterval: tt.interval,
			}

			ctx := context.Background()
			if tt.cancel {
				var cancel context.CancelFunc
				ctx, cancel = context.WithCancel(ctx)
				cancel()
			}

			timeChan := task.searchRepetition(ctx)

			gotLength := 0
			for range timeChan {
				gotLength++
			}

			assert.Equal(t, tt.wantLength, gotLength)
		})
	}
}

type MockTaskIDFactory struct {
	mock.Mock
}

func (m *MockTaskIDFactory) CreateTaskID(primaryID string) *TaskID {
	args := m.Called(primaryID)
	return args.Get(0).(*TaskID)
}

func TestCreateTasksFromDates(t *testing.T) {
	// Setup
	ogTaskID := uuid.New()
	times := []time.Time{
		time.Date(2022, 1, 1, 0, 0, 0, 0, time.UTC),
		time.Date(2022, 1, 2, 0, 0, 0, 0, time.UTC),
	}

	task := &Task{
		id: &TaskID{
			primaryId:   ogTaskID,
			secondaryId: ogTaskID,
			original:    true},
		title:             "Task1",
		description:       "Description1",
		time:              time.Date(2022, 1, 1, 0, 0, 0, 0, time.UTC),
		repeating:         true,
		repeatingInterval: time.Duration(24 * time.Hour),
		dayOfWeek:         1,
	}

	tests := []struct {
		name          string
		dates         []time.Time
		taskIDFactory *MockTaskIDFactory
		expectedDates []time.Time
		cancelContext bool
	}{
		{
			name:          "Create tasks from dates",
			dates:         times,
			taskIDFactory: &MockTaskIDFactory{},
			expectedDates: times,
			cancelContext: false,
		},
		{
			name:          "Context cancelled",
			dates:         times,
			taskIDFactory: &MockTaskIDFactory{},
			expectedDates: nil,
			cancelContext: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()
			if tt.cancelContext {
				var cancel context.CancelFunc
				ctx, cancel = context.WithCancel(ctx)
				cancel()
			}

			datesChan := make(chan time.Time, len(tt.dates))
			for _, date := range tt.dates {
				datesChan <- date
			}
			close(datesChan)

			tt.taskIDFactory.On("CreateTaskID", mock.Anything).Return(&TaskID{
				primaryId:   ogTaskID,
				secondaryId: uuid.New(),
				original:    false,
			})

			tasksChan := task.createTasksFromDates(ctx, datesChan, tt.taskIDFactory)

			for newTask := range tasksChan {
				assert.Contains(t, tt.expectedDates, newTask.GetTime())
			}

			if tt.cancelContext {
				assert.Empty(t, tasksChan)
			} else {
				tt.taskIDFactory.AssertExpectations(t)
			}

		})
	}
}
