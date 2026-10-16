// Package files supplies operations for creating project files.
package files

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
)

type Manager struct {
	RootDir string
}

func NewManager(dirname string) (*Manager, error) {
	if err := os.MkdirAll(dirname, 0o755); err != nil {
		return nil, fmt.Errorf("create project directory %q: %w", dirname, err)
	}
	return &Manager{RootDir: dirname}, nil
}

func (m *Manager) CreateFile(path string) (*os.File, error) {
	fullPath := filepath.Join(m.RootDir, path)
	file, err := os.OpenFile(fullPath, os.O_WRONLY|os.O_CREATE|os.O_EXCL, 0o644)
	if errors.Is(err, os.ErrExist) {
		return nil, fmt.Errorf("%s already exists: %w", fullPath, os.ErrExist)
	}
	if err != nil {
		return nil, fmt.Errorf("create %q: %w", fullPath, err)
	}
	return file, nil
}
