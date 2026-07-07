package task

import (
	"context"
	"fmt"
)

// PrintTask affiche un message sur la sortie standard.
type PrintTask struct {
	TaskID  string
	Message string
}

func (t *PrintTask) ID() string {
	return t.TaskID
}

func (t *PrintTask) Execute(ctx context.Context) error {
	select {
	case <-ctx.Done():
		return &TaskError{Code: 1, TaskID: t.TaskID, Err: ctx.Err()}
	default:
	}
	fmt.Println(t.Message)
	return nil
}
