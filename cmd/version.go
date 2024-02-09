package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

func NewVersionCommand(version string) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "version",
		Short: "Display version",
		Long:  "Display version",
		Args:  cobra.NoArgs,
		Run:   versionCommand(version),
	}

	return cmd
}

func versionCommand(version string) func(*cobra.Command, []string) {
	return func(cmd *cobra.Command, _ []string) {
		_, _ = fmt.Fprintln(cmd.OutOrStdout(), version)
	}
}
