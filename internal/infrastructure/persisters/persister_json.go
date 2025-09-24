package persisters

import (
	"errors"
	"fmt"
	"io"
	"log/slog"
	"os"
	"path/filepath"

	"github.com/goccy/go-json"

	"github.com/k6zma/avito-lab1/internal/domain/models"
	"github.com/k6zma/avito-lab1/internal/infrastructure/ciphers"
)

type jsonSnapshot struct {
	Students []*models.Student `json:"students"`
}

type JSONStudentPersister struct {
	path   string
	cipher ciphers.Cipher
}

func NewJSONStudentPersister(path string, c ciphers.Cipher) *JSONStudentPersister {
	return &JSONStudentPersister{
		path:   path,
		cipher: c,
	}
}

func (p *JSONStudentPersister) Save(students []*models.Student) error {
	if p.cipher == nil {
		return ErrInvalidCipher
	}

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

	ciphertext, err := p.cipher.Encrypt(payload)
	if err != nil {
		return fmt.Errorf("failed to encrypt snapshot: %w", err)
	}

	if n, err := tmp.Write(ciphertext); err != nil {
		if closeErr := tmp.Close(); closeErr != nil {
			slog.Error("failed to close temp file after write error",
				slog.String("path", tmp.Name()),
				slog.Any("error", closeErr),
			)
		}
		return fmt.Errorf("failed to write encrypted snapshot data: %w", err)
	} else if n != len(ciphertext) {
		if closeErr := tmp.Close(); closeErr != nil {
			slog.Error("failed to close temp file after short write",
				slog.String("path", tmp.Name()),
				slog.Any("error", closeErr),
			)
		}
		return ErrMismatchPayloadAndWriteLen
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

func (p *JSONStudentPersister) Load() ([]*models.Student, error) {
	if p.cipher == nil {
		return nil, ErrInvalidCipher
	}

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

	plaintext, err := p.cipher.Decrypt(data)
	if err != nil {
		return nil, fmt.Errorf("failed to decrypt snapshot: %w", err)
	}

	var snap jsonSnapshot
	if err := json.Unmarshal(plaintext, &snap); err != nil {
		return nil, fmt.Errorf("failed to unmarshal json snapshot: %w", err)
	}

	return snap.Students, nil
}
