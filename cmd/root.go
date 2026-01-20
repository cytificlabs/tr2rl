package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "tr2rl",
	Short: "Trees to Reality — messy text → real folders/files",
	Long: `tr2rl (Trees to Reality) is a "magic" parser that turns text-based tree specifications 
into actual directory structures.

It handles:
  - Standard 'tree' output (Unicode/ASCII)
  - Windows 'tree /a' output (with headers ignored)
  - Indented lists (Notion/TextEdit style)
  - Path lists (search results)
  - Mixed/messy input with comments

Default behavior is SAFE: operations are dry-run unless --dry-run=false is set.`,
	Example: `  # Build from a file
  tr2rl build spec.txt

  # Build from clipboard
  tr2rl build --clipboard

  # Format messy text into a clean tree
  tr2rl format list.txt`,
	SilenceUsage: true,
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.PersistentFlags().Bool("clipboard", false, "read input from clipboard")
	rootCmd.PersistentFlags().String("color", "auto", "color output: auto|always|never")
	rootCmd.PersistentFlags().Bool("json", false, "JSON output (for spec/format)")
}
