package cmd

import (
	"fmt"
	"os"

	"github.com/he8um/daryaft/internal/app"
	appconfig "github.com/he8um/daryaft/internal/config"
	"github.com/he8um/daryaft/internal/tui"
	"github.com/spf13/cobra"
)

var (
	noColor    bool
	noTUI      bool
	verbose    bool
	configPath string
)

var rootCmd = &cobra.Command{
	Use:           appconfig.BinaryName,
	Short:         "Daryaft is a modern terminal downloader.",
	SilenceUsage:  true,
	SilenceErrors: true,
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		if configPath == "" {
			return nil
		}
		// config init is allowed to create a new file; skip existence check.
		if cmd.Name() != "init" {
			if _, err := os.Stat(configPath); err != nil {
				if os.IsNotExist(err) {
					return fmt.Errorf("config file not found: %s", configPath)
				}
				return fmt.Errorf("config file error %s: %w", configPath, err)
			}
		}
		appconfig.SetConfigPath(configPath)
		return nil
	},
	Long: `Daryaft is a modern terminal downloader written in Go.

It is similar in spirit to wget, with a clean CLI foundation, an interactive
home screen, and planned packaging, self-update workflow, and expanded
downloader engine in future milestones.`,
	Example: `  daryaft https://example.com/file.zip
  daryaft https://example.com/file.zip --dry-run
  daryaft https://example.com/file.zip --checksum sha256:<hex>
  daryaft -f urls.txt --dry-run
  daryaft download https://example.com/file.zip
  daryaft download https://example.com/file.zip --dry-run
  daryaft inspect https://example.com/file.zip
  daryaft doctor
  daryaft version
  daryaft update`,
	Args: cobra.ArbitraryArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg, err := appconfig.LoadEffective()
		if err != nil {
			return err
		}

		if len(args) > 0 || hasDownloadFlagChanges(cmd) {
			return runDownload(cmd, args, rootDownloadFlags)
		}

		effectiveNoTUI := noTUI
		if !persistentFlagChanged(cmd, "no-tui") && cfg.NoTUI {
			effectiveNoTUI = true
		}
		effectiveNoColor := noColor
		if !persistentFlagChanged(cmd, "no-color") && (cfg.NoColor || appconfig.IsMonoTheme(cfg.Theme)) {
			effectiveNoColor = true
		}

		if !effectiveNoTUI {
			configPath, pathErr := appconfig.Path()
			if pathErr != nil {
				configPath = "(unavailable)"
			}
			configLoaded := false
			if pathErr == nil {
				loaded, err := appconfig.Exists()
				if err == nil {
					configLoaded = loaded
				}
			}
			return tui.Run(tui.Options{
				NoColor:           effectiveNoColor,
				Theme:             cfg.Theme,
				DownloadDir:       cfg.DownloadDir,
				Retries:           cfg.Retries,
				Resume:            cfg.Resume,
				UseConfigDefaults: true,
				NoTUI:             cfg.NoTUI,
				Animations:        cfg.Animations,
				Hyperlinks:        cfg.Hyperlinks,
				UserAgent:         cfg.UserAgent,
				Timeout:           cfg.Timeout,
				ConfigInfo: tui.ConfigInfo{
					Path:   configPath,
					Loaded: configLoaded,
				},
			})
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
	rootCmd.PersistentFlags().StringVar(&configPath, "config", "", "path to configuration file")

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

` + appconfig.FooterText + `
`)
}

func persistentFlagChanged(cmd *cobra.Command, name string) bool {
	flag := cmd.Root().PersistentFlags().Lookup(name)
	return flag != nil && flag.Changed
}
