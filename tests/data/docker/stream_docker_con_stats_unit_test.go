package docker_test

import (
	"context"
	"fmt"
	"io"
	"strings"
	"testing"
	"time"

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
    This test is for the StreamDockerConStats function, in this function we want
    to  return live stat updates
**/
func TestStreamDocker(t *testing.T) {
    ctx := context.Background()

    t.Run("Streams multiple containers stats", func(t *testing.T) {
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
            mockClient.On("ContainerList", mock.Anything, container.ListOptions{}).Return(containers, nil).Once()
            mockClient.On("ContainerStats", mock.Anything, "12345678901234abcd", false).Return(mockStats1, nil).Once()
            mockClient.On("ContainerStats", mock.Anything, "123456789abcd01234", false).Return(mockStats2, nil).Once()
            mockClient.On("Close").Return(nil).Once()
            return mockClient, nil
        }

        ctx, cancel := context.WithCancel(ctx)
        defer cancel()

        statsChan, errChan := data.StreamDockerConStats(ctx, start, 10*time.Millisecond)

        select {
        case stats := <-statsChan:
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

        case err := <-errChan:
            t.Fatalf("Unexpected error: %v", err)
        
        case <-time.After(100 * time.Millisecond):
            t.Fatalf("Time out waiting for stats")
        }

        mockClient.AssertExpectations(t)
    }) 

    t.Run("Streams empty stats for no containers", func(t *testing.T) {
        mockClient := &MockDockerClient{}

        start := func() (data.DockerClient, error) {
            mockClient.On("ContainerList", mock.Anything, container.ListOptions{}).Return([]container.Summary{}, nil).Once()
            mockClient.On("Close").Return(nil).Once()
            return mockClient, nil
        }

        ctx, cancel := context.WithCancel(ctx)
        defer cancel()

        statsChan, errChan := data.StreamDockerConStats(ctx, start, 10*time.Millisecond)

        select {
        case stats := <-statsChan:
            assert.Empty(t, stats)

        case err := <-errChan:
            t.Fatalf("Unexpected error: %v", err)

        case <-time.After(100 * time.Millisecond):
            t.Fatalf("Time out waiting for stats")
        }

        mockClient.AssertExpectations(t)
    })

    t.Run("Handle GetDockerConStats error", func(t *testing.T) {
        mockClient := &MockDockerClient{}

        start := func() (data.DockerClient, error) {
            mockClient.On("ContainerList", mock.Anything, container.ListOptions{}).Return([]container.Summary(nil), fmt.Errorf("docker error")).Once()
            mockClient.On("Close").Return(nil).Once()
            return mockClient, nil
        }

        ctx, cancel := context.WithCancel(context.Background())
        defer cancel()

        statsChan, errChan := data.StreamDockerConStats(ctx, start, 10*time.Millisecond)

        select {
        case err := <-errChan:
            assert.Contains(t, err.Error(), "failed to get all containers:")
        case stats := <-statsChan:
            t.Fatalf("Unexpected stats: %v", stats)
        case <-time.After(1 * time.Second):
            t.Fatal("Time out waiting for error")
        }

        mockClient.AssertExpectations(t)
    })

    /**
        t.Run("Stops streaming on context cancellation", func(t *testing.T) {})
    **/


}
