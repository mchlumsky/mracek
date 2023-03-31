package cmd

import (
	"fmt"

	"github.com/gophercloud/utils/openstack/clientconfig"
	"github.com/mchlumsky/mracek/config"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"gopkg.in/yaml.v3"
)

func NewShowCloudCommand() *cobra.Command {
	flags := showFlags{}

	cmd := &cobra.Command{
		Use:   "show-cloud CLOUD",
		Short: "Show cloud details",
		Long:  "Show cloud details",
		Args:  cobra.ExactArgs(1),
		RunE:  showCloudCommandRun,
		ValidArgsFunction: func() ValidArgsFunc {
			return func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
				opts := config.YAMLOpts{Directory: viper.GetString("os-config-dir")}

				return validArgsFunction(opts.AllCloudNames)(cmd, args, toComplete)
			}
		}(),
	}

	cmd.Flags().BoolVarP(&flags.unmask, "unmask", "u", false, "show password in clear text")

	return cmd
}

func showCloudCommandRun(cmd *cobra.Command, args []string) error {
	opts := config.YAMLOpts{Directory: viper.GetString("os-config-dir")}
	co := clientconfig.ClientOpts{Cloud: args[0], YAMLOpts: opts}

	cloud, err := clientconfig.GetCloudFromYAML(&co)
	if err != nil {
		return err
	}

	unmask, err := cmd.Flags().GetBool("unmask")
	if err != nil {
		return err
	}

	if !unmask && cloud.AuthInfo.Password != "" {
		cloud.AuthInfo.Password = "<masked>"
	}

	out, err := yaml.Marshal(cloud)
	if err != nil {
		return err
	}

	_, _ = fmt.Fprint(cmd.OutOrStdout(), "---\n"+string(out))

	return nil
}
