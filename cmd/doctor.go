package cmd

import (
	"fmt"

	"github.com/he8um/daryaft/internal/doctor"
	"github.com/spf13/cobra"
)

var doctorCmd = &cobra.Command{
	Use:   "doctor",
	Short: "Run local environment diagnostics",
	Long: `Run local diagnostics for the Daryaft environment.

The doctor command checks runtime details, version metadata, config loading,
default download directory writability, terminal environment hints, optional
tools, and currently skipped remote checks.`,
	Args: cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		report := doctor.Run(doctor.Options{})
		fmt.Fprint(cmd.OutOrStdout(), doctor.Format(report))
		if report.Failed() {
			return fmt.Errorf("doctor found critical issues")
		}
		return nil
	},
}

func init() {
	rootCmd.AddCommand(doctorCmd)
}
