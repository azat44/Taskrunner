package task

import (
	"context"
	"errors"
	"time"
)

// FakeTaskBehavior définit le comportement simulé.
type FakeTaskBehavior int

const (
	BehaviorSuccess FakeTaskBehavior = iota
	BehaviorFail
	BehaviorTimeout
)

// FakeTask simule une tâche avec un comportement configurable.
type FakeTask struct {
	TaskID   string
	Behavior FakeTaskBehavior
	Delay    time.Duration
}

// NewFakeTask crée une FakeTask avec les paramètres donnés.
func NewFakeTask(id string, behavior FakeTaskBehavior, delay time.Duration) *FakeTask {
	return &FakeTask{TaskID: id, Behavior: behavior, Delay: delay}
}

func (t *FakeTask) ID() string {
	return t.TaskID
}

func (t *FakeTask) Execute(ctx context.Context) error {
	// attend le délai ou l'annulation du context
	select {
	case <-time.After(t.Delay):
	case <-ctx.Done():
		return &TaskError{Code: 10, TaskID: t.TaskID, Err: ctx.Err()}
	}

	switch t.Behavior {
	case BehaviorFail:
		return &TaskError{Code: 11, TaskID: t.TaskID, Err: errors.New("fake failure")}
	case BehaviorTimeout:
		// bloque jusqu'au timeout du context
		select {
		case <-time.After(10 * time.Second):
		case <-ctx.Done():
			return &TaskError{Code: 10, TaskID: t.TaskID, Err: ctx.Err()}
		}
	}
	return nil
}
