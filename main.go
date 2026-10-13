package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/avycado13/kicad-init/files"
)

func main() {
	namePtr := flag.String("name", "", "Project Name to Initialize")

	flag.Parse()

	if *namePtr == "" {
		fmt.Println("Error: the --name flag is required.")
		flag.Usage()
		os.Exit(1)
	}

	fileMgr := files.NewManager(*namePtr)

	// Fixes the "declared and not used" error:
	// Replace 'SomeMethod' with an actual function inside your files package,
	// or print it out to test that it works.
	fmt.Printf("Initialized manager for project: %v\n", fileMgr.RootDir)
	files.CreateKicadGitignore(fileMgr, []string{})
}
