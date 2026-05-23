package cmd

import (
	"fmt"
	"os"

	"github.com/he8um/daryaft/internal/app"
	"github.com/he8um/daryaft/internal/config"
	"github.com/spf13/cobra"
)

var (
	noColor bool
	noTUI   bool
	verbose bool
)

var rootCmd = &cobra.Command{
	Use:           config.BinaryName,
	Short:         "Daryaft is a modern terminal downloader.",
	SilenceUsage:  true,
	SilenceErrors: true,
	Long: `Daryaft is a modern terminal downloader written in Go.

It is similar in spirit to wget, with a clean CLI foundation today and a planned
terminal UI, packaging, self-update workflow, and expanded downloader engine in
future milestones.`,
	Example: `  daryaft https://example.com/file.zip
  daryaft https://example.com/file.zip --dry-run
  daryaft -f urls.txt --dry-run
  daryaft download https://example.com/file.zip
  daryaft download https://example.com/file.zip --dry-run
  daryaft version
  daryaft update`,
	Args: cobra.ArbitraryArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) > 0 || hasDownloadFlagChanges(cmd) {
			return runDownload(cmd, args, rootDownloadFlags)
		}

		fmt.Fprintln(cmd.OutOrStdout(), app.InteractivePlaceholder())
		return nil
	},
	CompletionOptions: cobra.CompletionOptions{
		DisableDefaultCmd: true,
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func init() {
	rootCmd.PersistentFlags().BoolVar(&noColor, "no-color", false, "disable colored output")
	rootCmd.PersistentFlags().BoolVar(&noTUI, "no-tui", false, "disable terminal UI when available")
	rootCmd.PersistentFlags().BoolVarP(&verbose, "verbose", "v", false, "enable verbose output")

	rootCmd.SetUsageTemplate(`Usage:
  {{.UseLine}}{{if .HasAvailableSubCommands}}
  {{.CommandPath}} [command]{{end}}

Available Commands:{{range .Commands}}{{if (or .IsAvailableCommand (eq .Name "help"))}}
  {{rpad .Name .NamePadding }} {{.Short}}{{end}}{{end}}

Common Flags:
{{.LocalFlags.FlagUsages | trimTrailingWhitespaces}}{{if .HasAvailableInheritedFlags}}
{{.InheritedFlags.FlagUsages | trimTrailingWhitespaces}}{{end}}

Examples:
{{.Example}}

` + config.FooterText + `
`)
}
