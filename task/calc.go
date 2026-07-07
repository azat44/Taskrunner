package task

import (
	"context"
	"fmt"
	"math"
)

// CalcTask effectue un calcul simple sur une valeur.
type CalcTask struct {
	TaskID string
	Value  float64
}

func (t *CalcTask) ID() string {
	return t.TaskID
}

func (t *CalcTask) Execute(ctx context.Context) error {
	select {
	case <-ctx.Done():
		return &TaskError{Code: 1, TaskID: t.TaskID, Err: ctx.Err()}
	default:
	}

	result := math.Sqrt(t.Value)
	fmt.Printf("calc %s: sqrt(%g) = %g\n", t.TaskID, t.Value, result)
	return nil
}
