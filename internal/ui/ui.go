package ui

import (
	"fmt"
	"io"
	"strings"

	"github.com/charmbracelet/lipgloss"
)

type UI struct {
	out        io.Writer
	err        io.Writer
	styled     bool
	headerStyle lipgloss.Style
	headerText  lipgloss.Style
}

func New(out, err io.Writer, styled bool) *UI {
	headerStyle := lipgloss.NewStyle().
		BorderStyle(lipgloss.DoubleBorder()).
		BorderForeground(lipgloss.Color("212")).
		Align(lipgloss.Center).
		Width(50).
		Padding(1, 4)

	headerText := lipgloss.NewStyle().Bold(true)

	return &UI{
		out:         out,
		err:         err,
		styled:      styled,
		headerStyle: headerStyle,
		headerText:  headerText,
	}
}

func (u *UI) PrintHeader(name, version string) {
	if u.styled {
		content := fmt.Sprintf("%s\n%s", name, "v"+version)
		rendered := u.headerStyle.Render(u.headerText.Render(content))
		fmt.Fprintln(u.out)
		fmt.Fprintln(u.out, rendered)
		return
	}

	border := "╔════════════════════════════════════════════════╗"
	fmt.Fprintln(u.out)
	fmt.Fprintln(u.out, border)
	fmt.Fprintf(u.out, "║%s║\n", padCenter(name, 48))
	fmt.Fprintf(u.out, "║%s║\n", padCenter("v"+version, 48))
	fmt.Fprintln(u.out, "╚════════════════════════════════════════════════╝")
	fmt.Fprintln(u.out)
}

func (u *UI) PrintHelp(help string) {
	fmt.Fprintln(u.out, help)
}

func (u *UI) PrintError(msg string) {
	fmt.Fprintln(u.err, msg)
}

func padCenter(text string, width int) string {
	if width <= len(text) {
		return text
	}
	pad := width - len(text)
	left := pad / 2
	right := pad - left
	return strings.Repeat(" ", left) + text + strings.Repeat(" ", right)
}
