package cmd

import (
	"github.com/gophercloud/utils/openstack/clientconfig"
	"github.com/mchlumsky/mracek/config"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"gopkg.in/yaml.v3"
)

func NewDeleteCloudCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "delete-cloud <cloud>",
		Short: "Delete cloud",
		Long:  "Delete cloud",
		Args:  cobra.ExactArgs(1),
		RunE:  deleteCloudCommandRunE(),
		ValidArgsFunction: func() ValidArgsFunc {
			return func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
				opts := config.YAMLOpts{Directory: viper.GetString("os-config-dir")}

				return validArgsFunction(opts.AllCloudNames)(cmd, args, toComplete)
			}
		}(),
	}

	return cmd
}

func deleteCloudCommandRunE() func(cmd *cobra.Command, args []string) error {
	return func(cmd *cobra.Command, args []string) error {
		opts := config.YAMLOpts{Directory: viper.GetString("os-config-dir")}

		clouds, err := config.LoadAndCheckOSConfigfile("clouds.yaml", opts.LoadCloudsYAML, "")
		if err != nil {
			return err
		}

		secure, err := config.LoadAndCheckOSConfigfile("secure.yaml", opts.LoadSecureCloudsYAML, "")
		if err != nil {
			return err
		}

		delete(clouds, args[0])
		delete(secure, args[0])

		c := map[string]map[string]clientconfig.Cloud{"clouds": clouds}

		cloudsOut, err := yaml.Marshal(&c)
		if err != nil {
			return err
		}

		s := map[string]map[string]clientconfig.Cloud{"clouds": secure}

		secureOut, err := yaml.Marshal(&s)
		if err != nil {
			return err
		}

		if err := config.WriteOSConfig(viper.GetString("os-config-dir"), cloudsOut, secureOut, nil); err != nil {
			return err
		}

		return nil
	}
}
