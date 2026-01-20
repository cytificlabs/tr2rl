package cmd

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/tr2rl/tr2rl/internal/parser"
)

var specCmd = &cobra.Command{
	Use:   "spec [file]",
	Short: "Normalize/repair a tree spec (no disk writes)",
	Args:  cobra.MaximumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		in, err := readInputFromCmd(cmd, args)
		if err != nil {
			return err
		}

		res := parser.Parse(in)

		jsonOut, _ := cmd.Flags().GetBool("json")
		if jsonOut {
			enc := json.NewEncoder(os.Stdout)
			enc.SetIndent("", "  ")
			return enc.Encode(res)
		}

		fmt.Println(res.Normalized)
		if len(res.Warnings) > 0 {
			fmt.Fprintln(os.Stderr, "\nWarnings:")
			for _, w := range res.Warnings {
				fmt.Fprintln(os.Stderr, "- "+w)
			}
		}
		return nil
	},
}

func init() {
	rootCmd.AddCommand(specCmd)
	specCmd.Flags().BoolP("verbose", "v", false, "show parsing details")
	specCmd.Flags().Bool("json", false, "output result as JSON")
}
