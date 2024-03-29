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

func NewCreateCloudCommand() *cobra.Command {
	cloud := clientconfig.Cloud{AuthInfo: &clientconfig.AuthInfo{}, Verify: new(bool)}
	cmd := &cobra.Command{
		Use:   "create-cloud [flags] CLOUD",
		Short: "Create cloud",
		Long:  "Create cloud",
		Args:  cobra.ExactArgs(1),
		RunE:  createCloudCommandRunE(&cloud),
	}

	addAllFlags(cmd, &cloud)

	return cmd
}

//nolint:funlen,cyclop
func createCloudCommandRunE(cloud *clientconfig.Cloud) func(cmd *cobra.Command, args []string) error {
	return func(cmd *cobra.Command, args []string) (err error) {
		defer func() {
			if err != nil {
				err = fmt.Errorf("failed to create cloud: %w", err)
			}
		}()

		opts := config.YAMLOpts{Directory: viper.GetString("os-config-dir")}

		clouds, err := config.LoadAndCheckOSConfigfile("clouds.yaml", opts.LoadCloudsYAML, args[0])
		if err != nil {
			if errors.Is(err, fs.ErrNotExist) {
				clouds = make(map[string]clientconfig.Cloud)
			} else {
				return err
			}
		}

		secure, err := config.LoadAndCheckOSConfigfile("secure.yaml", opts.LoadSecureCloudsYAML, args[0])
		if err != nil {
			if errors.Is(err, fs.ErrNotExist) {
				secure = make(map[string]clientconfig.Cloud)
			} else {
				return err
			}
		}

		clouds[args[0]] = *cloud

		if cloud.AuthInfo.Password != "" {
			secure[args[0]] = clientconfig.Cloud{AuthInfo: &clientconfig.AuthInfo{Password: cloud.AuthInfo.Password}}
			cloud.AuthInfo.Password = ""
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

			secure[args[0]] = clientconfig.Cloud{AuthInfo: &clientconfig.AuthInfo{Password: string(bytepw)}}
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
