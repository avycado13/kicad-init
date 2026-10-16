// Package actions houses all miscellaneous actions that do not modify the filesystem significantly
package actions

import (
	"fmt"
	"os/exec"

	"github.com/cli/go-gh/v2"
	"github.com/cli/go-gh/v2/pkg/auth"
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

func CreateGithubRepo(settings GithubRepoSettings) (err error) {
	if settings.Name == "" {
		return fmt.Errorf("GitHub repository name cannot be empty")
	}
	publicString := "--public"
	if !settings.Public {
		publicString = "--private"
	}
	if _, err := exec.LookPath("gh"); err != nil {
		return fmt.Errorf("GitHub CLI (gh) is not installed")
	}
	token, _ := auth.TokenForHost("github.com")
	if token == "" {
		return fmt.Errorf("not authenticated; run `gh auth login`")
	}

	_, _, err = gh.Exec("repo", "create", settings.Name, "--description", settings.Description, publicString)
	if err != nil {
		return fmt.Errorf("create GitHub repository %q: %w", settings.Name, err)
	}
	return nil
}
