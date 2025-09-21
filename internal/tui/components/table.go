package components

import (
	"github.com/MikelGV/Contyard/internal/data/types"
	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/lipgloss/table"
)
var ()

func CreateTable(s []types.ContainerStats) string {
    t := table.New().
        Border(lipgloss.NormalBorder()).
        BorderStyle(lipgloss.
            NewStyle().
            Foreground(lipgloss.Color("255"))).
        /**
        StyleFunc(func(row, col int) lipgloss.Style {
            switch {
            case row == table.HeaderRow:
                return  lipgloss.NewStyle().
                    Foreground(lipgloss.
                        Color("99")).
                    Bold(true).
                    Align(lipgloss.Center)
            case row%2 == 0:
                return cell
            }
        }).
        **/
        Headers("containers", "pods").
        Row()

    return  t
}
