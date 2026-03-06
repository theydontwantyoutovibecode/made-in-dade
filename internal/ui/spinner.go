package ui

import (
	"fmt"
	"io"
	"time"

	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
)

type Spinner struct {
	out     io.Writer
	enabled bool
}

type spinnerResult struct {
	err error
}

type spinnerModel struct {
	spinner spinner.Model
	message string
	work    func() error
	result  spinnerResult
}

func NewSpinner(out io.Writer, enabled bool) *Spinner {
	return &Spinner{out: out, enabled: enabled}
}

func (s *Spinner) Run(message string, work func() error) error {
	if !s.enabled {
		return runPlainSpinner(s.out, message, work)
	}

	model := spinnerModel{
		spinner: spinner.New(spinner.WithSpinner(spinner.Dot)),
		message: message,
		work:    work,
	}
	p := tea.NewProgram(model, tea.WithOutput(s.out))
	final, err := p.Run()
	if err != nil {
		return err
	}
	if res, ok := final.(spinnerModel); ok {
		return res.result.err
	}
	return nil
}

func (m spinnerModel) Init() tea.Cmd {
	return tea.Batch(m.spinner.Tick, func() tea.Msg {
		return spinnerResult{err: m.work()}
	})
}

func (m spinnerModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		return m, tea.Quit
	case spinnerResult:
		m.result = msg
		return m, tea.Quit
	}

	var cmd tea.Cmd
	m.spinner, cmd = m.spinner.Update(msg)
	return m, cmd
}

func (m spinnerModel) View() string {
	return fmt.Sprintf("%s %s", m.spinner.View(), m.message)
}

func runPlainSpinner(out io.Writer, message string, work func() error) error {
	fmt.Fprintf(out, "%s... ", message)
	start := time.Now()
	err := work()
	if err != nil {
		fmt.Fprintln(out, "failed")
		return err
	}
	if time.Since(start) < 10*time.Millisecond {
		fmt.Fprintln(out, "done")
	} else {
		fmt.Fprintln(out, "done")
	}
	return nil
}
