package cmd

import (
	"fmt"

	appconfig "github.com/he8um/daryaft/internal/config"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
)

func newConfigCommand() *cobra.Command {
	command := &cobra.Command{
		Use:   "config",
		Short: "Manage Daryaft configuration",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			return cmd.Help()
		},
	}

	command.AddCommand(newConfigPathCommand())
	command.AddCommand(newConfigShowCommand())
	command.AddCommand(newConfigInitCommand())
	return command
}

func newConfigPathCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "path",
		Short: "Print the configuration file path",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			path, err := appconfig.Path()
			if err != nil {
				return fmt.Errorf("resolve config path: %w", err)
			}
			fmt.Fprintln(cmd.OutOrStdout(), path)
			return nil
		},
	}
}

func newConfigShowCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "show",
		Short: "Print the effective configuration",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			cfg, err := appconfig.LoadEffective()
			if err != nil {
				return err
			}
			data, err := yaml.Marshal(cfg)
			if err != nil {
				return fmt.Errorf("encode config: %w", err)
			}
			fmt.Fprint(cmd.OutOrStdout(), string(data))
			return nil
		},
	}
}

func newConfigInitCommand() *cobra.Command {
	var force bool
	command := &cobra.Command{
		Use:   "init",
		Short: "Create the default configuration file",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			path, err := appconfig.Init(force)
			if err != nil {
				return err
			}
			fmt.Fprintf(cmd.OutOrStdout(), "Created config: %s\n", path)
			return nil
		},
	}
	command.Flags().BoolVar(&force, "force", false, "overwrite an existing configuration file")
	return command
}

func init() {
	rootCmd.AddCommand(newConfigCommand())
}
