package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/tr2rl/tr2rl/internal/fs"
	"github.com/tr2rl/tr2rl/internal/parser"
)

var buildCmd = &cobra.Command{
	Use:   "build [file] [dir]",
	Short: "Create folders/files from a tree spec",
	Long: `Parses the input text and creates the corresponding directory structure on disk.

Input can be:
  - A file path
  - Stdin (pipe)
  - Clipboard (--clipboard)

Safety:
  - Defaults to DRY-RUN mode (preview only).
  - Use --dry-run=false to actually create files.
  - Will NOT overwrite existing files unless --force is used.    

Features:
  - --populate: Intelligently fills created files with boilerplate (e.g. package main for Go).`,
	Example: `  # Preview what would happen
  tr2rl build structure.txt

  # Actually create files in ./my-output
  tr2rl build structure.txt ./my-output --dry-run=false

  # Create from clipboard and auto-fill content
  tr2rl build --clipboard --populate --dry-run=false`,
	Args: cobra.MaximumNArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		in, err := readInputFromCmd(cmd, args[:min(1, len(args))])
		if err != nil {
			return err
		}

		outDir := "."
		if len(args) >= 2 {
			outDir = args[1]
		}

		dryRun, _ := cmd.Flags().GetBool("dry-run")
		force, _ := cmd.Flags().GetBool("force")

		res := parser.Parse(in)

		fmt.Printf("Building structure in: %s\n", outDir)
		if dryRun {
			fmt.Println("--- DRY RUN (No changes will be made) ---")
		}

		populate, _ := cmd.Flags().GetBool("populate")

		return fs.Apply(outDir, res.Nodes, fs.ApplyOptions{DryRun: dryRun, Force: force, Populate: populate})
	},
}

func init() {
	rootCmd.AddCommand(buildCmd)
	// Default behavior: WRITE to disk. Use --dry-run to preview.
	buildCmd.Flags().Bool("dry-run", false, "preview changes without writing to disk")
	buildCmd.Flags().Bool("force", false, "overwrite existing files")
	// Auto-populate is opt-in to avoid surprising users.
	buildCmd.Flags().Bool("populate", false, "auto-fill files with smart boilerplate")
}
