package logging

import (
	"io"

	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/log"
)

type Logger struct {
	out    *log.Logger
	err    *log.Logger
	styled bool
	silent bool
	styles map[level]lipgloss.Style
}

type level int

const (
	levelInfo level = iota
	levelSuccess
	levelWarn
	levelError
)

func New(out, err io.Writer, styled bool) *Logger {
	outLogger := log.NewWithOptions(out, log.Options{ReportTimestamp: false})
	errLogger := log.NewWithOptions(err, log.Options{ReportTimestamp: false})
	outLogger.SetFormatter(log.TextFormatter)
	errLogger.SetFormatter(log.TextFormatter)
	outLogger.SetStyles(minimalStyles())
	errLogger.SetStyles(minimalStyles())

	styles := map[level]lipgloss.Style{
		levelInfo:    lipgloss.NewStyle().Foreground(lipgloss.Color("4")),
		levelSuccess: lipgloss.NewStyle().Foreground(lipgloss.Color("2")),
		levelWarn:    lipgloss.NewStyle().Foreground(lipgloss.Color("3")),
		levelError:   lipgloss.NewStyle().Foreground(lipgloss.Color("1")),
	}

	return &Logger{
		out:    outLogger,
		err:    errLogger,
		styled: styled,
		styles: styles,
	}
}

func minimalStyles() *log.Styles {
	styles := log.DefaultStyles()
	styles.Timestamp = lipgloss.NewStyle()
	styles.Caller = lipgloss.NewStyle()
	styles.Prefix = lipgloss.NewStyle()
	styles.Message = lipgloss.NewStyle()
	styles.Key = lipgloss.NewStyle()
	styles.Value = lipgloss.NewStyle()
	styles.Separator = lipgloss.NewStyle()
	styles.Levels = map[log.Level]lipgloss.Style{}
	return styles
}

func (l *Logger) SetSilent(silent bool) {
	l.silent = silent
}

func (l *Logger) SetVerbose(verbose bool) {
	if !verbose {
		return
	}
	l.out.SetStyles(log.DefaultStyles())
	l.err.SetStyles(log.DefaultStyles())
}

func (l *Logger) Info(msg string) {
	l.print(l.out, levelInfo, "", msg)
}

func (l *Logger) Success(msg string) {
	l.print(l.out, levelSuccess, "✓", msg)
}

func (l *Logger) Warn(msg string) {
	l.print(l.out, levelWarn, "⚠", msg)
}

func (l *Logger) Error(msg string) {
	l.print(l.err, levelError, "✗", msg)
}

func (l *Logger) print(logger *log.Logger, lvl level, prefix, msg string) {
	if l.silent {
		return
	}
	text := msg
	if prefix != "" {
		text = prefix + " " + msg
	}
	if l.styled {
		text = l.styles[lvl].Render(text)
	}
	logger.Print(text)
}
