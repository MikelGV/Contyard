package tui

import (
	"github.com/MikelGV/Contyard/internal/data/types"
	"github.com/MikelGV/Contyard/internal/tui/components"
	tea "github.com/charmbracelet/bubbletea"
)

/**
    We handle things when things happen.
    In this case we need to handle when we get stats and when we get errors
**/
func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
    switch msg := msg.(type) {
    
    case []types.ContainerStats:
        m.stats = msg
        m.err = nil
        m.table = components.CreateTable(m.stats)
        return m, waitForStats(m) 

    case error:
        m.err = msg
        return m, waitForStats(m) 

    case tea.KeyMsg:
        switch msg.String() {
        case "q", "ctrl+c":
            m.cancel()
            return m, tea.Quit
        }
    }
    
    return m, nil
}

func waitForStats(m Model) tea.Cmd {
    return func() tea.Msg {
        select {
        case stats, ok := <-m.statsChan:
            if !ok {
                return nil
            }
            return stats
        case err, ok := <-m.errChan:
            if !ok {
                return nil
            }
            return  err
        }
    }
}
