package cmd

import (
	"github.com/gophercloud/utils/openstack/clientconfig"
	"github.com/mchlumsky/mracek/config"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"gopkg.in/yaml.v3"
)

func NewCreateCloudCommand() *cobra.Command {
	cloud := clientconfig.Cloud{AuthInfo: &clientconfig.AuthInfo{}, Verify: new(bool)}
	cmd := &cobra.Command{
		Use:   "create-cloud [flags] <cloud>",
		Short: "Create cloud",
		Long:  "Create cloud",
		Args:  cobra.ExactArgs(1),
		RunE:  createCloudCommandRunE(&cloud),
	}

	addAllFlags(cmd, &cloud)

	return cmd
}

func createCloudCommandRunE(cloud *clientconfig.Cloud) func(cmd *cobra.Command, args []string) error {
	return func(cmd *cobra.Command, args []string) error {
		opts := config.YAMLOpts{Directory: viper.GetString("os-config-dir")}

		clouds, err := config.LoadAndCheckOSConfigfile("clouds.yaml", opts.LoadCloudsYAML, args[0])
		if err != nil {
			return err
		}

		secure, err := config.LoadAndCheckOSConfigfile("secure.yaml", opts.LoadSecureCloudsYAML, args[0])
		if err != nil {
			return err
		}

		clouds[args[0]] = *cloud

		if cloud.AuthInfo.Password != "" {
			secure[args[0]] = clientconfig.Cloud{AuthInfo: &clientconfig.AuthInfo{Password: cloud.AuthInfo.Password}}
			cloud.AuthInfo.Password = ""
		}

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
