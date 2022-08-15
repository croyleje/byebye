package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var style = lipgloss.NewStyle().
	BorderStyle(lipgloss.RoundedBorder()).
	BorderForeground(lipgloss.Color("12")).
	Align(lipgloss.Center).
	Bold(true).
	Padding(2).
	Margin(1)

type model struct {
	cursor   int
	choices  []string
	selected map[int]struct{}
}

func end(choice string) {
	switch choice {
	case "Suspend":
		cmd := exec.Command("systemctl", "suspend")
		err := cmd.Run()

		if err != nil {
			log.Fatal(err)
		}
	case "Logout":
		cmd := exec.Command("kill", "-9", "-1")
		err := cmd.Run()

		if err != nil {
			log.Fatal(err)
		}
	case "Shutdown":
		cmd := exec.Command("shutdown", "now")
		err := cmd.Run()

		if err != nil {
			log.Fatal(err)
		}
	case "Restart":
		cmd := exec.Command("shutdown", "-r")
		err := cmd.Run()

		if err != nil {
			log.Fatal(err)
		}
	case "Hibernate":
		cmd := exec.Command("systemctl", "hibernate")
		err := cmd.Run()

		if err != nil {
			log.Fatal(err)
		}
	}
}

func initialModel() model {
	return model{
		choices: []string{"Suspend", "Logout", "Shutdown", "Restart", "Hibernate"},

		// A map which indicates which choices are selected. We're using
		// the  map like a mathematical set. The keys refer to the indexes
		// of the `choices` slice, above.
		selected: make(map[int]struct{}),
	}
}

func (m model) Init() tea.Cmd {
	return tea.Batch(nil, tea.EnterAltScreen)
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q", "esc":
			return m, tea.Quit
		case "up", "k":
			if m.cursor > 0 {
				m.cursor--
			}
		case "down", "j":
			if m.cursor < len(m.choices)-1 {
				m.cursor++
			}
		case "enter":
			end(m.choices[m.cursor])
			return m, tea.Quit
		}
	}

	return m, nil
}

func (m model) View() string {
	s := "\n\nWhere do you want to GO?\n\n"

	for i, choice := range m.choices {
		cursor := " "
		if m.cursor == i {
			cursor = "❯"
		}

		s += fmt.Sprintf("%s %s\n", cursor, choice)
	}

	s += "\nPress q to quit.\n\n\n"

	return style.Render(s)
}

func main() {
	p := tea.NewProgram(initialModel())
	if err := p.Start(); err != nil {
		fmt.Printf("Alas, there's been an error: %v", err)
		os.Exit(1)
	}
}
