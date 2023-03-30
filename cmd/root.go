package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"path"
	"strings"

	"github.com/gophercloud/utils/openstack/clientconfig"
	"github.com/mchlumsky/mracek/config"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

const (
	ColorGreen = "\033[32m"
	ColorReset = "\033[0m"
)

type rootFlags struct {
	cfgFile     string
	shell       string
	osCfgDir    string
	osCloudOnly bool
}

type showFlags struct {
	unmask bool
}

type AllCloudsOrProfilesFunc func() ([]string, error)

type ValidArgsFunc func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective)

func validArgsFunction(cloudsOrProfiles AllCloudsOrProfilesFunc) ValidArgsFunc {
	return func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		if len(args) != 0 {
			return nil, cobra.ShellCompDirectiveNoFileComp
		}

		acp, err := cloudsOrProfiles()
		if err != nil {
			_, _ = fmt.Fprintln(cmd.ErrOrStderr(), fmt.Errorf("%w", err))
		}

		return acp, cobra.ShellCompDirectiveNoFileComp
	}
}

// NewRootCommand creates the root command.
func NewRootCommand(flags *rootFlags) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "mracek [CLOUD]",
		Short: "Do things with your OpenStack client configuration",
		Long:  "Do things with your OpenStack client configuration",
		Run:   rootCommandRun(flags),
		Args:  cobra.MaximumNArgs(1),
		ValidArgsFunction: func() ValidArgsFunc {
			return func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
				opts := config.YAMLOpts{Directory: viper.GetString("os-config-dir")}

				return validArgsFunction(opts.AllCloudNames)(cmd, args, toComplete)
			}
		}(),
	}

	cmd.PersistentFlags().StringVar(&flags.cfgFile, "config", "", "config file (default \"$HOME/.mracek.yaml\")")

	cmd.Flags().StringVarP(
		&flags.shell,
		"shell",
		"s",
		os.Getenv("SHELL"),
		"full path to shell to use for sub-shell. Override it with $MRACEK_SHELL or shell config option",
	)
	cobra.CheckErr(viper.BindPFlag("shell", cmd.Flags().Lookup("shell")))

	home, err := os.UserHomeDir()
	cobra.CheckErr(err)

	cmd.PersistentFlags().StringVar(
		&flags.osCfgDir,
		"os-config-dir",
		path.Join(home, ".config/openstack"),
		"OpenStack configuration directory (default \"$HOME/.config/openstack/\")",
	)
	cobra.CheckErr(viper.BindPFlag("os-config-dir", cmd.PersistentFlags().Lookup("os-config-dir")))
	cobra.CheckErr(viper.BindEnv("os-config-dir", "MRACEK_OS_CONFIG_DIR"))

	cmd.Flags().BoolVar(
		&flags.osCloudOnly,
		"os-cloud-only",
		true,
		"if true, only OS_CLOUD will be set in the environment",
	)
	cobra.CheckErr(viper.BindPFlag("os-cloud-only", cmd.Flags().Lookup("os-cloud-only")))
	cobra.CheckErr(viper.BindEnv("os-cloud-only", "MRACEK_OS_CLOUD_ONLY"))

	cmd.CompletionOptions.DisableDescriptions = true
	cmd.SilenceUsage = true

	return cmd
}

// formatCloudsString returns a formatted string from the clouds argument. The cloud found in the environment variable
// OS_CLOUD is highlighted in green.
func formatCloudsString(clouds []string) string {
	if osCloud := os.Getenv("OS_CLOUD"); osCloud != "" {
		for i, c := range clouds {
			if c == osCloud {
				clouds[i] = ColorGreen + c + ColorReset
			}
		}
	}

	return strings.Join(clouds, "\n") + "\n"
}

