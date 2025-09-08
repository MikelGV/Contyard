package data

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"time"

	"github.com/MikelGV/Contyard/internal/data/types"
	"github.com/moby/moby/api/types/container"
	"github.com/moby/moby/api/types/network"
	"github.com/moby/moby/client"
	spec "github.com/opencontainers/image-spec/specs-go/v1"
)

type DockerClient interface {
    ContainerList(ctx context.Context, option container.ListOptions) ([]container.Summary, error)
    ContainerStats(ctx context.Context, containerID string, stream bool) (container.StatsResponse, error)
    Close() error
}

type RealDockerClient struct {
    *client.Client
}

func ClientStart() (*RealDockerClient, error) {
    cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
    if err != nil {
        return nil, err
    }

    return &RealDockerClient{cli}, nil
}

type DockerStats struct {
    CPUStats    struct {
        CPUUsage    struct {
            TotalUsage uint64 `json:"total_usage"`
        } `json:"cpu_usage"`
    } `json:"cpu_stats"`
    MemoryStats struct {
        Usage uint64 `json:"usage"`
        Limit uint64 `json:"limit"`
    } `json:"memory_stats"`
}

/**
    Docker will be our default when loocking containers health
**/
func GetDockerConStats(ctx context.Context, start ClientStart) ([]types.ContainerStats, error) {
    /**
        I have to add a break if no containers are found
    **/
    cli, err := start() 
    if err != nil {
        return nil, fmt.Errorf("Client setup failed: %w", err)
    }
    defer cli.Close()

    containers, err := cli.ContainerList(ctx, container.ListOptions{})
    if err != nil {
        return nil, fmt.Errorf("failed to get all containers: %w", err)
    }

    if len(containers) == 0 {
        return []types.ContainerStats{}, nil
    }

    var ContainerStats []types.ContainerStats 
    for _, ctr := range containers {
        stats, err := cli.ContainerStats(ctx, ctr.ID, false )
        if err != nil {
            return nil, fmt.Errorf("failed to get containers stats: %w", err)
        }

        dec := json.NewDecoder(stats.Body)
        var ds DockerStats
        if err :=dec.Decode(&ds); err != nil {
            stats.Body.Close()
            continue
        }

        stats.Body.Close()

        ContainerStats = append(ContainerStats, types.ContainerStats{
            ID: ctr.ID[:12],
            NAME: ctr.Names[0],
            CPUUSAGE: ds.CPUStats.CPUUsage.TotalUsage,
            MEMORYUSAGE: ds.MemoryStats.Usage,
            MEMORYLIMIT: ds.MemoryStats.Limit,
        })

    }

    return ContainerStats, nil
}

/**
    Returns a channel for real-time updates.
**/
func StreamDockerConStats(ctx context.Context, interval time.Duration) (chan([]types.ContainerStats), chan(error)) {
    statsChan := make(chan []types.ContainerStats)
    errChan := make(chan error)

    go func() {
        timer := time.NewTicker()
        defer timer.Stop()

        for {
            select {}
        }
    }

    return statsChan, nil
}

