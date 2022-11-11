package cmd

import (
	"testing"
)

func TestDeleteProfileCommand(t *testing.T) {
	confDir := t.TempDir()

	data := []testItem{
		{
			name:        "test profile creation",
			env:         nil,
			args:        []string{"create-profile", "--os-config-dir", confDir, "new-profile"},
			expected:    "",
			expectedErr: nil,
		},
		{
			name:        "test profile deletion",
			env:         nil,
			args:        []string{"delete-profile", "--os-config-dir", confDir, "new-profile"},
			expected:    "",
			expectedErr: nil,
		},
		{
			name:        "test profile deletion twice (idempotency)",
			env:         nil,
			args:        []string{"delete-profile", "--os-config-dir", confDir, "new-profile"},
			expected:    "",
			expectedErr: nil,
		},
		{
			// tif creation succeeds it's because the deletion was successful
			name:        "test profile creation again",
			env:         nil,
			args:        []string{"create-profile", "--os-config-dir", confDir, "new-profile"},
			expected:    "",
			expectedErr: nil,
		},
	}

	for _, item := range data {
		t.Run(item.name, testExecuteFunc(item))
	}
}
