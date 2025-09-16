package persisters

import (
	"context"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"os"
	"path/filepath"

	"github.com/goccy/go-json"

	"github.com/k6zma/avito-lab1/internal/domain/models"
	persisters2 "github.com/k6zma/avito-lab1/internal/domain/persisters"
)

type jsonSnapshot struct {
	Students []*models.Student `json:"students"`
}

type JSONStudentPersister struct {
	path string
}

func NewJSONStudentPersister(path string) *JSONStudentPersister {
	return &JSONStudentPersister{
		path: path,
	}
}

func (p *JSONStudentPersister) Save(_ context.Context, students []*models.Student) error {
	dir := filepath.Dir(p.path)
	if err := os.MkdirAll(dir, 0o750); err != nil {
		return fmt.Errorf("failed to create directory with json file: %w", err)
	}

	tmp, err := os.CreateTemp(dir, ".students-*.tmp")
	if err != nil {
		return fmt.Errorf("failed to create temp file to save snapshot data: %w", err)
	}

	defer func(name string) {
		if err := os.Remove(name); err != nil && !errors.Is(err, os.ErrNotExist) {
			slog.Error(
				"failed to remove temp file with snapshot data",
				slog.String("path", name),
				slog.Any("error", err),
			)
		}
	}(tmp.Name())

	payload, err := json.Marshal(jsonSnapshot{Students: students})
	if err != nil {
		if closeErr := tmp.Close(); closeErr != nil {
			slog.Error(
				"failed to close temp file with snapshot after marshal error",
				slog.String("path", tmp.Name()),
				slog.Any("error", closeErr),
			)
		}

		return fmt.Errorf("failed to marshal json with snapshot data: %w", err)
	}

	if n, err := tmp.Write(payload); err != nil {
		if closeErr := tmp.Close(); closeErr != nil {
			slog.Error(
				"failed to close temp file with snapshot after write error",
				slog.String("path", tmp.Name()),
				slog.Any("error", closeErr),
			)
		}

		return fmt.Errorf("failed to write marshalled snapshot data: %w", err)
	} else if n != len(payload) {
		if closeErr := tmp.Close(); closeErr != nil {
			slog.Error(
				"failed to close temp file with snapshot after short write",
				slog.String("path", tmp.Name()),
				slog.Any("err", closeErr),
			)
		}

		return persisters2.ErrMismatchPayloadAndWriteLen
	}

	if err := tmp.Sync(); err != nil {
		if closeErr := tmp.Close(); closeErr != nil {
			slog.Error(
				"failed to close temp file with snapshot after sync error",
				slog.String("path", tmp.Name()),
				slog.Any("error", closeErr),
			)
		}

		return fmt.Errorf("failed to sync temp file with snapshot data: %w", err)
	}

	if err := tmp.Close(); err != nil {
		return fmt.Errorf("failed to close temp file with snapshot data: %w", err)
	}

	if err := os.Rename(tmp.Name(), p.path); err != nil {
		return fmt.Errorf("failed to move snapshot file into final destination: %w", err)
	}

	return nil
}

func (p *JSONStudentPersister) Load(_ context.Context) ([]*models.Student, error) {
	file, err := os.Open(p.path)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return nil, nil
		}

		return nil, fmt.Errorf("failed to open json snapshot file: %w", err)
	}

	defer func() {
		if closeErr := file.Close(); closeErr != nil {
			slog.Error(
				"failed to close json snapshot file after load",
				slog.String("path", p.path),
				slog.Any("err", closeErr),
			)
		}
	}()

	data, err := io.ReadAll(file)
	if err != nil {
		return nil, fmt.Errorf("failed to read json snapshot file: %w", err)
	}

	if len(data) == 0 {
		return nil, nil
	}

	var snap jsonSnapshot
	if err := json.Unmarshal(data, &snap); err != nil {
		return nil, fmt.Errorf("failed to unmarshal json snapshot: %w", err)
	}

	return snap.Students, nil
}
