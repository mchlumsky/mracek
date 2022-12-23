package cmd

import (
	"github.com/gophercloud/utils/openstack/clientconfig"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

func addAllFlags(cmd *cobra.Command, cloud *clientconfig.Cloud) {
	cmd.Flags().StringVar(&cloud.AuthInfo.AuthURL, "auth-url", "", "")
	cmd.Flags().StringVar(&cloud.AuthInfo.Token, "token", "", "")
	cmd.Flags().StringVar(&cloud.AuthInfo.Username, "username", "", "")
	cmd.Flags().StringVar(&cloud.AuthInfo.UserID, "user-id", "", "")
	cmd.Flags().StringVar(&cloud.AuthInfo.Password, "password", "", "")
	cmd.Flags().StringVar(&cloud.AuthInfo.ApplicationCredentialID, "application-credential-id", "", "")
	cmd.Flags().StringVar(&cloud.AuthInfo.ApplicationCredentialName, "application-credential-name", "", "")
	cmd.Flags().StringVar(&cloud.AuthInfo.ApplicationCredentialSecret, "application-credential-secret", "", "")
	cmd.Flags().StringVar(&cloud.AuthInfo.ProjectName, "project-name", "", "")
	cmd.Flags().StringVar(&cloud.AuthInfo.ProjectID, "project-id", "", "")
	cmd.Flags().StringVar(&cloud.AuthInfo.UserDomainName, "user-domain-name", "", "")
	cmd.Flags().StringVar(&cloud.AuthInfo.UserDomainID, "user-domain-id", "", "")
	cmd.Flags().StringVar(&cloud.AuthInfo.ProjectDomainName, "project-domain-name", "", "")
	cmd.Flags().StringVar(&cloud.AuthInfo.ProjectDomainID, "project-domain-id", "", "")
	cmd.Flags().StringVar(&cloud.AuthInfo.DomainName, "domain-name", "", "")
	cmd.Flags().StringVar(&cloud.AuthInfo.DomainID, "domain-id", "", "")
	cmd.Flags().StringVar(&cloud.AuthInfo.DefaultDomain, "default-domain", "", "")
	cmd.Flags().BoolVar(
		&cloud.AuthInfo.AllowReauth,
		"allow-reauth",
		false,
		"allow Gophercloud to attempt to re-authenticate automatically if/when your token expires",
	)
	cmd.Flags().StringVar(&cloud.Profile, "profile", "", "")
	cmd.Flags().StringVar(&cloud.Cloud, "cloud", "", "")
	cmd.Flags().StringVar((*string)(&cloud.AuthType), "auth-type", "", "")
	cmd.Flags().StringVar(&cloud.RegionName, "region-name", "", "")
	cmd.Flags().StringVar(&cloud.EndpointType, "endpoint-type", "", "")
	cmd.Flags().StringVar(&cloud.Interface, "interface", "", "")
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
