package domain

import (
	"fmt"

	"github.com/google/uuid"
)

// TaskID is the unique identifier for a task.
//
// It is composed of a primaryId and a secondaryId.
//
// The primaryId is the original identifier for the task.
//
// The secondaryId is the identifier for the child tasks.
//
// The original flag indicates if the task is the original or a copy.
//
// The original task has the same primaryId and secondaryId.
type TaskID struct {
	// primaryId is the original task identifier
	primaryId uuid.UUID
	// secondaryId is the identifier for childs tasks
	secondaryId uuid.UUID
	// original indicates if the task is the original or a copy
	original bool
}

func (ti *TaskID) String() string {
	var taskType string

	if ti.original {
		taskType = "original"
	} else {
		taskType = "copy"
	}

	return fmt.Sprintf("%s-%s-%s", ti.primaryId, ti.secondaryId.String(), taskType)
}

func NewTaskID() *TaskID {
	id := uuid.New()
	return &TaskID{
		primaryId:   id,
		secondaryId: id,
		original:    true,
	}
}

// TaskIDFactory is the abstract factory interface for creating TaskID instances
type TaskIDFactory interface {
	CreateTaskID(identifier string) *TaskID
}

// OriginalTaskIDFactory is a concrete factory implementation for creating original TaskID instances
type OriginalTaskIDFactory struct{}

// CreateTaskID creates an original TaskID instance
func (f *OriginalTaskIDFactory) CreateTaskID(identifier string) *TaskID {
	if identifier == "" {
		return nil
	}

	parsedIdentifier, err := uuid.Parse(identifier)

	if err != nil {
		return nil
	}

	return &TaskID{
		primaryId:   parsedIdentifier,
		secondaryId: parsedIdentifier,
		original:    true,
	}
}

// CopyTaskIDFactory is a concrete factory implementation for creating copy TaskID instances
type CopyTaskIDFactory struct{}

// CreateTaskID creates a copy TaskID instance
func (f *CopyTaskIDFactory) CreateTaskID(identifier string) *TaskID {
	if identifier == "" {
		return nil
	}

	parsedIdentifier, err := uuid.Parse(identifier)

	if err != nil {
		return nil
	}

	return &TaskID{
		primaryId:   parsedIdentifier,
		secondaryId: uuid.New(),
		original:    false,
	}
}
