package cmd

import (
	"context"
	"fmt"

	"github.com/he8um/daryaft/internal/inspect"
	"github.com/spf13/cobra"
)

func newInspectCommand() *cobra.Command {
	var jsonOutput bool

	command := &cobra.Command{
		Use:          "inspect <url>",
		Short:        "Inspect URL metadata without downloading",
		SilenceUsage: true,
		Long: `Inspect one HTTP or HTTPS URL and print metadata without saving a file.

Daryaft tries a HEAD request first and may fall back to a small range probe when
servers do not support HEAD or omit useful metadata.`,
		Example: `  daryaft inspect https://example.com/file.zip
  daryaft inspect https://example.com/file.zip --json`,
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			result, err := inspect.URL(context.Background(), inspect.Options{URL: args[0]})
			if err != nil {
				return err
			}

			if jsonOutput {
				data, err := result.FormatJSON()
				if err != nil {
					return fmt.Errorf("encode inspect result: %w", err)
				}
				fmt.Fprintln(cmd.OutOrStdout(), string(data))
				return nil
			}

			fmt.Fprintln(cmd.OutOrStdout(), result.Format())
			return nil
		},
	}
	command.Flags().BoolVar(&jsonOutput, "json", false, "print inspect metadata as JSON")
	return command
}

func init() {
	rootCmd.AddCommand(newInspectCommand())
}
