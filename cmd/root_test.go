package cmd

import (
	"bytes"
	"os"
	"path"
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/mchlumsky/mracek/config"
	"github.com/spf13/viper"
)

func cleanEnv() {
	for _, e := range os.Environ() {
		if strings.HasPrefix(e, "OS_") {
			if err := os.Unsetenv(strings.Split(e, "=")[0]); err != nil {
				panic(err)
			}
		}
	}
}

func TestAllCloudNames(t *testing.T) { //nolint:paralleltest
	expected := []string{
		"alberta",
		"all_fields",
		"all_from_profile",
		"arizona",
		"california",
		"chicago",
		"chicago_legacy",
		"chicago_useprofile",
		"disconnected_clouds",
		"florida",
		"florida_insecure",
		"florida_secure",
		"hawaii",
		"nevada",
		"newmexico",
		"no_fields",
		"philadelphia",
		"philadelphia_complex",
		"region_has_null_char",
		"texas",
		"virginia",
		"yukon",
	}

	opts := config.YAMLOpts{Directory: "testdata/"}

	a, err := opts.AllCloudNames()
	if err != nil {
		t.Fatal(err)
	}

	if diff := cmp.Diff(a, expected); diff != "" {
		t.Error(diff)
	}
}

func TestFormatCloudsStringNoOSCloud(t *testing.T) { //nolint:paralleltest
	clouds := []string{"cloud1", "cloud2"}
	expected := "cloud1\ncloud2\n"

	if diff := cmp.Diff(formatCloudsString(clouds), expected); diff != "" {
		t.Error(diff)
	}
}

func TestFormatCloudsStringWithOSCloud(t *testing.T) {
	t.Setenv("OS_CLOUD", "cloud1")

	clouds := []string{"cloud1", "cloud2"}
	expected := ColorGreen + "cloud1" + ColorReset + "\n" + "cloud2\n"

	if diff := cmp.Diff(formatCloudsString(clouds), expected); diff != "" {
		t.Error(diff)
	}
}

func TestFormatCloudsStringWithInvalidOSCloud(t *testing.T) {
	t.Setenv("OS_CLOUD", "cloud3")

	clouds := []string{"cloud1", "cloud2"}
	expected := "cloud1\ncloud2\n"

	if diff := cmp.Diff(formatCloudsString(clouds), expected); diff != "" {
		t.Error(diff)
	}
}

func TestSetCloudEnvAllFieldsOSCloudOnly(t *testing.T) {
	t.Cleanup(cleanEnv)

	expected := map[string]string{
		"OS_AUTH_URL":                      "https://all.example.com:5000/v3",
		"OS_USERNAME":                      "jdoe",
		"OS_PASSWORD":                      "password",
		"OS_PROJECT_NAME":                  "Some Project",
		"OS_PROJECT_ID":                    "Some Project ID",
		"OS_TENANT_NAME":                   "Some Project",
		"OS_TENANT_ID":                     "Some Project ID",
		"OS_PROJECT_DOMAIN_NAME":           "default",
		"OS_PROJECT_DOMAIN_ID":             "fedcba",
		"OS_USER_DOMAIN_NAME":              "default",
		"OS_USER_DOMAIN_ID":                "abcde",
		"OS_DOMAIN_NAME":                   "default",
		"OS_DOMAIN_ID":                     "Default",
		"OS_APPLICATION_CREDENTIAL_ID":     "app-cred-id",
		"OS_APPLICATION_CREDENTIAL_SECRET": "secret",
		"OS_APPLICATION_CREDENTIAL_NAME":   "app-cred-name",
		"OS_REGION_NAME":                   "ALL",
	}

	for k := range expected {
		if v, ok := os.LookupEnv(k); ok {
			t.Errorf("%s should not be set but is set to \"%s\"", k, v)
		}
	}

	opts := config.YAMLOpts{Directory: "testdata"}
	if err := setCloudEnv("all_fields", opts, true); err != nil {
		t.Errorf("error setting environment: %v", err)
	}

	if v, ok := os.LookupEnv("OS_CLOUD"); !ok {
		t.Error("OS_CLOUD is not set, but it should be")
	} else if v != "all_fields" {
		t.Errorf("OS_CLOUD set to %s when it should be 'all_fields'", v)
	}

	for k := range expected {
		if _, ok := os.LookupEnv(k); ok {
			t.Errorf("%s is set but shouldn't", k)
		}
	}
}

func TestSetCloudEnvAllFields(t *testing.T) {
	t.Cleanup(cleanEnv)

	expected := map[string]string{
		"OS_CLOUD":                         "all_fields",
		"OS_AUTH_URL":                      "https://all.example.com:5000/v3",
		"OS_USERNAME":                      "jdoe",
		"OS_PASSWORD":                      "password",
		"OS_PROJECT_NAME":                  "Some Project",
		"OS_PROJECT_ID":                    "Some Project ID",
		"OS_TENANT_NAME":                   "Some Project",
		"OS_TENANT_ID":                     "Some Project ID",
		"OS_PROJECT_DOMAIN_NAME":           "default",
		"OS_PROJECT_DOMAIN_ID":             "fedcba",
		"OS_USER_DOMAIN_NAME":              "default",
		"OS_USER_DOMAIN_ID":                "abcde",
		"OS_DOMAIN_NAME":                   "default",
		"OS_DOMAIN_ID":                     "Default",
		"OS_APPLICATION_CREDENTIAL_ID":     "app-cred-id",
		"OS_APPLICATION_CREDENTIAL_SECRET": "secret",
		"OS_APPLICATION_CREDENTIAL_NAME":   "app-cred-name",
		"OS_REGION_NAME":                   "ALL",
	}

	for k := range expected {
		if v, ok := os.LookupEnv(k); ok {
			t.Errorf("%s should not be set but is set to \"%s\"", k, v)
		}
	}

	opts := config.YAMLOpts{Directory: "testdata"}
	if err := setCloudEnv("all_fields", opts, false); err != nil {
		t.Errorf("error setting environment: %v", err)
	}

	for k, v := range expected {
		if actual, ok := os.LookupEnv(k); !ok {
			t.Errorf("%s should be set but isn't", k)
		} else if actual != v {
			t.Errorf("Expected \"%s\" got %s", v, actual)
		}
	}
}

