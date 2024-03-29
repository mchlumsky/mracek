package cmd

import (
	"errors"
	"fmt"
	"io/fs"
	"sort"
	"strings"

	"github.com/mchlumsky/mracek/config"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// NewListProfileCommand creates the show command.
func NewListProfileCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "list-profiles",
		Short: "List profiles",
		RunE: func(cmd *cobra.Command, _ []string) error {
			opts := config.YAMLOpts{Directory: viper.GetString("os-config-dir")}

			profileNames, err := opts.AllProfileNames()
			sort.Strings(profileNames)

			if err != nil {
				if errors.Is(err, fs.ErrNotExist) {
					return fmt.Errorf("try creating some profiles first with 'mracek create-profile': %w", err)
				}

				return fmt.Errorf("failed to load profile names: %w", err)
			}

			out := strings.Join(profileNames, "\n")
			_, _ = fmt.Fprintln(cmd.OutOrStdout(), out)

			return nil
		},
	}
}
