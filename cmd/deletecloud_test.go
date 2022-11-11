package cmd

import (
	"testing"
)

func TestDeleteCloudCommand(t *testing.T) {
	confDir := t.TempDir()

	data := []testItem{
		{
			name:        "test cloud creation",
			env:         nil,
			args:        []string{"create-cloud", "--os-config-dir", confDir, "new-cloud"},
			expected:    "",
			expectedErr: nil,
		},
		{
			name:        "test cloud deletion",
			env:         nil,
			args:        []string{"delete-cloud", "--os-config-dir", confDir, "new-cloud"},
			expected:    "",
			expectedErr: nil,
		},
		{
			name:        "test cloud deletion twice (idempotency)",
			env:         nil,
			args:        []string{"delete-cloud", "--os-config-dir", confDir, "new-cloud"},
			expected:    "",
			expectedErr: nil,
		},
	}

	for _, item := range data {
		t.Run(item.name, testExecuteFunc(item))
	}
}
