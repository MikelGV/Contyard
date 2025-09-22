package tui

import (
	"time"

	"github.com/charmbracelet/lipgloss"
)

/**
    This is all of the cool "FrontEnd" or what the users will see :D
**/
var (
    errorStyle = lipgloss.NewStyle().
        Foreground(lipgloss.Color("#DB222A")).
        Padding(0, 1)
    warningStyle = lipgloss.NewStyle().
        Foreground(lipgloss.Color("#AAAAA")).
        Padding(0, 1)
)


func (m Model) View() string {
    var output []string
    currentTime := time.Now()
    if !m.docker {
        output = append(output, errorStyle.Render("No data sources selected (use -d)"))
    } else if m.err != nil {
        output = append(output, errorStyle.Render("Error:" + m.err.Error()))
    } else if len(m.stats) == 0 {
        if time.Since(currentTime) > 100 * time.Millisecond {
            output = append(output, errorStyle.Render("No containers found"))
        } else {
            output = append(output, warningStyle.Render("Waiting for containers"))
        }
    } else {
        output = append(output, m.table.Render())
    }

    return lipgloss.JoinVertical(lipgloss.Left, output...) 
}

