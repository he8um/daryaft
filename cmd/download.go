package cmd

import (
	"fmt"
	"net/url"
	"path/filepath"
	"strings"
	"time"

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

	verboseMode := effectiveVerbose(cmd)
	startedAt := time.Now()
	printVerbosePlan(cmd, plan, verboseMode)

	if len(plan.URLs) > 1 {
		savingPrinted := make(map[int]bool)
		result := downloader.New().DownloadBatch(plan, downloader.BatchHandlers{
			ItemStarted: func(item downloader.BatchItem) {
				fmt.Fprintf(cmd.OutOrStdout(), "[%d/%d] Downloading: %s\n", item.Index, item.Total, item.URL)
				printVerboseItem(cmd, item, verboseMode)
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
				printVerboseEvent(cmd, event, verboseMode, startedAt)
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
		printVerboseEvent(cmd, event, verboseMode, startedAt)
	})
	if err != nil {
		return err
	}

	return nil
}

func effectiveVerbose(cmd *cobra.Command) bool {
	flag := cmd.Root().PersistentFlags().Lookup("verbose")
	if flag != nil {
		return flag.Changed && verbose || !flag.Changed && verbose
	}
	return verbose
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

func printVerbosePlan(cmd *cobra.Command, plan download.Plan, enabled bool) {
	if !enabled {
		return
	}
	fmt.Fprintf(cmd.OutOrStdout(), "Verbose: output directory: %s\n", displayOutput(plan.Output))
	fmt.Fprintf(cmd.OutOrStdout(), "Verbose: selected filename: %s\n", displayFilename(plan.Name))
	fmt.Fprintf(cmd.OutOrStdout(), "Verbose: retries: %d\n", plan.Retries)
	fmt.Fprintf(cmd.OutOrStdout(), "Verbose: resume enabled: %t\n", plan.Resume)
	if len(plan.URLs) == 1 {
		fmt.Fprintf(cmd.OutOrStdout(), "Verbose: effective URL: %s\n", redactURL(plan.URLs[0]))
	}
}

func printVerboseItem(cmd *cobra.Command, item downloader.BatchItem, enabled bool) {
	if !enabled {
		return
	}
	fmt.Fprintf(cmd.OutOrStdout(), "Verbose: item %d/%d effective URL: %s\n", item.Index, item.Total, redactURL(item.URL))
}

func printVerboseEvent(cmd *cobra.Command, event downloader.Event, enabled bool, startedAt time.Time) {
	if !enabled {
		return
	}

	switch event.Type {
	case downloader.EventStarted:
		if event.Status != "" {
			fmt.Fprintf(cmd.OutOrStdout(), "Verbose: HTTP status: %s\n", event.Status)
		}
		if event.TargetPath != "" {
			fmt.Fprintf(cmd.OutOrStdout(), "Verbose: target path: %s\n", event.TargetPath)
			fmt.Fprintf(cmd.OutOrStdout(), "Verbose: selected filename: %s\n", filepath.Base(event.TargetPath))
		}
		if event.DownloadedBytes > 0 {
			fmt.Fprintf(cmd.OutOrStdout(), "Verbose: resume offset: %d bytes\n", event.DownloadedBytes)
		} else {
			fmt.Fprintln(cmd.OutOrStdout(), "Verbose: resume decision: starting from byte 0")
		}
	case downloader.EventResuming, downloader.EventRestarting:
		if event.Message != "" {
			fmt.Fprintf(cmd.OutOrStdout(), "Verbose: resume decision: %s\n", event.Message)
		}
	case downloader.EventRetrying:
		if event.Error != nil {
			fmt.Fprintf(cmd.OutOrStdout(), "Verbose: retry reason: %v\n", event.Error)
		}
		fmt.Fprintf(cmd.OutOrStdout(), "Verbose: next retry delay: %s\n", event.NextDelay)
	case downloader.EventCompleted:
		fmt.Fprintf(cmd.OutOrStdout(), "Verbose: completion duration: %s\n", time.Since(startedAt).Round(time.Millisecond))
	}
}

func displayOutput(output string) string {
	if strings.TrimSpace(output) == "" {
		return "."
	}
	return output
}

func displayFilename(name string) string {
	if strings.TrimSpace(name) == "" {
		return "auto-detect"
	}
	return strings.TrimSpace(name)
}

func redactURL(rawURL string) string {
	parsed, err := url.Parse(rawURL)
	if err != nil {
		return rawURL
	}
	parsed.User = nil
	parsed.RawQuery = ""
	parsed.Fragment = ""
	return parsed.String()
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
