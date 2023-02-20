package cmd

import (
	"github.com/gophercloud/utils/openstack/clientconfig"
	"github.com/mchlumsky/mracek/config"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

func cloudsFromFiles() []clientconfig.Cloud {
	opts := config.YAMLOpts{Directory: viper.GetString("os-config-dir")}

	clds, err := opts.LoadCloudsYAML()
	cobra.CheckErr(err)

	secs, err := opts.LoadSecureCloudsYAML()
	cobra.CheckErr(err)

	pubs, err := opts.LoadPublicCloudsYAML()
	cobra.CheckErr(err)

	clouds := make([]clientconfig.Cloud, 0, len(clds)+len(secs)+len(pubs))

	for _, c := range clds {
		clouds = append(clouds, c)
	}

	for _, s := range secs {
		clouds = append(clouds, s)
	}

	for _, p := range pubs {
		clouds = append(clouds, p)
	}

	return clouds
}

func domainNameCompletionFn(_ *cobra.Command, _ []string, _ string) ([]string, cobra.ShellCompDirective) {
	clouds := cloudsFromFiles()

	properties := make([]string, 0, len(clouds))

	for _, c := range clouds {
		a := c.AuthInfo
		for _, d := range []string{a.UserDomainName, a.ProjectDomainName, a.DomainName} {
			if d != "" {
				properties = append(properties, d)
			}
		}
	}

	return properties, cobra.ShellCompDirectiveNoFileComp
}

func domainIDCompletionFn(_ *cobra.Command, _ []string, _ string) ([]string, cobra.ShellCompDirective) {
	clouds := cloudsFromFiles()

	properties := make([]string, 0, len(clouds))

	for _, c := range clouds {
		a := c.AuthInfo
		for _, d := range []string{a.DomainID, a.UserDomainID, a.ProjectDomainID, a.DefaultDomain} {
			if d != "" {
				properties = append(properties, d)
			}
		}
	}

	return properties, cobra.ShellCompDirectiveNoFileComp
}

func completionFn(setFn func(clientconfig.Cloud, []string, int)) func(
	cmd *cobra.Command,
	args []string,
	toComplete string,
) ([]string, cobra.ShellCompDirective) {
	return func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		clouds := cloudsFromFiles()

		properties := make([]string, len(clouds))
		for i, c := range clouds {
			setFn(c, properties, i)
		}

		return properties, cobra.ShellCompDirectiveNoFileComp
	}
}

func setAuthURL(c clientconfig.Cloud, properties []string, idx int) {
	properties[idx] = c.AuthInfo.AuthURL
}

func setUserName(c clientconfig.Cloud, properties []string, idx int) {
	properties[idx] = c.AuthInfo.Username
}

func setUserID(c clientconfig.Cloud, properties []string, idx int) {
	properties[idx] = c.AuthInfo.UserID
}

func setAppCredID(c clientconfig.Cloud, properties []string, idx int) {
	properties[idx] = c.AuthInfo.ApplicationCredentialID
}

func setAppCredName(c clientconfig.Cloud, properties []string, idx int) {
	properties[idx] = c.AuthInfo.ApplicationCredentialName
}

func setProjectName(c clientconfig.Cloud, properties []string, idx int) {
	properties[idx] = c.AuthInfo.ProjectName
}

func setProjectID(c clientconfig.Cloud, properties []string, idx int) {
	properties[idx] = c.AuthInfo.ProjectID
}

func setRegionName(c clientconfig.Cloud, properties []string, idx int) {
	properties[idx] = c.RegionName
}

