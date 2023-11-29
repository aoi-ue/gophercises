package main

import (
	"fmt"
	"log"
	"os"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

// main func
func main() {
	p := tea.NewProgram(initialModel())

	if _, err := p.Run(); err != nil {
		log.Fatal(err)
		os.Exit(1)
	}

}

// model as state: question to prompt, correctcount to keep track, timer to end
type model struct {
	textInput textinput.Model
	// questions []question
	correctcount int
	timer        int
	err          error
}

// type question struct {
// 	Question string
// 	Answer string
// }

// initial state with bubbletea: to read textinput,
func initialModel() model {
	ti := textinput.New()
	ti.Placeholder = "Enter Answer Here"
	ti.Focus()
	ti.CharLimit = 156
	ti.Width = 20

	return model{
		textInput: ti,
		// questions []question
		correctcount: 0,
		// todo: have to take timer from cobra's flag
		timer: 30,
		err:   nil,
	}
}

// init actionable: blink and display
func (m model) Init() tea.Cmd {
	return textinput.Blink
}

// update: update the correctcount, readline to trigger view render
func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "esc":
			return m, tea.Quit
		case "enter", " ":
			// todo: evaluate if input is same as column[1]
			return m, textinput.Blink
		}
	}

	return m, cmd

}

// view: simple view of a quiz, render everytime for new readLine
func (m model) View() string {
	return fmt.Sprintf(
		"Question Goes here\n\n%s\n\n%s",
		m.textInput.View(),
		"(ctrl c or esc to quit)",
	)
}

// lipgloss: for styling

// readCSV: purely to read CSV
