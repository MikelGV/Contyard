package tui

import (
	"github.com/MikelGV/Contyard/internal/tui/components"
	"github.com/charmbracelet/lipgloss"
)

/**
    This is all of the cool "FrontEnd" or what the users will see :D
**/
var (
    errorStyle = lipgloss.NewStyle().
        Foreground(lipgloss.Color("#DB222A")).
        Padding(0, 1)
)


func (m Model) View() string {
    var output []string
    if !m.docker {
        output = append(output, errorStyle.Render("No data sources selected (use -d)"))
    } else {
        table := components.CreateTable(m.stats)
        output = append(output, table)
    }

    return lipgloss.JoinVertical(lipgloss.Left, output...) 
}

