package cmd

import (
	"errors"
	"fmt"
	"io/fs"
	"os"

	"github.com/mchlumsky/mracek/config"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"gopkg.in/yaml.v3"
)

func NewShowProfileCommand() *cobra.Command {
	flags := showFlags{}

	cmd := &cobra.Command{
		Use:   "show-profile PROFILE",
		Short: "Show profile details",
		Long:  "Show profile details",
		Args:  cobra.ExactArgs(1),
		RunE:  showProfileCommandRun,
		ValidArgsFunction: func() ValidArgsFunc {
			return func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
				opts := config.YAMLOpts{Directory: viper.GetString("os-config-dir")}

				return validArgsFunction(opts.AllProfileNames)(cmd, args, toComplete)
			}
		}(),
	}

	cmd.Flags().BoolVarP(&flags.unmask, "unmask", "u", false, "show password in clear text")

	return cmd
}

func showProfileCommandRun(cmd *cobra.Command, args []string) error {
	opts := config.YAMLOpts{Directory: viper.GetString("os-config-dir")}

	profiles, err := opts.LoadPublicCloudsYAML()
	if err != nil {
		if errors.Is(err, fs.ErrNotExist) {
			return fmt.Errorf("try creating some profiles first with 'mracek create-profile': %w", err)
		}

		return err
	}

	profile, ok := profiles[args[0]]
	if !ok {
		_, _ = fmt.Fprintf(cmd.ErrOrStderr(), "profile %s not found", args[0])

		os.Exit(1)
	}

	unmask, err := cmd.Flags().GetBool("unmask")
	if err != nil {
		return err
	}

	if profile.AuthInfo != nil && !unmask && profile.AuthInfo.Password != "" {
		profile.AuthInfo.Password = "<masked>"
	}

	out, err := yaml.Marshal(profile)
	if err != nil {
		return err
	}

	_, _ = fmt.Fprint(cmd.OutOrStdout(), "---\n"+string(out))

	return nil
}
