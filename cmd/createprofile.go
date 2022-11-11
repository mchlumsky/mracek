package cmd

import (
	"github.com/gophercloud/utils/openstack/clientconfig"
	"github.com/mchlumsky/mracek/config"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"gopkg.in/yaml.v3"
)

func NewCreateProfileCommand() *cobra.Command {
	profile := clientconfig.Cloud{AuthInfo: &clientconfig.AuthInfo{}, Verify: new(bool)}
	cmd := &cobra.Command{
		Use:   "create-profile [flags] <profile>",
		Short: "Create profile",
		Long:  "Create profile",
		Args:  cobra.ExactArgs(1),
		RunE:  createProfileCommandRunE(&profile),
	}

	addAllFlags(cmd, &profile)

	return cmd
}

func createProfileCommandRunE(profile *clientconfig.Cloud) func(cmd *cobra.Command, args []string) error {
	return func(cmd *cobra.Command, args []string) error {
		opts := config.YAMLOpts{Directory: viper.GetString("os-config-dir")}

		profiles, err := config.LoadAndCheckOSConfigfile("clouds-public.yaml", opts.LoadPublicCloudsYAML, args[0])
		if err != nil {
			return err
		}

		profiles[args[0]] = *profile

		c := map[string]map[string]clientconfig.Cloud{"public-clouds": profiles}

		publicOut, err := yaml.Marshal(&c)
		if err != nil {
			return err
		}

		if err := config.WriteOSConfig(viper.GetString("os-config-dir"), nil, nil, publicOut); err != nil {
			return err
		}

		return nil
	}
}
