package cmd

import (
	"testing"
)

func TestShowProfileCommand(t *testing.T) {
	data := []testItem{
		{
			name: "rackspace",
			env:  map[string]string{},
			args: []string{"--os-config-dir", "testdata", "show-profile", "rackspace"},
			expected: `---
auth:
    auth_url: https://identity.api.rackspacecloud.com/v2.0/
`,
		},
		{
			name: "all_from_profile",
			env:  map[string]string{},
			args: []string{"--os-config-dir", "testdata", "show-profile", "all_fields"},
			expected: `---
profile: some_profile
auth:
    auth_url: https://all.example.com:5000/v3
    token: bizbaz
    username: jdoe
    user_id: "12345"
    password: <masked>
    application_credential_id: app-cred-id
    application_credential_name: app-cred-name
    application_credential_secret: secret
    project_name: Some Project
    project_id: Some Project ID
    user_domain_name: default
    user_domain_id: abcde
    project_domain_name: default
    project_domain_id: fedcba
    domain_name: default
    domain_id: Default
auth_type: token
region_name: ALL
endpoint_type: public
interface: public
identity_api_version: "3"
volume_api_version: "3"
verify: true
cacert: foo.crt
cert: bar.crt
key: bar.key
`,
		},
	}

	for _, item := range data {
		t.Run(item.name, testExecuteFunc(item))
	}
}
