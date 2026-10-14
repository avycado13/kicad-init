package files

import (
	"fmt"
	"path/filepath"
	"time"

	"gopkg.in/yaml.v3"
)

type JournalMetadata struct {
	Title       string    `yaml:"title"`
	Github      string    `yaml:"github"`
	Description string    `yaml:"description"`
	CreatedAt   time.Time `yaml:"created_at"`
}

func CreateJournal(mgr *Manager) (err error) {
	now := time.Now()
	metadata := JournalMetadata{
		Title:       filepath.Base(filepath.Clean(mgr.RootDir)),
		Github:      "github/github",
		Description: "Lorem ipsum dolor sit amet, consectetur adipiscing elit",
		CreatedAt:   time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location()),
	}

	frontMatter, err := yaml.Marshal(metadata)
	if err != nil {
		return fmt.Errorf("marshal journal metadata: %w", err)
	}

	file, err := mgr.CreateFile("JOURNAL.md")
	if err != nil {
		return fmt.Errorf("create journal file: %w", err)
	}

	defer func() {
		if closeErr := file.Close(); err == nil {
			err = closeErr
		}
	}()

	_, err = fmt.Fprintf(
		file,
		"---\n%s---\n\n",
		frontMatter,
	)
	if err != nil {
		return fmt.Errorf("write journal: %w", err)
	}

	return nil
}
