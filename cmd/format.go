package cmd

import (
	"github.com/spf13/cobra"
	"github.com/tr2rl/tr2rl/internal/parser"
	"github.com/tr2rl/tr2rl/internal/printer"
)

var formatCmd = &cobra.Command{
	Use:   "format [file]",
	Short: "Pretty-print a standard tree from any input",
	Long: `Takes any messy input (indented lists, partial trees, path lists) and outputs 
a perfectly formatted Unicode tree structure. Useful for documentation or verifying 
how tr2rl interprets your input.`,
	Example: `  # Clean up a messy list from clipboard
  tr2rl format --clipboard

  # verify how a Windows tree is parsed
  tr2rl format windows_output.txt`,
	Args: cobra.MaximumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		in, err := readInputFromCmd(cmd, args)
		if err != nil {
			return err
		}

		res := parser.Parse(in)

		// If we encounter parse errors, warning user might be nice,
		// but typically format just does its best.

		printer.PrintTree(res.Nodes)
		return nil
	},
}

func init() {
	rootCmd.AddCommand(formatCmd)
}
