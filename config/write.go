package config

import (
	"os"
	"path"
)

func WriteOSConfig(directory string, clouds, secure, public []byte) error {
	if clouds != nil {
		fn := path.Join(directory, "clouds.yaml")

		err := os.WriteFile(fn, clouds, 0o600) //nolint:gomnd
		if err != nil {
			return err
		}
	}

	if secure != nil {
		fn := path.Join(directory, "secure.yaml")

		err := os.WriteFile(fn, secure, 0o600) //nolint:gomnd
		if err != nil {
			return err
		}
	}

	if public != nil {
		fn := path.Join(directory, "clouds-public.yaml")

		err := os.WriteFile(fn, public, 0o600) //nolint:gomnd
		if err != nil {
			return err
		}
	}

	return nil
}
