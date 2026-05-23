package cmd

import (
	"errors"
	"fmt"

	"github.com/he8um/daryaft/internal/download"
	"github.com/he8um/daryaft/internal/downloader"
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
	Short: "Download one URL or show a dry-run plan",
	Long: `Download one HTTP/HTTPS URL or show a dry-run plan.

Batch real downloads are planned but not implemented yet. Use --dry-run to
inspect how Daryaft will interpret URLs, batch files, output options, retries,
and resume settings before any network request is made.`,
	Example: `  daryaft download https://example.com/file.zip --dry-run
  daryaft download https://example.com/file.zip
  daryaft download -f urls.txt --dry-run
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
	command.Flags().IntVar(&flags.retries, "retries", 3, "planned retry count")
	command.Flags().BoolVar(&flags.resume, "resume", true, "enable planned resume support")
	command.Flags().BoolVar(&flags.noResume, "no-resume", false, "disable planned resume support")
}

func runDownload(cmd *cobra.Command, args []string, flags downloadFlagValues) error {
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
		return errors.New(download.BatchNotImplementedMessage)
	}

	fmt.Fprintf(cmd.OutOrStdout(), "Downloading: %s\n", plan.URLs[0])

	result, err := downloader.New().Download(plan)
	if err != nil {
		return err
	}

	fmt.Fprintf(cmd.OutOrStdout(), "Saving to: %s\n", result.Path)
	fmt.Fprintf(cmd.OutOrStdout(), "Completed: %s\n", result.Path)
	return nil
}

func hasDownloadFlagChanges(cmd *cobra.Command) bool {
	for _, name := range []string{"file", "output", "name", "dry-run", "retries", "resume", "no-resume"} {
		if cmd.Flags().Changed(name) {
			return true
		}
	}
	return false
}
