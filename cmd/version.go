package cmd

import (
	"encoding/json"
	"fmt"

	"github.com/he8um/daryaft/internal/config"
	"github.com/he8um/daryaft/pkg/version"
	"github.com/spf13/cobra"
)

func newVersionCommand() *cobra.Command {
	var jsonOutput bool

	command := &cobra.Command{
		Use:          "version",
		Short:        "Print Daryaft version information",
		SilenceUsage: true,
		Args:         cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			info := version.Info()
			out := cmd.OutOrStdout()

			if jsonOutput {
				data, err := json.MarshalIndent(info, "", "  ")
				if err != nil {
					return fmt.Errorf("encode version info: %w", err)
				}
				fmt.Fprintln(out, string(data))
				return nil
			}

			fmt.Fprintf(out, "%s version: %s\n", config.AppName, info.Version)
			fmt.Fprintf(out, "commit: %s\n", info.Commit)
			fmt.Fprintf(out, "build date: %s\n", info.Date)
			fmt.Fprintf(out, "built by: %s\n", info.BuiltBy)
			fmt.Fprintf(out, "go version: %s\n", info.GoVersion)

			return nil
		},
	}
	command.Flags().BoolVar(&jsonOutput, "json", false, "print version information as JSON")
	return command
}

func init() {
	rootCmd.AddCommand(newVersionCommand())
}
