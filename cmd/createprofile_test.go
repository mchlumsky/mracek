package cmd

import (
	"fmt"
	"testing"

	"github.com/mchlumsky/mracek/config"
)

func TestCreateProfileCommand(t *testing.T) {
	confDir := t.TempDir()

	data := []testItem{
		{
			name:        "test profile creation",
			env:         nil,
			args:        []string{"create-profile", "--os-config-dir", confDir, "--password", "secret", "new-profile"},
			expected:    "",
			expectedErr: nil,
		},
		{
			name:     "test profile creation fail",
			env:      nil,
			args:     []string{"create-profile", "--os-config-dir", confDir, "--password", "secret", "new-profile"},
			expected: "",
			expectedErr: fmt.Errorf(
				"failed to create profile: %w", config.CloudAlreadyExistsError{
					Cloud:    "new-profile",
					Filename: "clouds-public.yaml",
				}),
		},
	}

	for _, item := range data {
		t.Run(item.name, testExecuteFunc(item))
	}
}
