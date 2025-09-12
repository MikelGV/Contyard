package docker

import (
	"context"
	"fmt"
	"io"
	"strings"
	"testing"

	data "github.com/MikelGV/Contyard/internal/data/docker"
	"github.com/moby/moby/api/types/container"
	"github.com/moby/moby/client"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockDockerClient struct {
    mock.Mock
}

func (m *MockDockerClient) ContainerList(ctx context.Context, options container.ListOptions) ([]container.Summary, error) {
    args := m.Called(ctx, options)
    return args.Get(0).([]container.Summary), args.Error(1)
}
func (m *MockDockerClient) ContainerStats(ctx context.Context, containerID string, stream bool) (client.StatsResponseReader, error) {
    args := m.Called(ctx, containerID, stream)
    return args.Get(0).(client.StatsResponseReader), args.Error(1)
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

        stats, err := data.GetDockerConStats(ctx, start)
        assert.NoError(t, err)
        assert.Empty(t, stats)
        mockClient.AssertExpectations(t)
    })

    t.Run("Single container with valid stats", func(t *testing.T) {
        mockClient := &MockDockerClient{}
        containers := []container.Summary{
            {
                ID: "12345678901234abcd",
                Names: []string{"/test-container"},
            },
        }

        statsJSON := `{
            "cpu_stats": {"cpu_usage": {"total_usage": 1000}},
            "memory_stats": {"usage": 1024, "limit": 2048}
        }`

        mockStats := client.StatsResponseReader{
            Body: io.NopCloser(strings.NewReader(statsJSON)),
        }

        start := func() (data.DockerClient, error) {
            mockClient.On("ContainerList", ctx, container.ListOptions{}).Return(containers, nil)
            mockClient.On("ContainerStats", ctx, "12345678901234abcd", false).Return(mockStats, nil)
            mockClient.On("Close").Return(nil)
            return mockClient, nil
        }

        stats, err := data.GetDockerConStats(ctx, start)
        assert.NoError(t, err)
        assert.Len(t, stats, 1)
        assert.Equal(t, "123456789012", stats[0].ID)
        assert.Equal(t, "/test-container", stats[0].NAME)
        assert.Equal(t, uint64(1000), stats[0].CPUUSAGE)
        assert.Equal(t, uint64(1024), stats[0].MEMORYUSAGE)
        assert.Equal(t, uint64(2048), stats[0].MEMORYLIMIT)
        mockClient.AssertExpectations(t)
    })

    t.Run("multiple containers with valid stats", func(t *testing.T) {
        mockClient := &MockDockerClient{}
        containers := []container.Summary{
            {
                ID: "12345678901234abcd",
                Names: []string{"/test-container1"},
            },
            {
                ID: "123456789abcd01234",
                Names: []string{"/test-container2"},
            },
        }

        statsJSON1 := `{
            "cpu_stats": {"cpu_usage": {"total_usage": 1000}},
            "memory_stats": {"usage": 1024, "limit": 2048}
        }`

        statsJSON2 := `{
            "cpu_stats": {"cpu_usage": {"total_usage": 2000}},
            "memory_stats": {"usage": 2028, "limit": 4096}
        }`

        mockStats1 := client.StatsResponseReader{
            Body: io.NopCloser(strings.NewReader(statsJSON1)),
        }

        mockStats2 := client.StatsResponseReader{
            Body: io.NopCloser(strings.NewReader(statsJSON2)),
        }

        start := func() (data.DockerClient, error) {
            mockClient.On("ContainerList", ctx, container.ListOptions{}).Return(containers, nil)
            mockClient.On("ContainerStats", ctx, "12345678901234abcd", false).Return(mockStats1, nil)
            mockClient.On("ContainerStats", ctx, "123456789abcd01234", false).Return(mockStats2, nil)
            mockClient.On("Close").Return(nil)
            return mockClient, nil
        }

        stats, err := data.GetDockerConStats(ctx, start)
        assert.NoError(t, err)
        assert.Len(t, stats, 2)
        assert.Equal(t, "123456789012", stats[0].ID)
        assert.Equal(t, "/test-container1", stats[0].NAME)
        assert.Equal(t, uint64(1000), stats[0].CPUUSAGE)
        assert.Equal(t, uint64(1024), stats[0].MEMORYUSAGE)
        assert.Equal(t, uint64(2048), stats[0].MEMORYLIMIT)
        assert.Equal(t, "123456789abc", stats[1].ID)
        assert.Equal(t, "/test-container2", stats[1].NAME)
        assert.Equal(t, uint64(2000), stats[1].CPUUSAGE)
        assert.Equal(t, uint64(2028), stats[1].MEMORYUSAGE)
        assert.Equal(t, uint64(4096), stats[1].MEMORYLIMIT)
        mockClient.AssertExpectations(t)
    })

    t.Run("multiple containers with invalid stats json", func(t *testing.T) {
        mockClient := &MockDockerClient{}
        containers := []container.Summary{
            {
                ID: "12345678901234abcd",
                Names: []string{"/test-container1"},
            },
        }

        mockStats1 := client.StatsResponseReader{
            Body: io.NopCloser(strings.NewReader("invalid json")),
        }

        start := func() (data.DockerClient, error) {
            mockClient.On("ContainerList", ctx, container.ListOptions{}).Return(containers, nil)
            mockClient.On("ContainerStats", ctx, "12345678901234abcd", false).Return(mockStats1, nil)
            mockClient.On("Close").Return(nil)
            return mockClient, nil
        }

        stats, err := data.GetDockerConStats(ctx, start)
        assert.NoError(t, err)
        assert.Empty(t, stats)
        mockClient.AssertExpectations(t)
    })

    t.Run("containers with no name", func(t *testing.T) {
        mockClient := &MockDockerClient{}
        containers := []container.Summary{
            {
                ID: "12345678901234abcd",
                Names: []string{},
            },
        }

        statsJSON := `{
            "cpu_stats": {"cpu_usage": {"total_usage": 1000}},
            "memory_stats": {"usage": 1024, "limit": 2048}
        }`

        mockStats := client.StatsResponseReader{
            Body: io.NopCloser(strings.NewReader(statsJSON)),
        }

        start := func() (data.DockerClient, error) {
            mockClient.On("ContainerList", ctx, container.ListOptions{}).Return(containers, nil)
            mockClient.On("ContainerStats", ctx, "12345678901234abcd", false).Return(mockStats, nil)
            mockClient.On("Close").Return(nil)
            return mockClient, nil
        }

        stats, err := data.GetDockerConStats(ctx, start)
        assert.NoError(t, err)
        assert.Len(t, stats, 1)
        assert.Equal(t, "123456789012", stats[0].ID)
        assert.Equal(t, "", stats[0].NAME)
        assert.Equal(t, uint64(1000), stats[0].CPUUSAGE)
        assert.Equal(t, uint64(1024), stats[0].MEMORYUSAGE)
        assert.Equal(t, uint64(2048), stats[0].MEMORYLIMIT)
        mockClient.AssertExpectations(t)
    })

    t.Run("client creation error", func(t *testing.T) {
        start := func() (data.DockerClient, error) {
            return nil, fmt.Errorf("client creation failed")
        }

        stats, err := data.GetDockerConStats(ctx, start)
        assert.Error(t, err)
        assert.Contains(t, err.Error(), "Client setup failed:")
        assert.Nil(t, stats)
    })

    t.Run("containerList error", func(t *testing.T) {
        mockClient := &MockDockerClient{}

        start := func() (data.DockerClient, error) {
            mockClient.On("ContainerList", ctx, container.ListOptions{}).Return([]container.Summary(nil), fmt.Errorf("docker error"))
            mockClient.On("Close").Return(nil)
            return mockClient, nil
        }

        stats, err := data.GetDockerConStats(ctx, start)
        assert.Error(t, err)
        assert.Contains(t, err.Error(), "failed to get all containers")
        assert.Nil(t, stats)
        mockClient.AssertExpectations(t)
    })

    t.Run("containerStats error", func(t *testing.T) {
        mockClient := &MockDockerClient{}
        containers := []container.Summary{
            {
                ID: "12345678901234abcd",
                Names: []string{"/test-container"},
            },
        }

        start := func() (data.DockerClient, error) {
            mockClient.On("ContainerList", ctx, container.ListOptions{}).Return(containers, nil)
            mockClient.On("ContainerStats", ctx, "12345678901234abcd", false).Return(client.StatsResponseReader{}, fmt.Errorf("stats error"))
            mockClient.On("Close").Return(nil)
            return mockClient, nil
        }

        stats, err := data.GetDockerConStats(ctx, start)
        assert.Error(t, err)
        assert.Contains(t, err.Error(), "failed to get containers stats")
        assert.Nil(t, stats)
        mockClient.AssertExpectations(t)
    })
    /**
    **/
}
