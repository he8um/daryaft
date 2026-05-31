package cmd

import (
	"fmt"

	appconfig "github.com/he8um/daryaft/internal/config"
	"github.com/spf13/cobra"
)

var supportedCompletionShells = map[string]struct{}{
	"bash":       {},
	"zsh":        {},
	"fish":       {},
	"powershell": {},
}

func newCompletionCommand() *cobra.Command {
	command := &cobra.Command{
		Use:   "completion [bash|zsh|fish|powershell]",
		Short: "Generate shell completion scripts",
		Long: `Generate shell completion scripts for Daryaft.

The generated script must be loaded by your shell. Installation paths vary by
operating system and shell setup.`,
		Example: `  daryaft completion zsh > "${fpath[1]}/_daryaft"
  daryaft completion bash > /etc/bash_completion.d/daryaft
  daryaft completion fish > ~/.config/fish/completions/daryaft.fish
  daryaft completion powershell`,
		Args:              validateCompletionShell,
		ValidArgs:         []string{"bash", "zsh", "fish", "powershell"},
		ValidArgsFunction: completionShellCompletions,
		RunE: func(cmd *cobra.Command, args []string) error {
			out := cmd.OutOrStdout()
			switch args[0] {
			case "bash":
				return cmd.Root().GenBashCompletion(out)
			case "zsh":
				return cmd.Root().GenZshCompletion(out)
			case "fish":
				return cmd.Root().GenFishCompletion(out, true)
			case "powershell":
				return cmd.Root().GenPowerShellCompletion(out)
			default:
				return unsupportedCompletionShellError(args[0])
			}
		},
	}
	return command
}

func validateCompletionShell(cmd *cobra.Command, args []string) error {
	if err := cobra.ExactArgs(1)(cmd, args); err != nil {
		return err
	}
	if _, ok := supportedCompletionShells[args[0]]; !ok {
		return unsupportedCompletionShellError(args[0])
	}
	return nil
}

func completionShellCompletions(cmd *cobra.Command, args []string, argToComplete string) ([]string, cobra.ShellCompDirective) {
	if len(args) > 0 {
		return nil, cobra.ShellCompDirectiveNoFileComp
	}
	return []string{"bash", "zsh", "fish", "powershell"}, cobra.ShellCompDirectiveNoFileComp
}

func unsupportedCompletionShellError(shell string) error {
	return fmt.Errorf("unsupported shell %q: supported shells are bash, zsh, fish, and powershell", shell)
}

func configKeyCompletions(cmd *cobra.Command, args []string, argToComplete string) ([]string, cobra.ShellCompDirective) {
	if len(args) > 0 {
		return nil, cobra.ShellCompDirectiveNoFileComp
	}
	return supportedConfigKeyNames(), cobra.ShellCompDirectiveNoFileComp
}

func configSetCompletions(cmd *cobra.Command, args []string, argToComplete string) ([]string, cobra.ShellCompDirective) {
	switch len(args) {
	case 0:
		return supportedConfigKeyNames(), cobra.ShellCompDirectiveNoFileComp
	case 1:
		if configKeyType(args[0]) == "bool" {
			return []string{"true", "false"}, cobra.ShellCompDirectiveNoFileComp
		}
		return nil, cobra.ShellCompDirectiveNoFileComp
	default:
		return nil, cobra.ShellCompDirectiveNoFileComp
	}
}

func supportedConfigKeyNames() []string {
	keys := appconfig.SupportedKeys()
	names := make([]string, 0, len(keys))
	for _, key := range keys {
		names = append(names, key.Key)
	}
	return names
}

func configKeyType(name string) string {
	for _, key := range appconfig.SupportedKeys() {
		if key.Key == name {
			return key.Type
		}
	}
	return ""
}

func init() {
	rootCmd.AddCommand(newCompletionCommand())
}
