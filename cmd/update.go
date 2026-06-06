package cmd

import (
	"encoding/json"
	"fmt"

	"github.com/he8um/daryaft/internal/update"
	"github.com/spf13/cobra"
)

func newUpdateCommand() *cobra.Command {
	var check bool
	var jsonOutput bool
	var includePrerelease bool
	var repo string

	command := &cobra.Command{
		Use:          "update",
		Short:        "Check for a new Daryaft release",
		SilenceUsage: true,
		Long: `Check whether a newer Daryaft release is available.

daryaft update --check queries the GitHub Releases API and compares
the current version to the latest stable release. It is read-only:
it never downloads, installs, or replaces the current binary.

Auto-update is not yet implemented. Use the appropriate command for
your install channel to upgrade after a new release is available.`,
		Example: `  daryaft update --check
  daryaft update --check --json
  daryaft update --check --include-prerelease`,
		Args: cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			if !check {
				return fmt.Errorf(
					"auto-update is not implemented\n\n" +
						"To check for a new release run:\n\n" +
						"  daryaft update --check",
				)
			}

			opts := update.CheckOptions{
				IncludePrerelease: includePrerelease,
			}
			if repo != "" {
				// Allow owner/repo override for testing; hidden flag.
				opts.Owner, opts.Repo = splitRepo(repo)
			}

			result, err := update.Check(cmd.Context(), opts)
			if err != nil {
				if jsonOutput {
					errJSON, _ := json.MarshalIndent(map[string]string{"error": err.Error()}, "", "  ")
					fmt.Fprintln(cmd.ErrOrStderr(), string(errJSON))
				}
				return fmt.Errorf("update check failed: %w", err)
			}

			if jsonOutput {
				data, err := result.FormatJSON()
				if err != nil {
					return fmt.Errorf("encode update result: %w", err)
				}
				fmt.Fprintln(cmd.OutOrStdout(), string(data))
				return nil
			}

			fmt.Fprint(cmd.OutOrStdout(), result.Format())
			return nil
		},
	}

	command.Flags().BoolVar(&check, "check", false, "check for a new release (required)")
	command.Flags().BoolVar(&jsonOutput, "json", false, "print result as JSON")
	command.Flags().BoolVar(&includePrerelease, "include-prerelease", false, "include pre-release versions in the check")
	command.Flags().StringVar(&repo, "repo", "", "override owner/repo for the release check (for testing)")
	_ = command.Flags().MarkHidden("repo")

	return command
}

// splitRepo splits "owner/repo" into two parts. If the input lacks a slash,
// both owner and repo are returned as empty strings so CheckOptions keeps its
// defaults.
func splitRepo(s string) (owner, repo string) {
	for i := 0; i < len(s); i++ {
		if s[i] == '/' {
			return s[:i], s[i+1:]
		}
	}
	return "", ""
}

func init() {
	rootCmd.AddCommand(newUpdateCommand())
}
