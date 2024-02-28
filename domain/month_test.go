package domain

import (
	"testing"
	"time"

	"github.com/google/uuid"
	domain_errors "github.com/sosalejandro/go-calendar/domain/errors"
	"github.com/stretchr/testify/assert"
)

func TestMonthValidate(t *testing.T) {
	tests := []struct {
		name    string
		month   *Month
		wantErr bool
		err     error
	}{
		{
			name: "Valid month",
			month: &Month{
				month: time.May,
				year:  2022,
			}, // May
			wantErr: false,
		},
		{
			name: "Month too low",
			month: &Month{
				month: 0,
				year:  2022,
			},
			wantErr: true,
			err:     domain_errors.ErrInvalidMonth,
		},
		{
			name: "Month too high",
			month: &Month{
				month: 13,
				year:  2022,
			},
			wantErr: true,
			err:     domain_errors.ErrInvalidMonth,
		},
		{
			name: "Year too low",
			month: &Month{
				month: time.February,
				year:  0,
			},
			wantErr: true,
			err:     domain_errors.ErrInvalidYear,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.month.Validate()
			if tt.wantErr {
				assert.Error(t, err)
				assert.ErrorIs(t, err, tt.err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestNewMonth(t *testing.T) {
	tests := []struct {
		name    string
		month   time.Month
		year    int
		wantErr bool
	}{
		{
			name:    "Valid month",
			month:   time.May, // May
			year:    2022,
			wantErr: false,
		},
		{
			name:    "Month too low",
			month:   0,
			year:    2022,
			wantErr: true,
		},
		{
			name:    "Month too high",
			month:   13,
			year:    2022,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := NewMonth(tt.month, tt.year)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestMonthAddDay(t *testing.T) {
	tests := []struct {
		name        string
		month       *Month
		dayToAdd    int
		expected    *Month
		expectedErr bool
		err         error
	}{
		{
			name: "Add day within same month",
			month: &Month{
				month: time.May,
				year:  2022,
				days:  map[int]*Day{},
			},
			dayToAdd: 15,
			expected: &Month{
				month: time.May,
				year:  2022,
				days: map[int]*Day{
					15: &Day{
						day:   15,
						tasks: []*Task{},
					},
				},
			}, // Still May
		},
		{
			name: "Add day to next month throws error",
			month: &Month{
				month: time.December,
				year:  2022,
				days:  map[int]*Day{},
			}, // December
			dayToAdd: 32,
			expected: &Month{
				month: time.December,
				year:  2022,
				days:  map[int]*Day{},
			},
			expectedErr: true,
			err:         domain_errors.ErrAddDay,
		},
		{
			name: "Add day returns existing day",
			month: &Month{
				month: time.May,
				year:  2022,
				days: map[int]*Day{
					15: &Day{
						day:   15,
						tasks: []*Task{},
					},
				},
			},
			dayToAdd: 15,
			expected: &Month{
				month: time.May,
				year:  2022,
				days: map[int]*Day{
					15: &Day{
						day:   15,
						tasks: []*Task{},
					},
				},
			}, // Still May
		},
		// {
		// 	name: "Add day to next month",
		// 	month: &Month{
		// 		month: time.January,
		// 		year:  2022,
		// 		days:  map[int]*Day{},
		// 	}, // January
		// 	dayToAdd: 32,
		// 	expected: &Month{
		// 		month: time.February,
		// 		year:  2022,
		// 		days: map[int]*Day{
		// 			15: &Day{
		// 				day:   15,
		// 				tasks: []*Task{},
		// 			},
		// 		},
		// 	}, // February
		// },
		// {
		// 	name: "Add day to next year",
		// 	month: &Month{
		// 		month: time.December,
		// 		year:  2022,
		// 		days:  map[int]*Day{},
		// 	}, // December
		// 	dayToAdd: 32,
		// 	expected: &Month{
		// 		month: time.January,
		// 		year:  2022,
		// 		days: map[int]*Day{
		// 			15: &Day{
		// 				day:   15,
		// 				tasks: []*Task{},
		// 			},
		// 		},
		// 	}, // January
		// },
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d, err := tt.month.addDay(tt.dayToAdd)

			if tt.expectedErr {
				assert.ErrorIs(t, err, tt.err)
			}

			if err == nil {
				assert.NotNil(t, d)
			}

			assert.Equal(t, tt.expected, tt.month)
		})
	}
}

func TestMonthAddTaskToDay(t *testing.T) {
	task := &Task{
		id: &TaskID{
			primaryId:   uuid.New(),
			secondaryId: uuid.New(),
		},
		title:       "Task1",
		description: "Description1",
	}

	tests := []struct {
		name     string
		month    *Month
		day      int
		task     *Task
		expected error
	}{
		{
			name:     "Add task to existing day",
			month:    &Month{days: map[int]*Day{15: &Day{day: 15, tasks: []*Task{}}}},
			day:      15,
			task:     task,
			expected: nil,
		},
		{
			name:     "Add task to non-existing day",
			month:    &Month{days: map[int]*Day{15: &Day{day: 15, tasks: []*Task{}}}},
			day:      16,
			task:     task,
			expected: domain_errors.ErrDayNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.month.addTaskToDay(tt.day, tt.task)
			assert.Equal(t, tt.expected, err)
		})
	}
}

func TestMonthGetDay(t *testing.T) {
	tests := []struct {
		name     string
		month    *Month
		day      int
		expected *Day
		wantErr  error
	}{
		{
			name:     "Get existing day",
			month:    &Month{days: map[int]*Day{15: &Day{day: 15, tasks: []*Task{}}}},
			day:      15,
			expected: &Day{day: 15, tasks: []*Task{}},
			wantErr:  nil,
		},
		{
			name:     "Get non-existing day",
			month:    &Month{days: map[int]*Day{15: &Day{day: 15, tasks: []*Task{}}}},
			day:      16,
			expected: nil,
			wantErr:  domain_errors.ErrDayNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			day, err := tt.month.getDay(tt.day)
			assert.Equal(t, tt.expected, day)
			if tt.wantErr != nil {
				assert.ErrorIs(t, err, tt.wantErr)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
