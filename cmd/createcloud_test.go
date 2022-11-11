package cmd

import (
	"fmt"
	"testing"
)

func TestCreateCloudCommand(t *testing.T) {
	confDir := t.TempDir()

	data := []testItem{
		{
			name:        "test cloud creation",
			env:         nil,
			args:        []string{"create-cloud", "--os-config-dir", confDir, "--password", "secret", "new-cloud"},
			expected:    "",
			expectedErr: nil,
		},
		{
			name:        "test cloud creation twice fail",
			env:         nil,
			args:        []string{"create-cloud", "--os-config-dir", confDir, "new-cloud"},
			expected:    "",
			expectedErr: fmt.Errorf("cloud new-cloud already exists in clouds.yaml"),
		},
	}

	for _, item := range data {
		t.Run(item.name, testExecuteFunc(item))
	}
}
