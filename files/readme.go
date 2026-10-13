package files

import "fmt"

func CreateREADME(mgr *Manager) {
	// Create README File
	file, err := mgr.CreateFile("README.md")

	defer file.Close()
	file.WriteString("# " + mgr.RootDir + "\n") // Heading with Project name inferred from name/root directory argument

}
