package cmd

import (
	"fmt"

	"github.com/he8um/daryaft/internal/config"
	"github.com/he8um/daryaft/pkg/version"
	"github.com/spf13/cobra"
)

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print Daryaft version information",
	Args:  cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		info := version.Info()
		out := cmd.OutOrStdout()

		fmt.Fprintf(out, "%s version: %s\n", config.AppName, info.Version)
		fmt.Fprintf(out, "commit: %s\n", info.Commit)
		fmt.Fprintf(out, "build date: %s\n", info.Date)
		fmt.Fprintf(out, "go version: %s\n", info.GoVersion)

		return nil
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
}
