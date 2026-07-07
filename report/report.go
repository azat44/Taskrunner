package report

import (
	"encoding/json"
	"fmt"
	"io"
	"runtime"
)

// TaskResult contient le résultat d'une tâche exécutée.
type TaskResult struct {
	ID       string `json:"id"`
	Status   string `json:"status"`
	Duration string `json:"duration"`
	Attempts int    `json:"attempts"`
}

// Report contient tous les résultats.
type Report struct {
	Results []TaskResult `json:"results"`
}

// WriteTo implémente io.WriterTo : sérialise le rapport en JSON.
func (r Report) WriteTo(w io.Writer) (int64, error) {
	data, err := json.MarshalIndent(r, "", "  ")
	if err != nil {
		return 0, fmt.Errorf("sérialisation rapport: %w", err)
	}
	data = append(data, '\n')
	n, err := w.Write(data)
	return int64(n), err
}

// WriteMetrics génère le contenu Markdown du fichier METRICS.md.
func WriteMetrics(results []TaskResult) string {
	total := len(results)
	success := 0
	failed := 0
	timeout := 0

	for _, r := range results {
		switch r.Status {
		case "success":
			success++
		case "failed":
			failed++
		case "timeout":
			timeout++
		}
	}

	goroutines := runtime.NumGoroutine()

	return fmt.Sprintf(`# Métriques d'exécution

- Goroutines actives à la fin : %d
- Tâches exécutées : %d
- Tâches réussies : %d
- Tâches en échec : %d
- Tâches en timeout : %d
`, goroutines, total, success, failed, timeout)
}
