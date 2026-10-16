package main

import (
	"flag"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	tea "charm.land/bubbletea/v2"

	"github.com/avycado13/kicad-init/actions"
	"github.com/avycado13/kicad-init/files"
	"github.com/avycado13/kicad-init/files/kicad"
	"github.com/muesli/mango"
	"github.com/muesli/mango/mflag"
	"github.com/muesli/roff"
)

type actionStatus int

const (
	statusPending actionStatus = iota
	statusRunning
	statusSuccess
	statusFailed
)

const (
	stepSelect = iota
	stepConfigure
	stepReview
	stepRun
	stepDone
)

type projectConfig struct {
	GitHub actions.GithubRepoSettings
}

type actionFinishedMsg struct {
	index int
	err   error
}

type model struct {
	projectName string
	rootDir     string

	mgr *files.Manager

	choices  []string
	cursor   int
	selected map[int]struct{}

	config      projectConfig
	configField int // 0: repository name, 1: visibility
	configError string

	step   int
	status []actionStatus
	errors map[int]error
	queue  []int
}

var (
	version = "dev"
	commit  = "none"
	date    = "unknown"
)

func initialModel(name string) model {
	choices := []string{
		".gitignore",
		"README.md",
		"Git repository",
		"JOURNAL.md",
		"KiCad Project",
		"GitHub Repository",
	}

	return model{
		projectName: name,
		rootDir:     name,
		config: projectConfig{GitHub: actions.GithubRepoSettings{
			Name:   filepath.Base(filepath.Clean(name)),
			Public: true,
		}},
		choices: choices,
		selected: map[int]struct{}{
			0: {},
			1: {},
			2: {},
			3: {},
			4: {},
			5: {},
		},
		status: make([]actionStatus, len(choices)),
		errors: make(map[int]error),
	}
}

func (m model) runAction(index int) error {
	if index < 0 || index >= len(m.choices) {
		return fmt.Errorf("unknown action: %d", index)
	}

	switch index {
	case 0:
		return files.CreateKicadGitignore(m.mgr, nil)
	case 1:
		return files.CreateREADME(m.mgr)
	case 2:
		if _, err := exec.LookPath("git"); err != nil {
			return fmt.Errorf("git is not installed")
		}
		output, err := exec.Command("git", "-C", m.mgr.RootDir, "init").CombinedOutput()
		if err != nil {
			return fmt.Errorf("git init: %w: %s", err, strings.TrimSpace(string(output)))
		}
		return nil
	case 3:
		return files.CreateJournal(m.mgr)
	case 4:
		return kicad.CreateKiCadProject(m.mgr)
	case 5:
		return actions.CreateGithubRepo(m.config.GitHub)
	default:
		return fmt.Errorf("unknown action: %d", index)
	}
}

func runActionCmd(m model, index int) tea.Cmd {
	return func() tea.Msg {
		return actionFinishedMsg{
			index: index,
			err:   m.runAction(index),
		}
	}
}

func (m *model) prepareActions() tea.Cmd {
	m.queue = nil

	for i := range m.choices {
		if _, ok := m.selected[i]; ok {
			m.queue = append(m.queue, i)
		}
	}

	if len(m.queue) == 0 {
		m.step = stepDone
		return nil
	}

	mgr, err := files.NewManager(m.rootDir)
	if err != nil {
		for _, index := range m.queue {
			m.status[index] = statusFailed
			m.errors[index] = err
		}
		m.step = stepDone
		return nil
	}
	m.mgr = mgr

	return m.startNextAction()
}

func (m *model) startNextAction() tea.Cmd {
	index := m.queue[0]
	m.status[index] = statusRunning
	return runActionCmd(*m, index)
}

func main() {
	namePtr := flag.String("name", "", "Project Name/Directory to Initialize")
	help := flag.Bool("help", false, "Show the help message")
	man := flag.Bool("man", false, "Show the man page")
	versionFlag := flag.Bool("version", false, "Show the installed version of kicad-init")
	flag.Parse()

	manPage := mango.NewManPage(1, "kicad-init", "a KiCad project generator").
		WithLongDescription("kicad-init is a project generator and setter-upper for the KiCad EDA.\n"+
			"Features:\n"+
			"* User-friendly\n"+
			"* Configurable").
		WithSection("Copyright", "(C) 2026 avycado13.\n"+
			"Released under MIT license.")

	flag.VisitAll(mflag.FlagVisitor(manPage))

	if *help {
		// fmt.Println(manPage.Build(roff.NewDocument()))
		flag.Usage()
		os.Exit(0)
	}
	if *man {
		fmt.Println(manPage.Build(roff.NewDocument()))
		os.Exit(0)
	}
	if *versionFlag {
		fmt.Printf("my app %s, commit %s, built at %s", version, commit, date)
		os.Exit(0)
	}
	if *namePtr == "" {
		fmt.Fprintln(os.Stderr, "Error: the --name flag is required.")
		flag.Usage()
		os.Exit(1)
	}

	p := tea.NewProgram(initialModel(*namePtr))

	result, err := p.Run()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error running TUI: %v\n", err)
		os.Exit(1)
	}
	if final, ok := result.(model); ok && len(final.errors) > 0 {
		os.Exit(1)
	}
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case actionFinishedMsg:
		if m.step != stepRun || len(m.queue) == 0 || m.queue[0] != msg.index {
			return m, nil
		}
		if msg.err != nil {
			m.status[msg.index] = statusFailed
			m.errors[msg.index] = msg.err
		} else {
			m.status[msg.index] = statusSuccess
		}

		m.queue = m.queue[1:]

		if len(m.queue) == 0 {
			m.step = stepDone
			return m, nil
		}

		return m, m.startNextAction()

	case tea.KeyPressMsg:
		if msg.String() == "ctrl+c" {
			return m, tea.Quit
		}
		if m.step == stepConfigure {
			return m.updateConfig(msg), nil
		}
		switch msg.String() {
		case "q":
			return m, tea.Quit

		case "up", "k":
			if m.step == stepSelect && m.cursor > 0 {
				m.cursor--
			}

		case "down", "j":
			if m.step == stepSelect && m.cursor < len(m.choices)-1 {
				m.cursor++
			}

		case "space":
			if m.step == stepSelect {
				if _, ok := m.selected[m.cursor]; ok {
					delete(m.selected, m.cursor)
				} else {
					m.selected[m.cursor] = struct{}{}
				}
			}

		case "enter":
			switch m.step {
			case stepSelect:
				if _, ok := m.selected[5]; ok {
					m.step = stepConfigure
				} else {
					m.step = stepReview
				}
			case stepReview:
				m.step = stepRun
				return m, m.prepareActions()
			case stepDone:
				return m, tea.Quit
			}
		}
	}

	return m, nil
}

