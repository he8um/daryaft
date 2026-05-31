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
	command.AddCommand(newConfigGetCommand())
	command.AddCommand(newConfigSetCommand())
	command.AddCommand(newConfigResetCommand())
	command.AddCommand(newConfigKeysCommand())
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

func newConfigGetCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "get <key>",
		Short: "Print an effective configuration value",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			cfg, err := appconfig.LoadEffective()
			if err != nil {
				return err
			}
			value, err := appconfig.Get(cfg, args[0])
			if err != nil {
				return err
			}
			fmt.Fprintln(cmd.OutOrStdout(), value)
			return nil
		},
	}
}

func newConfigSetCommand() *cobra.Command {
	command := &cobra.Command{
		Use:   "set <key> <value>",
		Short: "Set a configuration file value",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			cfg, err := appconfig.Load()
			if err != nil {
				return err
			}
			updated, err := appconfig.Set(cfg, args[0], args[1])
			if err != nil {
				return err
			}
			if err := appconfig.Save(updated); err != nil {
				return err
			}
			value, err := appconfig.Get(updated, args[0])
			if err != nil {
				return err
			}
			fmt.Fprintf(cmd.OutOrStdout(), "Updated config: %s=%s\n", args[0], value)
			return nil
		},
	}
	command.Flags().SetInterspersed(false)
	return command
}

func newConfigResetCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "reset",
		Short: "Reset the configuration file to defaults",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			path, err := appconfig.Reset()
			if err != nil {
				return err
			}
			fmt.Fprintf(cmd.OutOrStdout(), "Reset config: %s\n", path)
			return nil
		},
	}
}

func newConfigKeysCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "keys",
		Short: "List supported configuration keys",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			for _, key := range appconfig.SupportedKeys() {
				fmt.Fprintf(cmd.OutOrStdout(), "%s %s\n", key.Key, key.Type)
			}
			return nil
		},
	}
}

func init() {
	rootCmd.AddCommand(newConfigCommand())
}
