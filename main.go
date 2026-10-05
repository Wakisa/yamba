package main

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/spf13/cobra"
	"github.com/wakisa/yamba/schema"
)

// version string for Yamba
// You can override this at build time with:
// go build -ldflags "-X main.version=v1.0.1" -o yamba .
var version = "v1.0.0"

func main() {
	var rootCmd = &cobra.Command{
		Use:   "yamba",
		Short: "Yamba is a project generator for Go applications.",
		Run: func(cmd *cobra.Command, args []string) {
			// Handle global --version / -v flags
			v, _ := cmd.Flags().GetBool("version")
			if v {
				fmt.Println(version)
				os.Exit(0)
			}
		},
	}

	// Add global flags for --version and -v
	rootCmd.PersistentFlags().BoolP("version", "v", false, "Print the version number of Yamba")

	// init command
	var initCmd = &cobra.Command{
		Use:   "init",
		Short: "Initialize a new project.",
		Run: func(cmd *cobra.Command, args []string) {
			p := tea.NewProgram(schema.NewModel())
			finalModel, err := p.Run()
			if err != nil {
				fmt.Println("Error running the program:", err)
				os.Exit(1)
			}

			_ = finalModel // optional: inspect or use the final model
		},
	}

	// version command (subcommand)
	var versionCmd = &cobra.Command{
		Use:     "version",
		Aliases: []string{"v"},
		Short:   "Print the version number of Yamba",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println(version)
		},
	}

	rootCmd.AddCommand(initCmd)
	rootCmd.AddCommand(versionCmd)

	if err := rootCmd.Execute(); err != nil {
		fmt.Println("Error:", err)
		os.Exit(1)
	}
}
