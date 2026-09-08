package generator

import (
	"fmt"
	"os"
	"path/filepath"
)

func GenerateCobraRoot(cwd, projectName string) {
	cobraFile := filepath.Join(cwd, "cmd", "root.go")
	content := fmt.Sprintf(`package cmd

import (
    "fmt"
    "os"

    "github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
    Use:   "%s",
    Short: "%s CLI",
    Run: func(cmd *cobra.Command, args []string) {
        fmt.Println("Welcome to %s CLI")
    },
}

func Execute() {
    if err := rootCmd.Execute(); err != nil {
        fmt.Println(err)
        os.Exit(1)
    }
}
`, projectName, projectName, projectName)
	os.WriteFile(cobraFile, []byte(content), 0644)
	fmt.Println("[✔] Added Cobra root.go")
}
