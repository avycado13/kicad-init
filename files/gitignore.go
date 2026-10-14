package files

import "fmt"

func CreateKicadGitignore(mgr *Manager, extraPaths []string) {
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
	combined := append(paths, extraPaths...)
	file, err := mgr.CreateFile(".gitignore")
	if err != nil {
		fmt.Println("Error Creating Gitignore:", err)
	}

	defer file.Close()
	for _, value := range combined {
		_, err = file.WriteString(value + "\n")
		if err != nil {
			fmt.Println("Error writing path to gitignore:", err)
		}
	}
	fmt.Printf("Gitignore created: \n")
}
