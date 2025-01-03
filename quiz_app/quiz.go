//todo: file and time prompt
//todo: styling issue on alignment
// init > model > update > view > helper func

package main

import (
	"encoding/csv"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type model struct {
	questions    []question
	currentIndex int
	correctCount int
	textInput    textinput.Model
	spinner      spinner.Model
	timeLimit    time.Duration
	quizEnded    bool 
	err          error
}

type question struct {
	Question string
	Answer   string
}

func main() {
	m, err := initialModel()
	if err != nil {
		log.Fatal(err)
	}

	p := tea.NewProgram(m)
	quitMsg := time.After(3 * time.Second)
	go func() {
		<-quitMsg
		m.quizEnded = true 
		p.Quit()
	}()

	_, err = p.Run()
	if err != nil {
		log.Fatal(err)
	}
}

func initialModel() (model, error) {
	filePath := "problems.csv"
	timeLimit := 30 * time.Second

	questions, err := readCSV(filePath)
	if err != nil {
		return model{}, err
	}

	ti := textinput.New()
	ti.Placeholder = "Enter your answer"
	ti.Focus()
	ti.CharLimit = 100
	ti.Width = 30

	s := spinner.New()
	s.Spinner = spinner.Dot
	s.Style = lipgloss.NewStyle().Foreground(lipgloss.Color("205"))

	return model{
		questions:    questions,
		currentIndex: 0,
		correctCount: 0,
		textInput:    ti,
		spinner:      s,
		timeLimit:    timeLimit,
		err:          nil,
	}, nil
}

func (m model) Init() tea.Cmd {
	return m.spinner.Tick
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "esc":
			return m, tea.Quit
		case "enter":
			if m.textInput.Value() == m.questions[m.currentIndex].Answer {
				m.correctCount++
			}
			m.currentIndex++
			if m.currentIndex >= len(m.questions) {
				return m, tea.Quit
			}
			m.textInput.Reset()
			return m, nil
		default:
			var cmd tea.Cmd
			m.textInput, cmd = m.textInput.Update(msg)
			return m, cmd
		}

	default:
		var cmd tea.Cmd
		m.spinner, cmd = m.spinner.Update(msg)
		return m, cmd
	}
}

// todo: view model not showing ops messages 
func (m model) View() string {
	if m.quizEnded {
		return "Time's up!"
	}
	
	if len(m.questions) == 0 {
		return "No questions found."
	}

	if m.currentIndex >= len(m.questions) {
		return fmt.Sprintf("Quiz completed! Your score: %d/%d", m.correctCount, len(m.questions))
	}

	question := m.questions[m.currentIndex]
	return fmt.Sprintf(
		"\n   %s\n\n%s\n\n   Correct: %d\n\n   (press Enter to answer)\n",
		question.Question,
		m.textInput.View(),
		m.correctCount,
	)
}

func readCSV(filename string) ([]question, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	reader := csv.NewReader(file)
	records, err := reader.ReadAll()
	if err != nil {
		return nil, err
	}

	questions := make([]question, len(records))
	for i, record := range records {
		questions[i] = question{
			Question: record[0],
			Answer:   record[1],
		}
	}

	return questions, nil
}

