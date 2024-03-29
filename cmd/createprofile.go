package cmd

import (
	"errors"
	"fmt"
	"io/fs"
	"syscall"

	"github.com/gophercloud/utils/openstack/clientconfig"
	"github.com/mchlumsky/mracek/config"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"golang.org/x/term"
	"gopkg.in/yaml.v3"
)

func NewCreateProfileCommand() *cobra.Command {
	profile := clientconfig.Cloud{AuthInfo: &clientconfig.AuthInfo{}, Verify: new(bool)}
	cmd := &cobra.Command{
		Use:   "create-profile [flags] PROFILE",
		Short: "Create profile",
		Long:  "Create profile",
		Args:  cobra.ExactArgs(1),
		RunE:  createProfileCommandRunE(&profile),
	}

	addAllFlags(cmd, &profile)

	return cmd
}

func createProfileCommandRunE(profile *clientconfig.Cloud) func(cmd *cobra.Command, args []string) error {
	return func(cmd *cobra.Command, args []string) (err error) {
		defer func() {
			if err != nil {
				err = fmt.Errorf("failed to create profile: %w", err)
			}
		}()

		opts := config.YAMLOpts{Directory: viper.GetString("os-config-dir")}

		profiles, err := config.LoadAndCheckOSConfigfile("clouds-public.yaml", opts.LoadPublicCloudsYAML, args[0])
		if err != nil {
			if errors.Is(err, fs.ErrNotExist) {
				profiles = make(map[string]clientconfig.Cloud)
			} else {
				return err
			}
		}

		passPrompt, err := cmd.Flags().GetBool("password-prompt")
		if err != nil {
			return err
		}

		if passPrompt {
			_, _ = fmt.Fprintf(cmd.OutOrStdout(), "Enter password:")

			bytepw, err := term.ReadPassword(syscall.Stdin)
			if err != nil {
				return err
			}

			_, _ = fmt.Fprintf(cmd.OutOrStdout(), "\n")

			if profile.AuthInfo == nil {
				profile.AuthInfo = &clientconfig.AuthInfo{}
			}

			profile.AuthInfo.Password = string(bytepw)
		}

		profiles[args[0]] = *profile

		c := map[string]map[string]clientconfig.Cloud{"public-clouds": profiles}

		publicOut, err := yaml.Marshal(&c)
		if err != nil {
			return err
		}

		err = config.WriteOSConfig(viper.GetString("os-config-dir"), nil, nil, publicOut)

		return err
	}
}
