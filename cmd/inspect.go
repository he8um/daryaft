package cmd

import (
	"context"
	"fmt"

	"github.com/he8um/daryaft/internal/httpopts"
	"github.com/he8um/daryaft/internal/inspect"
	"github.com/spf13/cobra"
)

func newInspectCommand() *cobra.Command {
	var jsonOutput bool
	var proxy string
	var headers []string
	var userAgent string
	var username string
	var password string

	command := &cobra.Command{
		Use:          "inspect <url>",
		Short:        "Inspect URL metadata without downloading",
		SilenceUsage: true,
		Long: `Inspect one HTTP or HTTPS URL and print metadata without saving a file.

Daryaft tries a HEAD request first and may fall back to a small range probe when
servers do not support HEAD or omit useful metadata.`,
		Example: `  daryaft inspect https://example.com/file.zip
  daryaft inspect https://example.com/file.zip --json
  daryaft inspect https://example.com/file.zip --header "X-Token: abc"
  daryaft inspect https://example.com/file.zip --username alice --password secret`,
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			parsedHeaders, err := httpopts.ParseHeaders(headers)
			if err != nil {
				return err
			}
			httpOpts := applyHTTPCredentialEnv(httpopts.Options{
				ProxyURL:  proxy,
				Headers:   parsedHeaders,
				UserAgent: userAgent,
				Username:  username,
				Password:  password,
			})

			result, err := inspect.URL(context.Background(), inspect.Options{
				URL:         args[0],
				HTTPOptions: httpOpts,
			})
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
	command.Flags().StringVar(&proxy, "proxy", "", "proxy URL (http or https)")
	command.Flags().StringArrayVar(&headers, "header", nil, "custom request header in \"Name: Value\" format (repeatable)")
	command.Flags().StringVar(&userAgent, "user-agent", "", "override the default User-Agent header")
	command.Flags().StringVar(&username, "username", "", "HTTP Basic Auth username")
	command.Flags().StringVar(&password, "password", "", "HTTP Basic Auth password")
	return command
}

func init() {
	rootCmd.AddCommand(newInspectCommand())
}