//nolint:funlen
func addAllFlags(cmd *cobra.Command, cloud *clientconfig.Cloud) {
	cmd.Flags().StringVar(&cloud.AuthInfo.AuthURL, "auth-url", "", "")
	cobra.CheckErr(cmd.RegisterFlagCompletionFunc("auth-url", completionFn(setAuthURL)))

	cmd.Flags().StringVar(&cloud.AuthInfo.Token, "token", "", "")

	cmd.Flags().StringVar(&cloud.AuthInfo.Username, "username", "", "")
	cobra.CheckErr(cmd.RegisterFlagCompletionFunc("username", completionFn(setUserName)))

	cmd.Flags().StringVar(&cloud.AuthInfo.UserID, "user-id", "", "")
	cobra.CheckErr(cmd.RegisterFlagCompletionFunc("user-id", completionFn(setUserID)))

	cmd.Flags().StringVar(&cloud.AuthInfo.Password, "password", "", "")

	cmd.Flags().StringVar(&cloud.AuthInfo.ApplicationCredentialID, "application-credential-id", "", "")
	cobra.CheckErr(cmd.RegisterFlagCompletionFunc("application-credential-id", completionFn(setAppCredID)))

	cmd.Flags().StringVar(&cloud.AuthInfo.ApplicationCredentialName, "application-credential-name", "", "")
	cobra.CheckErr(cmd.RegisterFlagCompletionFunc("application-credential-name", completionFn(setAppCredName)))

	cmd.Flags().StringVar(&cloud.AuthInfo.ApplicationCredentialSecret, "application-credential-secret", "", "")

	cmd.Flags().StringVar(&cloud.AuthInfo.ProjectName, "project-name", "", "")
	cobra.CheckErr(cmd.RegisterFlagCompletionFunc("project-name", completionFn(setProjectName)))

	cmd.Flags().StringVar(&cloud.AuthInfo.ProjectID, "project-id", "", "")
	cobra.CheckErr(cmd.RegisterFlagCompletionFunc("project-id", completionFn(setProjectID)))

	cmd.Flags().StringVar(&cloud.AuthInfo.UserDomainName, "user-domain-name", "", "")
	cobra.CheckErr(cmd.RegisterFlagCompletionFunc("user-domain-name", domainNameCompletionFn))

	cmd.Flags().StringVar(&cloud.AuthInfo.UserDomainID, "user-domain-id", "", "")
	cobra.CheckErr(cmd.RegisterFlagCompletionFunc("user-domain-id", domainIDCompletionFn))

	cmd.Flags().StringVar(&cloud.AuthInfo.ProjectDomainName, "project-domain-name", "", "")
	cobra.CheckErr(cmd.RegisterFlagCompletionFunc("project-domain-name", domainNameCompletionFn))

	cmd.Flags().StringVar(&cloud.AuthInfo.ProjectDomainID, "project-domain-id", "", "")
	cobra.CheckErr(cmd.RegisterFlagCompletionFunc("project-domain-id", domainIDCompletionFn))

	cmd.Flags().StringVar(&cloud.AuthInfo.DomainName, "domain-name", "", "")
	cobra.CheckErr(cmd.RegisterFlagCompletionFunc("domain-name", domainNameCompletionFn))

	cmd.Flags().StringVar(&cloud.AuthInfo.DomainID, "domain-id", "", "")
	cobra.CheckErr(cmd.RegisterFlagCompletionFunc("domain-id", domainIDCompletionFn))

	cmd.Flags().StringVar(&cloud.AuthInfo.DefaultDomain, "default-domain", "", "")
	cobra.CheckErr(cmd.RegisterFlagCompletionFunc("default-domain", domainIDCompletionFn))

	cmd.Flags().BoolVar(
		&cloud.AuthInfo.AllowReauth,
		"allow-reauth",
		false,
		"allow Gophercloud to attempt to re-authenticate automatically if/when your token expires",
	)

	cmd.Flags().StringVar(&cloud.Profile, "profile", "", "")
	cobra.CheckErr(cmd.RegisterFlagCompletionFunc("profile",
		func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
			profiles, err := clientconfig.LoadPublicCloudsYAML()
			cobra.CheckErr(err)

			result := make([]string, 0, len(profiles))

			for k := range profiles {
				result = append(result, k)
			}

			return result, cobra.ShellCompDirectiveNoFileComp
		}))

	cmd.Flags().StringVar(&cloud.Cloud, "cloud", "", "")

	cmd.Flags().StringVar((*string)(&cloud.AuthType), "auth-type", "", "")

	cmd.Flags().StringVar(&cloud.RegionName, "region-name", "", "")
	cobra.CheckErr(cmd.RegisterFlagCompletionFunc("region-name", completionFn(setRegionName)))

	cmd.Flags().StringVar(&cloud.EndpointType, "endpoint-type", "", "")

	interfaceFn := func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return []string{"public", "internal", "admin"}, cobra.ShellCompDirectiveNoFileComp
	}
	cobra.CheckErr(cmd.RegisterFlagCompletionFunc("endpoint-type", interfaceFn))

	cmd.Flags().StringVar(&cloud.Interface, "interface", "", "")
	cobra.CheckErr(cmd.RegisterFlagCompletionFunc("interface", interfaceFn))

	cmd.Flags().StringVar(&cloud.IdentityAPIVersion, "identity-api-version", "", "")

	cmd.Flags().StringVar(&cloud.VolumeAPIVersion, "volume-api-version", "", "")

	cmd.Flags().BoolVar(
		cloud.Verify,
		"verify",
		false,
		"whether or not SSL API requests should be verified",
	)

	cmd.Flags().StringVar(&cloud.CACertFile, "ca-cert-file", "", "")

	cmd.Flags().StringVar(&cloud.ClientCertFile, "client-cert-file", "", "")

	cmd.Flags().StringVar(&cloud.ClientKeyFile, "client-key-file", "", "")

	cmd.Flags().Bool("password-prompt", false, "Prompt for password")
}

func isFlagPassed(cmd *cobra.Command, name string) bool {
	found := false

	cmd.Flags().Visit(func(f *pflag.Flag) {
		if f.Name == name {
			found = true
		}
	})

	return found
}
