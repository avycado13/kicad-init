package files

import "fmt"

func CreateREADME(mgr *Manager) {
	// Create README File
	file, err := mgr.CreateFile("README.md")
	if err != nil {
		fmt.Println("Error Creating README:", err)
	}
	defer file.Close()
	_, err = file.WriteString("# " + mgr.RootDir + "\n") // Heading with Project name inferred from name/root directory argument
	if err != nil {
		fmt.Println("Error writing to README:", err)
	}
	fmt.Printf("README created: \n")
}
