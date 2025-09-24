package cmd

import (
	"fmt"
	"os"

	"github.com/MikelGV/Contyard/internal/tui"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/spf13/cobra"
)

// rootCmd represents the base command when called without any subcommands
var (
    docker bool

    podman bool

    kubernetes bool

    all bool

    version = "0.0.1-beta"

    rootCmd = &cobra.Command{
        Use:   "contyard",
        Short: "A tui application that monitors and manages Docker/Podman containers",
        Long: `
            A tui application that monitors and manages Docker/Podman containers on a local  
            or remote host. It displays real-time metrics (CPU, memory, network usage) for  
            running containers, it allows you stopping/starting containers, and it shows logs  
            interactively.
    `,
        Version: version,
        // Here i should run the defautl command(AKA: no flags so it should show 
        // every docker/podman, and kubernetes stats)
        Run: func(cmd *cobra.Command, args []string) { 
            if !docker && !podman && !kubernetes && !all {
                docker = true
            }
            t := tea.NewProgram(tui.NewModel(docker, podman, kubernetes, all), tea.WithAltScreen())
            if _, err := t.Run(); err != nil {
                fmt.Printf("there has been an error: %v", err)
                os.Exit(1)
            }
        },
    }
)

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	// rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.Contyard.yaml)")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
    // i need another flag for running it as default and that displays docker, podman, and kubernetes
    rootCmd.Flags().BoolVarP(&docker, "docker", "d", false, "Gets all docker container stats")
    rootCmd.Flags().BoolVarP(&podman, "podman", "p", false, "Gets all podman container stats")
    rootCmd.Flags().BoolVarP(&kubernetes, "kubernetes", "k", false, "Gets all kubernetes pods stats")
    rootCmd.Flags().BoolVarP(&all, "all", "a", false, "Gets all containers and pods stats")
}


