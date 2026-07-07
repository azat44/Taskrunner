package task

import (
	"context"
	"io"
	"net/http"
	"os"
)

// DownloadTask télécharge un fichier depuis une URL.
type DownloadTask struct {
	TaskID string
	URL    string
	Dest   string
}

func (t *DownloadTask) ID() string {
	return t.TaskID
}

func (t *DownloadTask) Execute(ctx context.Context) error {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, t.URL, nil)
	if err != nil {
		return &TaskError{Code: 1, TaskID: t.TaskID, Err: err}
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return &TaskError{Code: 2, TaskID: t.TaskID, Err: err}
	}
	defer resp.Body.Close()

	f, err := os.Create(t.Dest)
	if err != nil {
		return &TaskError{Code: 3, TaskID: t.TaskID, Err: err}
	}
	defer f.Close()

	if _, err = io.Copy(f, resp.Body); err != nil {
		return &TaskError{Code: 4, TaskID: t.TaskID, Err: err}
	}

	return nil
}
