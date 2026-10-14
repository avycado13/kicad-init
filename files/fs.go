// Package files supplies operations for creating project files.
package files

import (
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
	file, err := os.Create(fullPath)
	if err != nil {
		return nil, fmt.Errorf("create %q: %w", fullPath, err)
	}
	return file, nil
}