func TestSetCloudEnvWithNullChar(t *testing.T) {
	t.Cleanup(cleanEnv)

	opts := config.YAMLOpts{}
	if err := setCloudEnv("region_has_null_char", opts, false); err == nil {
		t.Errorf("setting environment should have failed but didn't: %v", err)
	}
}

func TestSetCloudEnv(t *testing.T) {
	t.Cleanup(cleanEnv)

	expected := map[string]string{
		"OS_CLOUD":        "hawaii",
		"OS_USERNAME":     "jdoe",
		"OS_PASSWORD":     "password",
		"OS_PROJECT_NAME": "Some Project",
		"OS_AUTH_URL":     "https://hi.example.com:5000/v3",
	}

	for k := range expected {
		if v, ok := os.LookupEnv(k); ok {
			t.Errorf("%s should not be set but is set to \"%s\"", k, v)
		}
	}

	opts := config.YAMLOpts{Directory: "testdata"}
	if err := setCloudEnv("foo", opts, false); err == nil {
		t.Error("\"foo\" is not a valid cloud, an error should have been returned by setCloudEnv()")
	}

	if err := setCloudEnv("hawaii", opts, false); err != nil {
		t.Errorf("error setting environment: %v", err)
	}

	for k, v := range expected {
		if actual, ok := os.LookupEnv(k); !ok {
			t.Errorf("%s should be set but isn't", k)
		} else if actual != v {
			t.Errorf("Expected \"%s\" got %s", v, actual)
		}
	}
}

func TestRun(t *testing.T) {
	t.Cleanup(cleanEnv)

	data := []testItem{
		{
			name: "no args",
			env:  map[string]string{},
			args: []string{"--os-config-dir", "testdata"},
			expected: `alberta
all_fields
all_from_profile
arizona
california
chicago
chicago_legacy
chicago_useprofile
disconnected_clouds
florida
florida_insecure
florida_secure
hawaii
nevada
newmexico
no_fields
philadelphia
philadelphia_complex
region_has_null_char
texas
virginia
yukon
`,
			expectedErr: nil,
		},
		{
			name:     "with cloud",
			env:      map[string]string{"MRACEK_SHELL": "/bin/bash"},
			args:     []string{"--os-config-dir", "testdata", "nevada"},
			expected: "Switching to cloud nevada\n",
		},
	}

	for _, item := range data {
		t.Run(item.name, testExecuteFunc(item))
	}
}

func TestAllCloudNamesWithInvalidCloudsYaml(t *testing.T) {
	workDir := t.TempDir()

	badYaml := []byte("------")

	cy := path.Join(workDir, "clouds.yaml")
	if err := os.WriteFile(cy, badYaml, 0o400); err != nil {
		t.Fatal(err)
	}

	opts := config.YAMLOpts{Directory: workDir}
	if _, err := opts.AllCloudNames(); err == nil {
		t.Error("opts.AllCloudNames() should fail with invalid clouds.yaml")
	}
}

func TestAllCloudNamesWithInvalidSecureYaml(t *testing.T) { //nolint:paralleltest
	workdir := t.TempDir()

	cy := path.Join(workdir, "clouds.yaml")
	// "---" is valid yaml
	if err := os.WriteFile(cy, []byte("---"), 0o400); err != nil {
		t.Fatal(err)
	}

	badYaml := []byte("-----")
	sy := path.Join(workdir, "secure.yaml")

	if err := os.WriteFile(sy, badYaml, 0o400); err != nil {
		t.Fatal(err)
	}

	opts := config.YAMLOpts{Directory: workdir}
	if _, err := opts.AllCloudNames(); err == nil {
		t.Error("allCloudNames() should fail with invalid secure.yaml")
	}
}

func TestShell(t *testing.T) {
	cmd := buildRootCommand("")
	buffer := bytes.NewBufferString("")
	cmd.SetOut(buffer)
	cmd.SetArgs([]string{"--os-config-dir", "testdata"})

	if err := cmd.Execute(); err != nil {
		t.Fatal(err)
	}

	if diff := cmp.Diff(viper.GetString("shell"), os.Getenv("SHELL")); diff != "" {
		t.Error(diff)
	}
}

func TestShellOverrideFromEnv(t *testing.T) {
	t.Setenv("MRACEK_SHELL", "foo")

	cmd := buildRootCommand("")
	buffer := bytes.NewBufferString("")
	cmd.SetOut(buffer)
	cmd.SetArgs([]string{"--os-config-dir", "testdata"})

	if err := cmd.Execute(); err != nil {
		t.Fatal(err)
	}

	if diff := cmp.Diff(viper.GetString("shell"), "foo"); diff != "" {
		t.Error(diff)
	}
}
