package cmd

import (
	"testing"
)

func TestListProfilesCommand(t *testing.T) {
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
			name:        "test profile creation 2",
			env:         nil,
			args:        []string{"create-profile", "--os-config-dir", confDir, "--password", "secret", "new-profile1"},
			expected:    "",
			expectedErr: nil,
		},
		{
			name:        "test profile creation 3",
			env:         nil,
			args:        []string{"create-profile", "--os-config-dir", confDir, "--password", "secret", "new-profile2"},
			expected:    "",
			expectedErr: nil,
		},
		{
			name:        "test profile listing",
			env:         nil,
			args:        []string{"list-profiles", "--os-config-dir", confDir, "list-profiles"},
			expected:    "new-profile\nnew-profile1\nnew-profile2\n",
			expectedErr: nil,
		},
	}

	for _, item := range data {
		t.Run(item.name, testExecuteFunc(item))
	}
}
