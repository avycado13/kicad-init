// Package actions houses all miscellaneous actions that do not modify the filesystem significantly
package actions

import (
	"fmt"

	"github.com/cli/go-gh/v2"
)

type GithubRepoSettings struct {
	Name        string
	Description string
	Public      bool
}

func GenerateGithubRepoSettings() GithubRepoSettings {
	return GithubRepoSettings{
		Name:        "",
		Description: "",
		Public:      true,
	}
}

func CreateGithubRepo() (err error) {
	settings := GenerateGithubRepoSettings()
	publicString := "--public"
	if !settings.Public {
		publicString = "--private"
	}
	gh.Exec("repo", "create", settings.Name, "--description", settings.Description, publicString)
	if err != nil {
		return fmt.Errorf(": %w", err)
	}
	return nil
}
