package tui

import (
	"context"
	"time"

	data "github.com/MikelGV/Contyard/internal/data/docker"
	"github.com/MikelGV/Contyard/internal/data/types"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss/table"
)

/**
    In the model we store the state
**/
type Model struct {
    all bool
    docker bool
    podman bool
    kubernetes bool
    statsChan chan []types.ContainerStats
    errChan chan error
    err error
    stats []types.ContainerStats
    table *table.Table
    cancel context.CancelFunc
}

/**
    Here we define our initial state
**/
func NewModel(docker, all, podman, kubernetes bool) Model {
    ctx, cancel := context.WithCancel(context.Background())
    var statsChan chan []types.ContainerStats
    var errChan chan error
    if docker {
        statsChan, errChan = data.StreamDockerConStats(ctx, data.DefaultClientStart, 100*time.Millisecond)
    }
    // add podman, kubernetes, and all later
    return Model{
        docker: docker,
        podman: podman,
        kubernetes: kubernetes,
        all: all,
        statsChan: statsChan,
        errChan: errChan,
        table: table.New(),
        cancel: cancel,
    }
}

/**
    Initial I/O
**/
func (m Model) Init() tea.Cmd {
    return tea.Batch(waitForStats(m), tea.EnterAltScreen)
}