func (m model) updateConfig(msg tea.KeyPressMsg) model {
	switch msg.String() {
	case "esc":
		m.step = stepSelect
	case "up", "tab", "down", "shift+tab":
		m.configField = 1 - m.configField
	case "enter":
		if m.config.GitHub.Name == "" {
			m.configError = "Repository name cannot be empty"
		} else {
			m.configError = ""
			m.step = stepReview
		}
	case "space":
		if m.configField == 1 {
			m.config.GitHub.Public = !m.config.GitHub.Public
		}
	case "backspace":
		if m.configField == 0 && len(m.config.GitHub.Name) > 0 {
			name := []rune(m.config.GitHub.Name)
			m.config.GitHub.Name = string(name[:len(name)-1])
			m.configError = ""
		}
	default:
		if m.configField == 0 && msg.Text != "" {
			for _, r := range msg.Text {
				if r >= 'a' && r <= 'z' || r >= 'A' && r <= 'Z' || r >= '0' && r <= '9' || r == '-' || r == '_' || r == '.' {
					m.config.GitHub.Name += string(r)
				} else {
					m.configError = "Use letters, numbers, hyphens, underscores or dots"
				}
			}
		}
	}
	return m
}

func statusLabel(s actionStatus) string {
	switch s {
	case statusPending:
		return "Pending"
	case statusRunning:
		return "Running..."
	case statusSuccess:
		return "Done"
	case statusFailed:
		return "FAILED"
	default:
		return "Unknown"
	}
}

func (m model) View() tea.View {
	var s strings.Builder

	s.WriteString("\n  kicad-init\n")
	s.WriteString("  Initialize a new project\n")
	fmt.Fprintf(&s, "  Project: %s\n\n", m.projectName)

	switch m.step {
	case stepSelect:
		s.WriteString("  Select what to include:\n\n")

		for i, choice := range m.choices {
			cursor := "  "
			if m.cursor == i {
				cursor = "❯ "
			}

			checked := " "
			if _, ok := m.selected[i]; ok {
				checked = "x"
			}

			fmt.Fprintf(&s, "%s[%s] %s\n", cursor, checked, choice)
		}

		s.WriteString("\n  ↑/↓ Navigate • Space Toggle\n")
		s.WriteString("  Enter Continue • q Quit\n")

	case stepConfigure:
		s.WriteString("  Configure GitHub repository:\n\n")
		nameCursor, visibilityCursor := "  ", "  "
		if m.configField == 0 {
			nameCursor = "❯ "
		} else {
			visibilityCursor = "❯ "
		}
		visibility := "Private"
		if m.config.GitHub.Public {
			visibility = "Public"
		}
		fmt.Fprintf(&s, "%sName: %s\n%sVisibility: %s\n", nameCursor, m.config.GitHub.Name, visibilityCursor, visibility)
		if m.configError != "" {
			fmt.Fprintf(&s, "\n  %s\n", m.configError)
		}
		s.WriteString("\n  ↑/↓ Select • Type/Backspace Edit Name • Space Toggle Visibility\n")
		s.WriteString("  Enter Review • Esc Back • Ctrl+C Quit\n")

	case stepReview:
		s.WriteString("  Review selected actions:\n\n")

		if len(m.selected) == 0 {
			s.WriteString("  No actions selected.\n")
		}

		for i, choice := range m.choices {
			if _, ok := m.selected[i]; ok {
				fmt.Fprintf(&s, "  • %s\n", choice)
				if i == 5 {
					visibility := "private"
					if m.config.GitHub.Public {
						visibility = "public"
					}
					fmt.Fprintf(&s, "    Name: %s (%s)\n", m.config.GitHub.Name, visibility)
				}
			}
		}

		s.WriteString("\n  Enter Initialize • q Quit\n")

	case stepRun, stepDone:
		s.WriteString("  Initialization progress:\n\n")

		for i, choice := range m.choices {
			if _, ok := m.selected[i]; !ok {
				continue
			}

			fmt.Fprintf(&s, "  %-24s %s\n", choice, statusLabel(m.status[i]))

			if err := m.errors[i]; err != nil {
				fmt.Fprintf(&s, "    Error: %v\n", err)
			}
		}

		if m.step == stepRun {
			s.WriteString("\n  Working...\n")
		} else {
			s.WriteString("\n  Initialization finished.\n")
			s.WriteString("  Press Enter to exit.\n")
		}
	}

	return tea.NewView(s.String())
}
