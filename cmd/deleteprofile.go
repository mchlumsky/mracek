package cmd

import (
	"errors"
	"fmt"
	"io/fs"

	"github.com/gophercloud/utils/openstack/clientconfig"
	"github.com/mchlumsky/mracek/config"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"gopkg.in/yaml.v3"
)

func NewDeleteProfileCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "delete-profile PROFILE",
		Short: "Delete profile",
		Long:  "Delete profile",
		RunE:  deleteProfileCommandRunE(),
		ValidArgsFunction: func() ValidArgsFunc {
			return func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
				opts := config.YAMLOpts{Directory: viper.GetString("os-config-dir")}

				return validArgsFunction(opts.AllProfileNames)(cmd, args, toComplete)
			}
		}(),
	}
}

func deleteProfileCommandRunE() func(cmd *cobra.Command, args []string) error {
	return func(cmd *cobra.Command, args []string) error {
		opts := config.YAMLOpts{Directory: viper.GetString("os-config-dir")}

		profiles, err := config.LoadAndCheckOSConfigfile("clouds-public.yaml", opts.LoadPublicCloudsYAML, "")
		if err != nil {
			if errors.Is(err, fs.ErrNotExist) {
				return fmt.Errorf("try creating some profiles first with 'mracek create-profile': %w", err)
			}

			return err
		}

		delete(profiles, args[0])

		pc := map[string]map[string]clientconfig.Cloud{"public-clouds": profiles}

		publicCloudsOut, err := yaml.Marshal(&pc)
		if err != nil {
			return err
		}

		err = config.WriteOSConfig(viper.GetString("os-config-dir"), nil, nil, publicCloudsOut)

		return err
	}
}
