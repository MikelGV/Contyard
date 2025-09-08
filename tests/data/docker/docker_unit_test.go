package docker

import (
	"context"
	"testing"

	data "github.com/MikelGV/Contyard/internal/data/docker"
	"github.com/moby/moby/api/types/container"
	"github.com/stretchr/testify/mock"
)

type MockDockerClient struct {
    mock.Mock
}

func (m *MockDockerClient) ContainerList(ctx context.Context, options container.ListOptions) ([]container.Summary, error) {
    args := m.Called(ctx, options)
    return args.Get(0).([]container.Summary), args.Error(1)
}
func (m *MockDockerClient) ContainerStats(ctx context.Context, containerID string, stream bool) (container.StatsResponse, error) {
    args := m.Called(ctx, containerID, stream)
    return args.Get(0).(container.StatsResponse), args.Error(1)
}

func (m *MockDockerClient) Close() error {
    args := m.Called()
    return args.Error(0)
}
/**
    This test is for the GetDockerConStats function, in this function we want to
    retrieve docker containers stat values if there are any containers running
**/
func TestDockerConStats(t *testing.T) {
    ctx := context.Background()
    /**
        I need to initialize an http server to mock moby's docker api connection 
        and then build all of the checks. This is what it should check: 
            1. Grab container and check if there are any containers,
            2. Can we stract the stats from a container.
            3. Can we decode the container stats.
            4. Can we format the stats into legible data. 
    **/
    t.Run("No containers", func(t *testing.T) {
        mockClient := &MockDockerClient{}
        start := func() (data.DockerClient, error) {
            mockClient.On("ContainerList", ctx, container.ListOptions{}).Return([]container.Summary{}, nil)
            mockClient.On("Close").Return(nil)
            return mockClient, nil
        }
    })
}

/**
    This test is for the StreamDockerConStats function, in this function we want
    to  return live stat updates 
**/
func StreamDockerTest() {
}
