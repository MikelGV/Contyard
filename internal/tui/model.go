package tui

import (
	"context"
	"time"

	data "github.com/MikelGV/Contyard/internal/data/docker"
	"github.com/MikelGV/Contyard/internal/data/types"
	tea "github.com/charmbracelet/bubbletea"
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
    cancel context.CancelFunc
}

/**
    Here we define our initial state
**/
func NewModel(docker, all, podman, kubernetes bool) Model {
    var statsChan chan []types.ContainerStats
    var errChan chan error
    ctx, cancel := context.WithCancel(context.Background())
    if docker {
        statsChan, errChan = data.StreamDockerConStats(ctx, data.DefaultClientStart, 2*time.Second)
    }
    // add podman, kubernetes, and all later
    return Model{
        docker: docker,
        podman: podman,
        kubernetes: kubernetes,
        all: all,
        statsChan: statsChan,
        errChan: errChan,
        cancel: cancel,
    }
}

/**
    Initial I/O
**/
func (m Model) Init() tea.Cmd {
    // i have to see how or if i need to return anything here so for now i will maintain it a nil
    return nil
}
