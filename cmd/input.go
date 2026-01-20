package cmd

import (
	"fmt"
	"io"
	"os"

	"github.com/spf13/cobra"
	"github.com/tr2rl/tr2rl/internal/clipboard"
)

// Helper to standardise input reading
func readInputFromCmd(cmd *cobra.Command, args []string) (string, error) {
	useClipboard, _ := cmd.Flags().GetBool("clipboard")
	if useClipboard {
		return clipboard.ReadAll()
	}

	if len(args) > 0 {
		content, err := os.ReadFile(args[0])
		if err != nil {
			return "", fmt.Errorf("failed to read file '%s': %w", args[0], err)
		}
		return string(content), nil
	}

	// Read from stdin if valid
	stat, _ := os.Stdin.Stat()
	if (stat.Mode() & os.ModeCharDevice) == 0 {
		content, err := io.ReadAll(os.Stdin)
		if err != nil {
			return "", err
		}
		return string(content), nil
	}

	return "", fmt.Errorf("no input provided.\nTry:\n  tr2rl build file.txt\n  cat file.txt | tr2rl build -\n  tr2rl build --clipboard")
}
