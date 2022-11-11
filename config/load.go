package config

import (
	"fmt"
	"os"
	"os/user"
	"path"
	"path/filepath"
	"sort"

	"github.com/gophercloud/utils/openstack/clientconfig"
	"github.com/spf13/viper"
	"gopkg.in/yaml.v3"
)

// fileExists checks for the existence of a file at a given location.
func fileExists(filename string) bool {
	if _, err := os.Stat(filename); err == nil {
		return true
	}

	return false
}

// FindAndReadYAML looks for filename in home directory if dir is empty, otherwise looks in specified directory.
func FindAndReadYAML(dir, filename string) (string, []byte, error) {
	if dir == "" {
		currentUser, err := user.Current()
		if err != nil {
			return "", nil, fmt.Errorf("could not get current user: %w", err)
		}

		dir = currentUser.HomeDir
		if dir == "" {
			return "", nil, fmt.Errorf("no home directory found for user %s", currentUser)
		}

		dir = filepath.Join(dir, ".config/openstack/")
	}

	fullPath := filepath.Join(dir, filename)

	if ok := fileExists(fullPath); ok {
		content, err := os.ReadFile(fullPath)
		if err != nil {
			return "", nil, fmt.Errorf("%w", err)
		}

		return fullPath, content, nil
	}

	return "", nil, fmt.Errorf("file %s does not exist", fullPath)
}

type Clouds struct {
	Clouds       map[string]clientconfig.Cloud `yaml:"clouds" json:"clouds"`
	PublicClouds map[string]clientconfig.Cloud `yaml:"public-clouds" json:"public-clouds"` //nolint:tagliatelle
}

func LoadYAML(dir, filename string) (Clouds, error) {
	_, content, err := FindAndReadYAML(dir, filename)
	if err != nil {
		return Clouds{}, err
	}

	clouds := Clouds{}
	if err = yaml.Unmarshal(content, &clouds); err != nil {
		return Clouds{}, fmt.Errorf("failed to unmarshal yaml: %w", err)
	}

	return clouds, nil
}

// YAMLOpts implements gophercloud util's YAMLOptsBuilder interface.
type YAMLOpts struct {
	Directory string
}

// AllCloudNames returns all the possible cloud names.
func (opts YAMLOpts) AllCloudNames() ([]string, error) {
	clouds, err := opts.LoadCloudsYAML()
	if err != nil {
		return nil, fmt.Errorf("failed to load clouds.yaml: %w", err)
	}

	sclouds, err := opts.LoadSecureCloudsYAML()
	if err != nil {
		return nil, fmt.Errorf("failed to load secure.yaml: %w", err)
	}

	names := make(map[string]bool, len(clouds)+len(sclouds))

	for c := range clouds {
		names[c] = true
	}

	for sc := range sclouds {
		names[sc] = true
	}

	all := make([]string, 0, len(names))
	for k := range names {
		all = append(all, k)
	}

	sort.Strings(all)

	return all, nil
}

// AllProfileNames returns all the possible profile names.
func (opts YAMLOpts) AllProfileNames() ([]string, error) {
	profiles, err := opts.LoadPublicCloudsYAML()
	if err != nil {
		return nil, fmt.Errorf("failed to load clouds-public.yaml: %w", err)
	}

	names := make([]string, 0, len(profiles))
	for name := range profiles {
		names = append(names, name)
	}

	return names, nil
}

// LoadCloudsYAML loads ~/.config/openstack/clouds.yaml.
func (opts YAMLOpts) LoadCloudsYAML() (map[string]clientconfig.Cloud, error) {
	clouds, err := LoadYAML(opts.Directory, "clouds.yaml")
	if err != nil {
		return nil, err
	}

	return clouds.Clouds, err
}

// LoadSecureCloudsYAML loads ~/.config/openstack/secure.yaml.
func (opts YAMLOpts) LoadSecureCloudsYAML() (map[string]clientconfig.Cloud, error) {
	clouds, err := LoadYAML(opts.Directory, "secure.yaml")
	if err != nil {
		return nil, err
	}

	return clouds.Clouds, err
}

// LoadPublicCloudsYAML loads ~/.config/openstack/clouds-public.yaml.
func (opts YAMLOpts) LoadPublicCloudsYAML() (map[string]clientconfig.Cloud, error) {
	clouds, err := LoadYAML(opts.Directory, "clouds-public.yaml")
	if err != nil {
		return nil, err
	}

	return clouds.PublicClouds, err
}

// LoadAndCheckOSConfigfile loads and returns all the clouds found the specified filename. If
// cloudName is different from "" then cloudName is searched for in the clouds and an error
// is returned if cloudName is found.
func LoadAndCheckOSConfigfile(
	filename string,
	loader func() (map[string]clientconfig.Cloud, error),
	cloudName string,
) (map[string]clientconfig.Cloud, error) {
	yamlPath := path.Join(viper.GetString("os-config-dir"), filename)

	clouds, err := loader()
	if err != nil {
		if err.Error() != fmt.Sprintf("file %s does not exist", yamlPath) {
			return nil, fmt.Errorf("%w", err)
		}

		clouds = make(map[string]clientconfig.Cloud)
	}

	if cloudName != "" {
		if _, present := clouds[cloudName]; present {
			return nil, fmt.Errorf("cloud %s already exists in %s", cloudName, filename)
		}
	}

	return clouds, nil
}
