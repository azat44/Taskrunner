package task

import (
	"context"
	"fmt"
)

// Task représente une tâche exécutable par l'orchestrateur.
type Task interface {
	ID() string
	Execute(ctx context.Context) error
}

// TaskError est le type d'erreur custom pour toutes les tâches.
type TaskError struct {
	Code   int
	TaskID string
	Err    error
}

func (e *TaskError) Error() string {
	return fmt.Sprintf("task %s (code %d): %v", e.TaskID, e.Code, e.Err)
}

func (e *TaskError) Unwrap() error {
	return e.Err
}
