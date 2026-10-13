package files

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
)

type Manager struct {
	RootDir string
}

func NewManager(dirname string) *Manager {
	err := os.MkdirAll(dirname, 0755)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Directory %s created\n", dirname)

	return &Manager{RootDir: dirname}
}

func (m *Manager) mkDir(dirname string) {
	fullPath := filepath.Join(m.RootDir, dirname)
	err := os.MkdirAll(fullPath, 0755)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Directory %s created\n", dirname)
}

func (m *Manager) CreateFile(path string) (*os.File, error) {
	fullPath := filepath.Join(m.RootDir, path)
	file, err := os.Create(fullPath)
	if err != nil {
		fmt.Println("Error creating file:", err)
		return file, err
	}
	return file, err
}
