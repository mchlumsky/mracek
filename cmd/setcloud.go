package cmd

import (
	"errors"
	"fmt"
	"io/fs"
	"syscall"

	"dario.cat/mergo"
	"github.com/gophercloud/utils/openstack/clientconfig"
	"github.com/mchlumsky/mracek/config"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"golang.org/x/term"
	"gopkg.in/yaml.v3"
)

func NewSetCloudCommand() *cobra.Command {
	cloudFlags := clientconfig.Cloud{AuthInfo: &clientconfig.AuthInfo{}, Verify: new(bool)}

	cmd := &cobra.Command{
		Use:   "set-cloud [flags] CLOUD",
		Short: "Set cloud details",
		Long:  "Set cloud details",
		RunE:  setCloudCommandRun(&cloudFlags),
		Args:  cobra.ExactArgs(1),
		ValidArgsFunction: func() ValidArgsFunc {
			return func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
				opts := config.YAMLOpts{Directory: viper.GetString("os-config-dir")}

				return validArgsFunction(opts.AllCloudNames)(cmd, args, toComplete)
			}
		}(),
	}

	addAllFlags(cmd, &cloudFlags)

	return cmd
}

//nolint:funlen,gocognit,cyclop
func setCloudCommandRun(cloudFlags *clientconfig.Cloud) func(cmd *cobra.Command, args []string) error {
	return func(cmd *cobra.Command, args []string) (err error) {
		defer func() {
			if err != nil {
				err = fmt.Errorf("failed to set cloud: %w", err)
			}
		}()

		opts := config.YAMLOpts{Directory: viper.GetString("os-config-dir")}
		co := clientconfig.ClientOpts{Cloud: args[0], YAMLOpts: opts}

		// gets the cloud constructed from clouds.yaml+secure.yaml
		cloud, err := clientconfig.GetCloudFromYAML(&co)
		if err != nil {
			if errors.Is(err, fs.ErrNotExist) {
				return fmt.Errorf("try creating some clouds first with 'mracek create-cloud': %w", err)
			}

			return err
		}

		// override cloud's fields with flags that are not zero-valued
		err = mergo.Merge(cloud, cloudFlags, mergo.WithOverride)
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

			cloud.Verify = &v
		}

		// only change the allowReauth field if --allow-reauth is explicitly passed
		// allowReauth is a special case because it's a boolean, and it's zero value is false (a valid value for
		// allowReauth)
		if isFlagPassed(cmd, "allow-reauth") {
			ar, err := cmd.Flags().GetBool("allow-reauth")
			if err != nil {
				return err
			}

			cloud.AuthInfo.AllowReauth = ar
		}

		clouds, err := config.LoadAndCheckOSConfigfile("clouds.yaml", opts.LoadCloudsYAML, "")
		if err != nil {
			return err
		}

		// replace the cloud in clouds.yaml with our newly constructed cloud struct. Anything that was in secure.yaml
		// and didn't belong there is injected into clouds.yaml here.
		clouds[args[0]] = *cloud

		secure, err := config.LoadAndCheckOSConfigfile("secure.yaml", opts.LoadSecureCloudsYAML, "")
		if err != nil {
			return err
		}

		delete(secure, args[0])

		if cloud.AuthInfo != nil {
			secure[args[0]] = clientconfig.Cloud{AuthInfo: &clientconfig.AuthInfo{}}
		}

		// move password from clouds.yaml to secure.yaml
		if cloud.AuthInfo.Password != "" {
			secure[args[0]].AuthInfo.Password = cloud.AuthInfo.Password

			cloud.AuthInfo.Password = ""
		}

		// override password if provided by the flag
		if cloudFlags.AuthInfo.Password != "" {
			secure[args[0]].AuthInfo.Password = cloudFlags.AuthInfo.Password
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

			secure[args[0]].AuthInfo.Password = string(bytepw)
		}

		// move application credential secret from clouds.yaml to secure.yaml
		if cloud.AuthInfo.ApplicationCredentialSecret != "" {
			secure[args[0]].AuthInfo.ApplicationCredentialSecret = cloud.AuthInfo.ApplicationCredentialSecret

			cloud.AuthInfo.ApplicationCredentialSecret = ""
		}

		// override application credential secret if provided by the flag
		if cloudFlags.AuthInfo.ApplicationCredentialSecret != "" {
			secure[args[0]].AuthInfo.ApplicationCredentialSecret = cloudFlags.AuthInfo.ApplicationCredentialSecret
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

		err = config.WriteOSConfig(viper.GetString("os-config-dir"), cloudsOut, secureOut, nil)
		if err != nil {
			return err
		}

		return nil
	}
}
