package files

import (
	"fmt"
	"path/filepath"
)

func CreateREADME(mgr *Manager) (err error) {
	file, err := mgr.CreateFile("README.md")
	if err != nil {
		return err
	}
	defer func() {
		if closeErr := file.Close(); err == nil {
			err = closeErr
		}
	}()
	_, err = fmt.Fprintf(file, "# %s\n", filepath.Base(filepath.Clean(mgr.RootDir)))
	return err
}
