package orchestrator

import (
	"context"
	"testing"
	"time"

	"taskrunner/task"
)

func TestValidateWorkers(t *testing.T) {
	tests := []struct {
		input   int
		want    int
		wantErr bool
	}{
		{3, 3, false},
		{1, 1, false},
		{100, 100, false},
		{0, 3, true},
		{-1, 3, true},
		{101, 3, true},
	}

	for _, tt := range tests {
		got, err := ValidateWorkers(tt.input)
		if (err != nil) != tt.wantErr {
			t.Errorf("ValidateWorkers(%d): erreur = %v, attendu erreur = %v", tt.input, err, tt.wantErr)
		}
		if got != tt.want {
			t.Errorf("ValidateWorkers(%d) = %d, attendu %d", tt.input, got, tt.want)
		}
	}
}

func TestOrchestrateSuccess(t *testing.T) {
	tasks := []task.LoadedTask{
		{
			Task:    task.NewFakeTask("ok-1", task.BehaviorSuccess, 10*time.Millisecond),
			Timeout: 1 * time.Second,
			Retries: 0,
		},
		{
			Task:    task.NewFakeTask("ok-2", task.BehaviorSuccess, 10*time.Millisecond),
			Timeout: 1 * time.Second,
			Retries: 0,
		},
	}

	rep, err := Orchestrate(context.Background(), tasks, 2)
	if err != nil {
		t.Fatalf("Orchestrate a échoué: %v", err)
	}

	if len(rep.Results) != 2 {
		t.Fatalf("attendu 2 résultats, obtenu %d", len(rep.Results))
	}

	for _, r := range rep.Results {
		if r.Status != "success" {
			t.Errorf("task %s: attendu success, obtenu %s", r.ID, r.Status)
		}
	}
}

func TestOrchestrateWithRetry(t *testing.T) {
	tasks := []task.LoadedTask{
		{
			Task:    task.NewFakeTask("fail-1", task.BehaviorFail, 5*time.Millisecond),
			Timeout: 1 * time.Second,
			Retries: 2,
		},
	}

	rep, err := Orchestrate(context.Background(), tasks, 1)
	if err != nil {
		t.Fatalf("Orchestrate a échoué: %v", err)
	}

	if len(rep.Results) != 1 {
		t.Fatalf("attendu 1 résultat, obtenu %d", len(rep.Results))
	}

	r := rep.Results[0]
	if r.Status != "failed" {
		t.Errorf("attendu failed, obtenu %s", r.Status)
	}
	if r.Attempts != 3 {
		t.Errorf("attendu 3 tentatives, obtenu %d", r.Attempts)
	}
}

func TestOrchestrateTimeout(t *testing.T) {
	tasks := []task.LoadedTask{
		{
			Task:    task.NewFakeTask("slow-1", task.BehaviorTimeout, 5*time.Millisecond),
			Timeout: 100 * time.Millisecond,
			Retries: 0,
		},
	}

	rep, err := Orchestrate(context.Background(), tasks, 1)
	if err != nil {
		t.Fatalf("Orchestrate a échoué: %v", err)
	}

	r := rep.Results[0]
	if r.Status != "timeout" {
		t.Errorf("attendu timeout, obtenu %s", r.Status)
	}
}

func TestWithVerboseOption(t *testing.T) {
	cfg := &OrchestratorConfig{}
	WithVerbose(true)(cfg)
	if !cfg.Verbose {
		t.Error("WithVerbose(true) devrait mettre Verbose à true")
	}
}

func TestWithWorkersOption(t *testing.T) {
	cfg := &OrchestratorConfig{}
	WithWorkers(5)(cfg)
	if cfg.Workers != 5 {
		t.Errorf("WithWorkers(5) devrait mettre Workers à 5, obtenu %d", cfg.Workers)
	}
}
