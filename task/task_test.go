package task

import (
	"context"
	"errors"
	"testing"
	"time"
)

func TestTaskErrorImplementsError(t *testing.T) {
	err := &TaskError{Code: 1, TaskID: "test-1", Err: errors.New("oops")}
	msg := err.Error()
	if msg == "" {
		t.Fatal("Error() ne doit pas retourner une chaîne vide")
	}
}

func TestTaskErrorUnwrap(t *testing.T) {
	original := errors.New("erreur originale")
	err := &TaskError{Code: 1, TaskID: "test-1", Err: original}

	if !errors.Is(err, original) {
		t.Fatal("Unwrap devrait permettre errors.Is de trouver l'erreur originale")
	}
}

func TestPrintTaskExecute(t *testing.T) {
	pt := &PrintTask{TaskID: "p1", Message: "hello"}
	if pt.ID() != "p1" {
		t.Fatalf("ID() attendu p1, obtenu %s", pt.ID())
	}

	err := pt.Execute(context.Background())
	if err != nil {
		t.Fatalf("PrintTask ne devrait pas échouer: %v", err)
	}
}

func TestCalcTaskExecute(t *testing.T) {
	ct := &CalcTask{TaskID: "c1", Value: 16}
	err := ct.Execute(context.Background())
	if err != nil {
		t.Fatalf("CalcTask ne devrait pas échouer: %v", err)
	}
}

func TestFakeTaskSuccess(t *testing.T) {
	ft := NewFakeTask("f1", BehaviorSuccess, 10*time.Millisecond)
	err := ft.Execute(context.Background())
	if err != nil {
		t.Fatalf("FakeTask success ne devrait pas échouer: %v", err)
	}
}

func TestFakeTaskFail(t *testing.T) {
	ft := NewFakeTask("f2", BehaviorFail, 10*time.Millisecond)
	err := ft.Execute(context.Background())
	if err == nil {
		t.Fatal("FakeTask fail devrait retourner une erreur")
	}

	var te *TaskError
	if !errors.As(err, &te) {
		t.Fatal("l'erreur devrait être un *TaskError")
	}
}

func TestFakeTaskTimeout(t *testing.T) {
	ft := NewFakeTask("f3", BehaviorTimeout, 10*time.Millisecond)

	ctx, cancel := context.WithTimeout(context.Background(), 50*time.Millisecond)
	defer cancel()

	err := ft.Execute(ctx)
	if err == nil {
		t.Fatal("FakeTask timeout devrait retourner une erreur")
	}

	if !errors.Is(err, context.DeadlineExceeded) {
		t.Fatalf("attendu DeadlineExceeded, obtenu: %v", err)
	}
}

func TestPrintTaskCancelledContext(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	cancel() // annulé immédiatement

	pt := &PrintTask{TaskID: "p2", Message: "nope"}
	err := pt.Execute(ctx)
	if err == nil {
		t.Fatal("devrait échouer avec un context annulé")
	}
}
