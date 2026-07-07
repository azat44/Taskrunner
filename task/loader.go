package task

import (
	"encoding/json"
	"fmt"
	"os"
	"time"
)

// TaskDef représente une tâche telle que lue depuis le JSON.
type TaskDef struct {
	ID      string          `json:"id"`
	Type    string          `json:"type"`
	Params  json.RawMessage `json:"params"`
	Timeout string          `json:"timeout"`
	Retries int             `json:"retries"`
}

// TaskFile est la structure du fichier tasks.json.
type TaskFile struct {
	Tasks []TaskDef `json:"tasks"`
}

// LoadedTask regroupe la tâche, son timeout et ses retries.
type LoadedTask struct {
	Task    Task
	Timeout time.Duration
	Retries int
}

// LoadTasks lit le fichier JSON et instancie les tâches.
func LoadTasks(path string) ([]LoadedTask, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("lecture fichier: %w", err)
	}

	var tf TaskFile
	if err := json.Unmarshal(data, &tf); err != nil {
		return nil, fmt.Errorf("parsing JSON: %w", err)
	}

	loaded := make([]LoadedTask, 0, len(tf.Tasks))
	for _, def := range tf.Tasks {
		t, err := buildTask(def)
		if err != nil {
			return nil, fmt.Errorf("task %s: %w", def.ID, err)
		}

		timeout, err := time.ParseDuration(def.Timeout)
		if err != nil {
			return nil, fmt.Errorf("task %s timeout invalide: %w", def.ID, err)
		}

		loaded = append(loaded, LoadedTask{
			Task:    t,
			Timeout: timeout,
			Retries: def.Retries,
		})
	}

	return loaded, nil
}

func buildTask(def TaskDef) (Task, error) {
	switch def.Type {
	case "print":
		var p struct {
			Message string `json:"message"`
		}
		if err := json.Unmarshal(def.Params, &p); err != nil {
			return nil, fmt.Errorf("params print: %w", err)
		}
		return &PrintTask{TaskID: def.ID, Message: p.Message}, nil

	case "download":
		var p struct {
			URL  string `json:"url"`
			Dest string `json:"dest"`
		}
		if err := json.Unmarshal(def.Params, &p); err != nil {
			return nil, fmt.Errorf("params download: %w", err)
		}
		return &DownloadTask{TaskID: def.ID, URL: p.URL, Dest: p.Dest}, nil

	case "calc":
		var p struct {
			Value float64 `json:"value"`
		}
		if err := json.Unmarshal(def.Params, &p); err != nil {
			return nil, fmt.Errorf("params calc: %w", err)
		}
		return &CalcTask{TaskID: def.ID, Value: p.Value}, nil

	case "fake":
		var p struct {
			Behavior string `json:"behavior"`
			Delay    string `json:"delay"`
		}
		if err := json.Unmarshal(def.Params, &p); err != nil {
			return nil, fmt.Errorf("params fake: %w", err)
		}

		var b FakeTaskBehavior
		switch p.Behavior {
		case "success":
			b = BehaviorSuccess
		case "fail":
			b = BehaviorFail
		case "timeout":
			b = BehaviorTimeout
		default:
			return nil, fmt.Errorf("behavior inconnu: %s", p.Behavior)
		}

		delay, err := time.ParseDuration(p.Delay)
		if err != nil {
			return nil, fmt.Errorf("delay invalide: %w", err)
		}
		return NewFakeTask(def.ID, b, delay), nil

	default:
		return nil, fmt.Errorf("type de tâche inconnu: %s", def.Type)
	}
}
