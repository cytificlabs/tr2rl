package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/tr2rl/tr2rl/internal/templates"
)

var templateCmd = &cobra.Command{
	Use:   "template",
	Short: "Manage and view built-in project templates",
}

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List available templates",
	Run: func(cmd *cobra.Command, args []string) {
		names := templates.List()
		fmt.Println("Available Templates:")
		for _, n := range names {
			fmt.Printf("  - %s\n", n)
		}
	},
}

var showCmd = &cobra.Command{
	Use:   "show [name]",
	Short: "Output a template's content (pipeable)",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		name := args[0]
		content, ok := templates.Get(name)
		if !ok {
			return fmt.Errorf("template not found: %s", name)
		}
		fmt.Println(content)
		return nil
	},
}

func init() {
	rootCmd.AddCommand(templateCmd)
	templateCmd.AddCommand(listCmd)
	templateCmd.AddCommand(showCmd)
}
