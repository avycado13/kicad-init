package files

func CreateKicadGitignore(mgr *Manager, extraPaths []string) (err error) {
	paths := []string{
		"*.000",
		"*.bak",
		"*.bck",
		"*.kicad_pcb-bak",
		"*.kicad_sch-bak",
		"*-backups",
		"*-cache*",
		"*-bak",
		"*-bak*",
		"*~",
		"~*",
		"_autosave-*",
		"#auto_saved_files#",
		"*.tmp",
		"*-save.pro",
		"*-save.kicad_pcb",
		"fp-info-cache",
		"~*.lck",
		"#auto_saved_files#",
		"*.net",
		"*.dsn",
		"*.ses",
		"**/*-backups/*.zip",
		"*.kicad_prl",
		".DS_Store",
		"Thumbs.db",
	}
	file, err := mgr.CreateFile(".gitignore")
	if err != nil {
		return err
	}
	defer func() {
		if closeErr := file.Close(); err == nil {
			err = closeErr
		}
	}()
	for _, value := range append(paths, extraPaths...) {
		if _, err = file.WriteString(value + "\n"); err != nil {
			return err
		}
	}
	return nil
}
