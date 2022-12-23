package cmd

import (
	"fmt"
	"syscall"

	"github.com/gophercloud/utils/openstack/clientconfig"
	"github.com/imdario/mergo"
	"github.com/mchlumsky/mracek/config"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"golang.org/x/term"
	"gopkg.in/yaml.v3"
)

func NewSetProfileCommand() *cobra.Command {
	profileFlags := clientconfig.Cloud{AuthInfo: &clientconfig.AuthInfo{}, Verify: new(bool)}

	cmd := &cobra.Command{
		Use:   "set-profile <profile>",
		Short: "Set profile details",
		Long:  "Set profile details",
		RunE:  setProfileCommandRun(&profileFlags),
		Args:  cobra.ExactArgs(1),
		ValidArgsFunction: func() ValidArgsFunc {
			return func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
				opts := config.YAMLOpts{Directory: viper.GetString("os-config-dir")}

				return validArgsFunction(opts.AllProfileNames)(cmd, args, toComplete)
			}
		}(),
	}

	addAllFlags(cmd, &profileFlags)

	return cmd
}

//nolint:cyclop,funlen
func setProfileCommandRun(profileFlags *clientconfig.Cloud) func(cmd *cobra.Command, args []string) error {
	return func(cmd *cobra.Command, args []string) error {
		opts := config.YAMLOpts{Directory: viper.GetString("os-config-dir")}

		profiles, err := config.LoadAndCheckOSConfigfile("clouds-public.yaml", opts.LoadPublicCloudsYAML, "")
		if err != nil {
			return err
		}

		profile, ok := profiles[args[0]]
		if !ok {
			return fmt.Errorf("error: profile %s not found", args[0])
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

		// override cloud's fields with flags that are not zero-valued
		err = mergo.Merge(&profile, profileFlags, mergo.WithOverride)
		if err != nil {
			return err
		}

		// only change the verify field if --verify is explicitly passed
		// verify is a special case because it's a boolean, and it's zero value is false (a valid value for verify)
		if isFlagPassed(cmd, "verify") {
			v, err := cmd.Flags().GetBool("verify")
			if err != nil {
				return err
			}

			profile.Verify = &v
		}

		// only change the allowReauth field if --allow-reauth is explicitly passed
		// allowReauth is a special case because it's a boolean, and it's zero value is false (a valid value for
		// allowReauth)
		if isFlagPassed(cmd, "allow-reauth") {
			ar, err := cmd.Flags().GetBool("allow-reauth")
			if err != nil {
				return err
			}

			profile.AuthInfo.AllowReauth = ar
		}

		profiles[args[0]] = profile

		p := map[string]map[string]clientconfig.Cloud{"public-clouds": profiles}

		profilesOut, err := yaml.Marshal(&p)
		if err != nil {
			return err
		}

		err = config.WriteOSConfig(viper.GetString("os-config-dir"), nil, nil, profilesOut)
		if err != nil {
			return err
		}

		return nil
	}
}
