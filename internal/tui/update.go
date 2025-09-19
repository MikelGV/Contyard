package tui

import tea "github.com/charmbracelet/bubbletea"

/**
    We handle the message types here
**/
func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
    switch msg := msg.(type) {
    
    case statsMsg:
        return nil, nil

    case errMsg:
        return nil, nil

    case tea.KeyMsg:
        if msg.Type == tea.KeyCtrlC {
            return m, tea.Quit
        }
    }
    
    return m, nil
}
