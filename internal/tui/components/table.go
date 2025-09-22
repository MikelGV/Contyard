package components

import (
	"fmt"

	"github.com/MikelGV/Contyard/internal/data/types"
	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/lipgloss/table"
)
var ()

/**
    I need to get the data as a [][]string to be able to pass it to the row and 
    a StyleFunc
**/
func CreateTable(d []types.ContainerStats) *table.Table {

    t := table.New().
        Border(lipgloss.NormalBorder()).
        BorderStyle(lipgloss.
            NewStyle().
            Foreground(lipgloss.Color("255"))).
        StyleFunc(func(row, col int) lipgloss.Style {
            switch {
            case row == table.HeaderRow:
                return  lipgloss.NewStyle().
                    Foreground(lipgloss.
                        Color("255")).
                    Bold(true).
                    Align(lipgloss.Center)
            case row%2 == 0:
                return  lipgloss.NewStyle().
                    Foreground(lipgloss.
                        Color("255")).
                    Bold(true).
                    Align(lipgloss.Center)
            default:
            return  lipgloss.NewStyle().
                Foreground(lipgloss.
                    Color("#AAAAA")).
                Bold(true).
                Align(lipgloss.Center)
            }
        }).
        Headers("ID", "Name", "CPU Usage", "Memory Usage", "Memory Limit")

    for _, stats := range d {
        t.Row(
            stats.ID,
            truncate(stats.NAME, 20),
            fmt.Sprintf("%d ns", stats.CPUUSAGE),
            format(stats.MEMORYUSAGE),
            format(stats.MEMORYLIMIT),
        )
    }

    return t  
}

func truncate(s string, w int) string {
    if len(s) <= w {
        return s
    }

    return s[:w-3] + "..."
}

func format(u uint64) string {
    const w = 1024

    if u < w {
        return fmt.Sprintf("%d B", u)
    }

    div, exp := uint64(w), 0
    for n := u; n >= w; n/= w {
        div *= w
        exp++
    }

    return fmt.Sprintf("%.1f %cB", float64(u)/float64(div), "KMGTPE"[exp])
}
