package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"taskrunner/orchestrator"
	"taskrunner/report"
	"taskrunner/task"
)

func main() {
	filePath := flag.String("file", "", "chemin vers le fichier de tâches JSON")
	workers := flag.Int("workers", 3, "nombre de workers simultanés")
	verbose := flag.Bool("verbose", false, "affiche le statut en temps réel sur stderr")
	flag.Parse()

	if *filePath == "" {
		fmt.Fprintln(os.Stderr, "erreur: le flag -file est obligatoire")
		flag.Usage()
		os.Exit(1)
	}

	// valider le nombre de workers
	w, err := orchestrator.ValidateWorkers(*workers)
	if err != nil {
		fmt.Fprintf(os.Stderr, "attention: %v\n", err)
	}

	// charger les tâches
	loaded, err := task.LoadTasks(*filePath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "erreur chargement: %v\n", err)
		os.Exit(1)
	}

	// context avec gestion du signal SIGINT
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT)

	go func() {
		<-sigCh
		fmt.Fprintln(os.Stderr, "\nsignal reçu, arrêt en cours...")
		cancel()
	}()

	// exécuter les tâches
	rep, err := orchestrator.Orchestrate(ctx, loaded, w, orchestrator.WithVerbose(*verbose))
	if err != nil {
		fmt.Fprintf(os.Stderr, "erreur orchestration: %v\n", err)
	}

	// écrire le rapport JSON sur stdout
	if _, err := rep.WriteTo(os.Stdout); err != nil {
		fmt.Fprintf(os.Stderr, "erreur écriture rapport: %v\n", err)
		os.Exit(1)
	}

	// générer le fichier METRICS.md
	metrics := writeMetricsFile(rep.Results)
	if err := os.WriteFile("METRICS.md", []byte(metrics), 0644); err != nil {
		fmt.Fprintf(os.Stderr, "erreur écriture METRICS.md: %v\n", err)
	}
}

func writeMetricsFile(results []report.TaskResult) string {
	return report.WriteMetrics(results)
}
