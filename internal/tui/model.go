package tui

import (
    tea "github.com/charmbracelet/bubbletea"
)
/**
    In the model we store the state
**/
type model struct {
    choices []string
    cursor int
    selected map[int]struct{}
}

/**
    Here we define our initial state
**/
func initialModel() model {
    return model{
        choices: []string{"buy carrots", "buy celery", "buy kohlrabi"},

        // a map which indicates which choices are selected.
        selected: make(map[int]struct{}),
    }
}

/**
    Initial I/O
**/
func (m model) Init() tea.Cmd {
    return nil
}
