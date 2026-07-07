package orchestrator

import (
	"context"
	"errors"
	"fmt"
	"os"
	"sync"
	"time"

	"taskrunner/report"
	"taskrunner/task"
)

// OrchestratorConfig contient la configuration de l'orchestrateur.
type OrchestratorConfig struct {
	Workers int
	Verbose bool
}

// Option est une functional option pour configurer l'orchestrateur.
type Option func(*OrchestratorConfig)

// WithWorkers configure le nombre de workers.
func WithWorkers(n int) Option {
	return func(c *OrchestratorConfig) {
		c.Workers = n
	}
}

// WithVerbose active le mode verbose.
func WithVerbose(v bool) Option {
	return func(c *OrchestratorConfig) {
		c.Verbose = v
	}
}

// ValidateWorkers vérifie que le nombre de workers est valide.
func ValidateWorkers(n int) (int, error) {
	if n < 1 || n > 100 {
		return 3, fmt.Errorf("nombre de workers invalide (%d), utilisation de la valeur par défaut 3", n)
	}
	return n, nil
}

// job est une tâche à exécuter avec ses paramètres.
type job struct {
	task    task.Task
	timeout time.Duration
	retries int
}

// Orchestrate exécute les tâches avec un pool de workers.
func Orchestrate(ctx context.Context, tasks []task.LoadedTask, workers int, opts ...Option) (report.Report, error) {
	cfg := &OrchestratorConfig{
		Workers: workers,
		Verbose: false,
	}
	for _, opt := range opts {
		opt(cfg)
	}

	jobs := make(chan job, len(tasks))
	resultsCh := make(chan report.TaskResult, len(tasks))

	var wg sync.WaitGroup

	// lancer les workers
	for i := 0; i < cfg.Workers; i++ {
		wg.Add(1)
		go func(workerID int) {
			defer wg.Done()
			for j := range jobs {
				res := executeWithRetry(ctx, j, cfg.Verbose)
				resultsCh <- res
			}
		}(i)
	}

	// envoyer les jobs
	for _, lt := range tasks {
		jobs <- job{
			task:    lt.Task,
			timeout: lt.Timeout,
			retries: lt.Retries,
		}
	}
	close(jobs)

	// attendre la fin des workers puis fermer le channel de résultats
	go func() {
		wg.Wait()
		close(resultsCh)
	}()

	// collecter les résultats
	var results []report.TaskResult
	for r := range resultsCh {
		results = append(results, r)
	}

	return report.Report{Results: results}, nil
}

func executeWithRetry(ctx context.Context, j job, verbose bool) report.TaskResult {
	maxAttempts := j.retries + 1
	var lastErr error

	start := time.Now()

	for attempt := 1; attempt <= maxAttempts; attempt++ {
		if verbose {
			fmt.Fprintf(os.Stderr, "[%s] tentative %d/%d...\n", j.task.ID(), attempt, maxAttempts)
		}

		err := runWithTimeout(ctx, j.task, j.timeout)
		if err == nil {
			if verbose {
				fmt.Fprintf(os.Stderr, "[%s] succès\n", j.task.ID())
			}
			return report.TaskResult{
				ID:       j.task.ID(),
				Status:   "success",
				Duration: time.Since(start).String(),
				Attempts: attempt,
			}
		}

		lastErr = err
		if verbose {
			fmt.Fprintf(os.Stderr, "[%s] échec: %v\n", j.task.ID(), err)
		}

		// si le context principal est annulé, on arrête les retries
		if ctx.Err() != nil {
			break
		}
	}

	status := "failed"
	if errors.Is(lastErr, context.DeadlineExceeded) {
		status = "timeout"
	}

	return report.TaskResult{
		ID:       j.task.ID(),
		Status:   status,
		Duration: time.Since(start).String(),
		Attempts: maxAttempts,
	}
}

func runWithTimeout(parent context.Context, t task.Task, timeout time.Duration) error {
	ctx, cancel := context.WithTimeout(parent, timeout)
	defer cancel()
	return t.Execute(ctx)
}