func setCloudEnv(cloudName string, opts config.YAMLOpts, osCloudOnly bool) error {
	co := clientconfig.ClientOpts{Cloud: cloudName, YAMLOpts: opts}

	cloud, err := clientconfig.GetCloudFromYAML(&co)
	if err != nil {
		return fmt.Errorf("error getting cloud from configuration: %w", err)
	}

	var vars map[string]string
	if osCloudOnly {
		vars = map[string]string{
			"OS_CLOUD": cloudName,
		}
	} else {
		vars = map[string]string{
			"OS_CLOUD":                         cloudName,
			"OS_REGION_NAME":                   cloud.RegionName,
			"OS_USERNAME":                      cloud.AuthInfo.Username,
			"OS_USER_ID":                       cloud.AuthInfo.UserID,
			"OS_PASSWORD":                      cloud.AuthInfo.Password,
			"OS_PROJECT_NAME":                  cloud.AuthInfo.ProjectName,
			"OS_PROJECT_ID":                    cloud.AuthInfo.ProjectID,
			"OS_TENANT_NAME":                   cloud.AuthInfo.ProjectName,
			"OS_TENANT_ID":                     cloud.AuthInfo.ProjectID,
			"OS_AUTH_URL":                      cloud.AuthInfo.AuthURL,
			"OS_DOMAIN_NAME":                   cloud.AuthInfo.DomainName,
			"OS_DOMAIN_ID":                     cloud.AuthInfo.DomainID,
			"OS_USER_DOMAIN_NAME":              cloud.AuthInfo.UserDomainName,
			"OS_USER_DOMAIN_ID":                cloud.AuthInfo.UserDomainID,
			"OS_PROJECT_DOMAIN_NAME":           cloud.AuthInfo.ProjectDomainName,
			"OS_PROJECT_DOMAIN_ID":             cloud.AuthInfo.ProjectDomainID,
			"OS_APPLICATION_CREDENTIAL_SECRET": cloud.AuthInfo.ApplicationCredentialSecret,
			"OS_APPLICATION_CREDENTIAL_ID":     cloud.AuthInfo.ApplicationCredentialID,
			"OS_APPLICATION_CREDENTIAL_NAME":   cloud.AuthInfo.ApplicationCredentialName,
			"OS_TOKEN":                         cloud.AuthInfo.Token,
			"OS_DEFAULT_DOMAIN":                cloud.AuthInfo.DefaultDomain,
		}
	}

	for key, value := range vars {
		if value == "" {
			continue
		}

		if err = os.Setenv(key, value); err != nil {
			return fmt.Errorf("error setting %s environment variable: %w", key, err)
		}
	}

	return nil
}

func runShell(shell string) error {
	cmd := exec.Cmd{
		Path:   shell,
		Env:    os.Environ(),
		Stdin:  os.Stdin,
		Stdout: os.Stdout,
		Stderr: os.Stderr,
	}

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to run shell: %w", err)
	}

	return nil
}

func rootCommandRun(_ *rootFlags) func(cmd *cobra.Command, args []string) {
	return func(cmd *cobra.Command, args []string) {
		opts := config.YAMLOpts{Directory: viper.GetString("os-config-dir")}
		if len(args) == 0 {
			cloudNames, err := opts.AllCloudNames()
			if err != nil {
				_, _ = fmt.Fprintln(cmd.ErrOrStderr(), fmt.Errorf("failed to load cloud names: %w", err))

				os.Exit(1)
			}

			_, _ = fmt.Fprint(cmd.OutOrStdout(), formatCloudsString(cloudNames))
		} else {
			_, _ = fmt.Fprintf(cmd.OutOrStdout(), "Switching to cloud %s\n", args[0])

			err := setCloudEnv(args[0], opts, viper.GetBool("os-cloud-only"))
			if err != nil {
				_, _ = fmt.Fprintln(cmd.ErrOrStderr(), fmt.Errorf("failed to set cloud environment: %w", err))

				os.Exit(1)
			}

			err = runShell(viper.GetString("shell"))
			if err != nil {
				_, _ = fmt.Fprintln(cmd.ErrOrStderr(), fmt.Errorf("failed to run shell: %w", err))

				os.Exit(1)
			}
		}
	}
}

func buildRootCommand() *cobra.Command {
	rf := rootFlags{}
	cobra.OnInitialize(initConfig(&rf))

	rootCmd := NewRootCommand(&rf)

	rootCmd.AddCommand(NewListProfileCommand())
	rootCmd.AddCommand(NewShowCloudCommand())
	rootCmd.AddCommand(NewShowProfileCommand())
	rootCmd.AddCommand(NewCreateCloudCommand())
	rootCmd.AddCommand(NewCreateProfileCommand())
	rootCmd.AddCommand(NewSetCloudCommand())
	rootCmd.AddCommand(NewSetProfileCommand())
	rootCmd.AddCommand(NewDeleteCloudCommand())
	rootCmd.AddCommand(NewDeleteProfileCommand())
	rootCmd.AddCommand(NewCopyCloudCommand())
	rootCmd.AddCommand(NewSmokeTestsCommand())

	return rootCmd
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	rootCmd := buildRootCommand()

	if rootCmd.Execute() != nil {
		os.Exit(1)
	}
}

// initConfig reads in config file and ENV variables if set.
func initConfig(flags *rootFlags) func() {
	return func() {
		viper.SetEnvPrefix("MRACEK")

		if flags.cfgFile != "" {
			// Use config file from the flag.
			viper.SetConfigFile(flags.cfgFile)
		} else {
			// Find home directory.
			home, err := os.UserHomeDir()
			cobra.CheckErr(err)

			// Search config in home directory with name ".mracek" (without extension).
			viper.AddConfigPath(home)
			viper.SetConfigType("yaml")
			viper.SetConfigName(".mracek")
		}

		viper.AutomaticEnv() // read in environment variables that match

		_ = viper.ReadInConfig()
	}
}
