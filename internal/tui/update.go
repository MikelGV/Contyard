package tui

import (
	"github.com/MikelGV/Contyard/internal/data/types"
	tea "github.com/charmbracelet/bubbletea"
)

/**
    We handle things when things happen.
    In this case we need to handle when we get stats and when we get errors
**/
func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
    switch msg := msg.(type) {
    
    case []types.ContainerStats:
        /**
            We have to handle the stats so we have to compare the model stats
            with the msg adn check that the error is nil
        **/
        m.stats = msg
        m.err = nil
        return m, waitForStats(m) 

    case error:
        /**
            We have to handle if there is an error  
        **/
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
