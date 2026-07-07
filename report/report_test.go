package report

import (
	"bytes"
	"encoding/json"
	"strings"
	"testing"
)

func TestReportWriteTo(t *testing.T) {
	r := Report{
		Results: []TaskResult{
			{ID: "t1", Status: "success", Duration: "12ms", Attempts: 1},
			{ID: "t2", Status: "failed", Duration: "50ms", Attempts: 3},
		},
	}

	var buf bytes.Buffer
	n, err := r.WriteTo(&buf)
	if err != nil {
		t.Fatalf("WriteTo a échoué: %v", err)
	}
	if n == 0 {
		t.Fatal("WriteTo a écrit 0 octets")
	}

	// vérifier que le JSON est valide
	var parsed Report
	if err := json.Unmarshal(buf.Bytes(), &parsed); err != nil {
		t.Fatalf("le JSON produit est invalide: %v", err)
	}

	if len(parsed.Results) != 2 {
		t.Fatalf("attendu 2 résultats, obtenu %d", len(parsed.Results))
	}

	if parsed.Results[0].ID != "t1" {
		t.Errorf("premier résultat attendu t1, obtenu %s", parsed.Results[0].ID)
	}
}

func TestWriteMetrics(t *testing.T) {
	results := []TaskResult{
		{ID: "t1", Status: "success"},
		{ID: "t2", Status: "failed"},
		{ID: "t3", Status: "timeout"},
		{ID: "t4", Status: "success"},
	}

	md := WriteMetrics(results)

	if !strings.Contains(md, "Tâches exécutées : 4") {
		t.Error("devrait contenir le nombre total de tâches")
	}
	if !strings.Contains(md, "Tâches réussies : 2") {
		t.Error("devrait contenir le nombre de succès")
	}
	if !strings.Contains(md, "Tâches en échec : 1") {
		t.Error("devrait contenir le nombre d'échecs")
	}
	if !strings.Contains(md, "Tâches en timeout : 1") {
		t.Error("devrait contenir le nombre de timeouts")
	}
	if !strings.Contains(md, "Goroutines actives") {
		t.Error("devrait contenir le nombre de goroutines")
	}
}
