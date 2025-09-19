package tui

import (
	"github.com/charmbracelet/lipgloss"
)

/**
    This is all of the cool "FrontEnd" or what the users will see :D
**/

var baseStyle = lipgloss.NewStyle().
    BorderStyle(lipgloss.NormalBorder()).
    BorderForeground(lipgloss.Color("250"))

func (m Model) View() string {
    if m.docker {
        return baseStyle.Render() 
    }
    return "No data to display"
}

