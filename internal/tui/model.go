package tui

import (
    tea "github.com/charmbracelet/bubbletea"
)
/**
    In the model we store the state
**/
type model struct {
    stats []string
    err error
}

/**
    Here we define our initial state
**/
func InitialModel() model {
    return model{

    }
}

/**
    Initial I/O
**/
func (m model) Init() tea.Cmd {
    // i have to see how or if i need to return anything here so for now i will maintain it a nil
    return nil
}
