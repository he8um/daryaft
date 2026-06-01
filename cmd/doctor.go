package cmd

import (
	"fmt"

	"github.com/he8um/daryaft/internal/doctor"
	"github.com/spf13/cobra"
)

func newDoctorCommand() *cobra.Command {
	var jsonOutput bool
	var strict bool

	command := &cobra.Command{
		Use:          "doctor",
		Short:        "Run local environment diagnostics",
		SilenceUsage: true,
		Long: `Run local diagnostics for the Daryaft environment.

The doctor command checks runtime details, version metadata, config loading,
default download directory writability, terminal environment hints, optional
tools, and currently skipped remote checks.`,
		Args: cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			report := doctor.Run(doctor.Options{})
			if jsonOutput {
				data, err := doctor.FormatJSONWithOptions(report, strict)
				if err != nil {
					return fmt.Errorf("encode doctor report: %w", err)
				}
				fmt.Fprintln(cmd.OutOrStdout(), string(data))
			} else {
				fmt.Fprint(cmd.OutOrStdout(), doctor.Format(report))
				if strict && !report.Failed() && report.Warned() {
					fmt.Fprintln(cmd.OutOrStdout(), "Strict mode: warnings treated as failures")
				}
			}
			if !report.OK(strict) {
				return fmt.Errorf("doctor found issues")
			}
			return nil
		},
	}
	command.Flags().BoolVar(&jsonOutput, "json", false, "print diagnostics as JSON")
	command.Flags().BoolVar(&strict, "strict", false, "return non-zero when warnings are present")
	return command
}

func init() {
	rootCmd.AddCommand(newDoctorCommand())
}
