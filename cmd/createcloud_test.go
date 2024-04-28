package cmd

import (
	"fmt"
	"testing"

	"github.com/mchlumsky/mracek/config"
)

func TestCreateCloudCommand(t *testing.T) {
	confDir := t.TempDir()

	data := []testItem{
		{
			name: "test cloud creation",
			args: []string{"create-cloud", "--os-config-dir", confDir, "--password", "secret", "new-cloud"},
		},
		{
			name: "test cloud creation twice fail",
			args: []string{"create-cloud", "--os-config-dir", confDir, "new-cloud"},
			expectedErr: fmt.Errorf(
				"failed to create cloud: %w", config.CloudAlreadyExistsError{
					Cloud:    "new-cloud",
					Filename: "clouds.yaml",
				}),
		},
	}

	for _, item := range data {
		t.Run(item.name, testExecuteFunc(item))
	}
}
