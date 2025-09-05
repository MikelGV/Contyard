package data

import (
	"context"
	"encoding/json"
	"fmt"
	"io"

	"github.com/moby/moby/api/types/container"
	"github.com/moby/moby/client"
)

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
func DockerConStats() error {
    /**
        I have to add a break if no containers are found
    **/
    cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
    if err != nil {
        return fmt.Errorf("Client setup failed: %w", err)
    }
    defer cli.Close()

    containers, err := cli.ContainerList(context.Background(), container.ListOptions{All: true})
    if err != nil {
        return fmt.Errorf("failed to get all containers: %w", err)
    }

    for _, ctr := range containers {
        stats, err := cli.ContainerStats(context.Background(), ctr.ID, false )
        if err != nil {
            return fmt.Errorf("failed to get containers stats: %w", err)
        }

        defer stats.Body.Close()

        body, err := io.ReadAll(stats.Body)
        if err != nil {
            return fmt.Errorf("failed to read all stats: %w", err)
        }

        var ContainerStats DockerStats 
        err = json.Unmarshal(body, &ContainerStats)
        if err != nil {
            return fmt.Errorf("failed to Unmarshal stats: %w", err)
        }

        fmt.Printf("Container: %s\n",  ctr.ID)
        fmt.Printf("CPU usage: %d\n",  ContainerStats.CPUStats.CPUUsage.TotalUsage)
        fmt.Printf("Memory usage: %d\nMemory limit: %d\n",  ContainerStats.MemoryStats.Usage, ContainerStats.MemoryStats.Limit)

    }
    return nil
}

