package cmd

import (
	"fmt"

	appconfig "github.com/he8um/daryaft/internal/config"
	"github.com/he8um/daryaft/internal/download"
	"github.com/he8um/daryaft/internal/downloader"
	"github.com/he8um/daryaft/internal/utils"
	"github.com/spf13/cobra"
)

type downloadFlagValues struct {
	file     string
	output   string
	name     string
	dryRun   bool
	retries  int
	resume   bool
	noResume bool
}

var (
	rootDownloadFlags downloadFlagValues
	subDownloadFlags  downloadFlagValues
)

var downloadCmd = &cobra.Command{
	Use:   "download [url...]",
	Short: "Download URLs or show a dry-run plan",
	Long: `Download one or more HTTP/HTTPS URLs or show a dry-run plan.

Multiple URLs are downloaded sequentially. Use --dry-run to inspect how Daryaft
will interpret URLs, batch files, output options, retries, and resume settings
before any network request is made.`,
	Example: `  daryaft download https://example.com/file.zip --dry-run
  daryaft download https://example.com/file.zip
  daryaft download https://example.com/a.txt https://example.com/b.txt
  daryaft download -f urls.txt --dry-run
  daryaft download -f urls.txt
  daryaft download https://example.com/file.zip --output ~/Downloads`,
	RunE: func(cmd *cobra.Command, args []string) error {
		return runDownload(cmd, args, subDownloadFlags)
	},
}

func init() {
	addDownloadFlags(rootCmd, &rootDownloadFlags)
	addDownloadFlags(downloadCmd, &subDownloadFlags)
	rootCmd.AddCommand(downloadCmd)
}

func addDownloadFlags(command *cobra.Command, flags *downloadFlagValues) {
	flags.resume = true

	command.Flags().StringVarP(&flags.file, "file", "f", "", "read URLs from a file")
	command.Flags().StringVarP(&flags.output, "output", "o", "", "output directory")
	command.Flags().StringVar(&flags.name, "name", "", "filename for a single URL")
	command.Flags().BoolVar(&flags.dryRun, "dry-run", false, "validate inputs and print the download plan")
	command.Flags().IntVar(&flags.retries, "retries", 3, "retry attempts after the initial attempt, 0-20")
	command.Flags().BoolVar(&flags.resume, "resume", true, "resume interrupted partial downloads")
	command.Flags().BoolVar(&flags.noResume, "no-resume", false, "disable resume and restart partial downloads")
}

func runDownload(cmd *cobra.Command, args []string, flags downloadFlagValues) error {
	cfg, err := appconfig.LoadEffective()
	if err != nil {
		return err
	}
	flags = applyConfigDefaultsToDownloadFlags(cmd, flags, cfg)

	options := download.Options{
		URLs:     args,
		File:     flags.file,
		Output:   flags.output,
		Name:     flags.name,
		DryRun:   flags.dryRun,
		Retries:  flags.retries,
		Resume:   flags.resume,
		NoResume: flags.noResume,
	}

	plan, err := download.BuildPlan(options)
	if err != nil {
		return err
	}

	if options.DryRun {
		fmt.Fprintln(cmd.OutOrStdout(), plan.DryRunString())
		return nil
	}

	if len(plan.URLs) > 1 {
		savingPrinted := make(map[int]bool)
		result := downloader.New().DownloadBatch(plan, downloader.BatchHandlers{
			ItemStarted: func(item downloader.BatchItem) {
				fmt.Fprintf(cmd.OutOrStdout(), "[%d/%d] Downloading: %s\n", item.Index, item.Total, item.URL)
			},
			Event: func(item downloader.BatchItem, event downloader.Event) {
				switch event.Type {
				case downloader.EventStarted:
					if !savingPrinted[item.Index] {
						fmt.Fprintf(cmd.OutOrStdout(), "Saving to: %s\n", event.TargetPath)
						savingPrinted[item.Index] = true
					}
				case downloader.EventProgress:
					printProgress(cmd, event)
				case downloader.EventRetrying:
					printRetrying(cmd, event)
				case downloader.EventResuming, downloader.EventRestarting, downloader.EventWarning:
					printMessage(cmd, event)
				case downloader.EventCompleted:
					fmt.Fprintf(cmd.OutOrStdout(), "Completed: %s\n", event.TargetPath)
				case downloader.EventFailed:
					fmt.Fprintf(cmd.OutOrStdout(), "Failed: %s\n", event.Error)
				}
			},
		})

		fmt.Fprintln(cmd.OutOrStdout(), result.SummaryString())
		return result.Err()
	}

	fmt.Fprintf(cmd.OutOrStdout(), "Downloading: %s\n", plan.URLs[0])

	savingPrinted := false
	_, err = downloader.New().DownloadWithEvents(plan, func(event downloader.Event) {
		switch event.Type {
		case downloader.EventStarted:
			if !savingPrinted {
				fmt.Fprintf(cmd.OutOrStdout(), "Saving to: %s\n", event.TargetPath)
				savingPrinted = true
			}
		case downloader.EventProgress:
			printProgress(cmd, event)
		case downloader.EventRetrying:
			printRetrying(cmd, event)
		case downloader.EventResuming, downloader.EventRestarting, downloader.EventWarning:
			printMessage(cmd, event)
		case downloader.EventCompleted:
			fmt.Fprintf(cmd.OutOrStdout(), "Completed: %s\n", event.TargetPath)
		}
	})
	if err != nil {
		return err
	}

	return nil
}

func applyConfigDefaultsToDownloadFlags(cmd *cobra.Command, flags downloadFlagValues, cfg appconfig.Config) downloadFlagValues {
	if !localFlagChanged(cmd, "output") && cfg.DownloadDir != "" {
		flags.output = cfg.DownloadDir
	}
	if !localFlagChanged(cmd, "retries") {
		flags.retries = cfg.Retries
	}
	if !localFlagChanged(cmd, "resume") && !localFlagChanged(cmd, "no-resume") {
		flags.resume = cfg.Resume
	}
	return flags
}

func printProgress(cmd *cobra.Command, event downloader.Event) {
	if event.TotalBytes > 0 {
		fmt.Fprintf(
			cmd.OutOrStdout(),
			"Progress: %d / %d bytes (%.1f%%) | %s\n",
			event.DownloadedBytes,
			event.TotalBytes,
			event.Percent,
			utils.FormatSpeed(event.SpeedBytesPerSecond),
		)
		return
	}

	fmt.Fprintf(
		cmd.OutOrStdout(),
		"Progress: %d bytes | %s\n",
		event.DownloadedBytes,
		utils.FormatSpeed(event.SpeedBytesPerSecond),
	)
}

func printRetrying(cmd *cobra.Command, event downloader.Event) {
	fmt.Fprintf(
		cmd.OutOrStdout(),
		"Retrying %d/%d in %s: %v\n",
		event.Attempt,
		event.MaxAttempts,
		event.NextDelay,
		event.Error,
	)
}

func printMessage(cmd *cobra.Command, event downloader.Event) {
	if event.Message == "" {
		return
	}
	fmt.Fprintln(cmd.OutOrStdout(), event.Message)
}

func hasDownloadFlagChanges(cmd *cobra.Command) bool {
	for _, name := range []string{"file", "output", "name", "dry-run", "retries", "resume", "no-resume"} {
		if localFlagChanged(cmd, name) {
			return true
		}
	}
	return false
}

func localFlagChanged(cmd *cobra.Command, name string) bool {
	flag := cmd.Flags().Lookup(name)
	return flag != nil && flag.Changed
}
