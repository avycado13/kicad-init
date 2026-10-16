package main

import (
	"flag"
	"fmt"
	"os"
	"os/exec"
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
		choices:     choices,
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
		return actions.CreateGithubRepo(m.projectName)
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
		m.step = 3
		return nil
	}

	mgr, err := files.NewManager(m.rootDir)
	if err != nil {
		for _, index := range m.queue {
			m.status[index] = statusFailed
			m.errors[index] = err
		}
		m.step = 3
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
		if m.step != 2 || len(m.queue) == 0 || m.queue[0] != msg.index {
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
			m.step = 3
			return m, nil
		}

		return m, m.startNextAction()

	case tea.KeyPressMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit

		case "up", "k":
			if m.step == 0 && m.cursor > 0 {
				m.cursor--
			}

		case "down", "j":
			if m.step == 0 && m.cursor < len(m.choices)-1 {
				m.cursor++
			}

		case "space":
			if m.step == 0 {
				if _, ok := m.selected[m.cursor]; ok {
					delete(m.selected, m.cursor)
				} else {
					m.selected[m.cursor] = struct{}{}
				}
			}

		case "enter":
			switch m.step {
			case 0:
				m.step = 1
			case 1:
				m.step = 2
				return m, m.prepareActions()
			case 3:
				return m, tea.Quit
			}
		}
	}

	return m, nil
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
	case 0:
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
		s.WriteString("  Enter Review • q Quit\n")

	case 1:
		s.WriteString("  Review selected actions:\n\n")

		if len(m.selected) == 0 {
			s.WriteString("  No actions selected.\n")
		}

		for i, choice := range m.choices {
			if _, ok := m.selected[i]; ok {
				fmt.Fprintf(&s, "  • %s\n", choice)
			}
		}

		s.WriteString("\n  Enter Initialize • q Quit\n")

	case 2, 3:
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

		if m.step == 2 {
			s.WriteString("\n  Working...\n")
		} else {
			s.WriteString("\n  Initialization finished.\n")
			s.WriteString("  Press Enter to exit.\n")
		}
	}

	return tea.NewView(s.String())
}
