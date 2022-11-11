package cmd

import "testing"

//nolint:funlen
func TestSetCloudCommand(t *testing.T) {
	confDir := t.TempDir()

	data := []testItem{
		{
			name:        "create cloud new-cloud",
			env:         nil,
			args:        []string{"--os-config-dir", confDir, "create-cloud", "new-cloud"},
			expected:    "",
			expectedErr: nil,
		},
		{
			name:        "show cloud new-cloud before set-cloud",
			env:         nil,
			args:        []string{"--os-config-dir", confDir, "show-cloud", "new-cloud"},
			expected:    "---\nauth: {}\nverify: false\n",
			expectedErr: nil,
		},
		{
			name: "set cloud fields",
			env:  nil,
			args: []string{
				"--os-config-dir",
				confDir,
				"set-cloud",
				"--allow-reauth",
				"--application-credential-id",
				"appid",
				"--application-credential-name",
				"appname",
				"--application-credential-secret",
				"appsecret",
				"--auth-type",
				"token",
				"--auth-url",
				"http://example.com:5000/v3",
				"--ca-cert-file",
				"ca.cert",
				"--client-cert-file",
				"client.cert",
				"--client-key-file",
				"client.key",
				"--cloud",
				"cloud1",
				"--default-domain",
				"domain1",
				"--domain-id",
				"domainid1",
				"--domain-name",
				"domainname1",
				"--endpoint-type",
				"admin",
				"--identity-api-version",
				"3",
				"--interface",
				"internal",
				"--password",
				"very_secret",
				"--profile",
				"profile1",
				"--project-domain-id",
				"domainid2",
				"--project-domain-name",
				"domainname2",
				"--project-id",
				"project1",
				"--project-name",
				"projectname1",
				"--region-name",
				"region1",
				"--token",
				"token1",
				"--user-domain-id",
				"domainid3",
				"--user-domain-name",
				"domainname3",
				"--user-id",
				"user1",
				"--username",
				"username1",
				"--verify",
				"--volume-api-version",
				"3",
				"new-cloud",
			},
			expected:    "",
			expectedErr: nil,
		},
		{
			name:        "create profile cloud1",
			env:         nil,
			args:        []string{"--os-config-dir", confDir, "create-profile", "cloud1"},
			expected:    "",
			expectedErr: nil,
		},
		{
			name:        "create profile profile1",
			env:         nil,
			args:        []string{"--os-config-dir", confDir, "create-profile", "profile1"},
			expected:    "",
			expectedErr: nil,
		},
		{
			name: "show cloud new-cloud after set-cloud",
			env:  nil,
			args: []string{"--os-config-dir", confDir, "show-cloud", "new-cloud"},
			expected: `---
cloud: cloud1
profile: profile1
auth:
    auth_url: http://example.com:5000/v3
    token: token1
    username: username1
    user_id: user1
    password: <masked>
    application_credential_id: appid
    application_credential_name: appname
    application_credential_secret: appsecret
    project_name: projectname1
    project_id: project1
    user_domain_name: domainname3
    user_domain_id: domainid3
    project_domain_name: domainname2
    project_domain_id: domainid2
    domain_name: domainname1
    domain_id: domainid1
    default_domain: domain1
    allow_reauth: true
auth_type: token
region_name: region1
endpoint_type: admin
identity_api_version: "3"
volume_api_version: "3"
verify: true
cacert: ca.cert
cert: client.cert
key: client.key
`,
			expectedErr: nil,
		},
	}

	for _, item := range data {
		t.Run(item.name, testExecuteFunc(item))
	}
}
