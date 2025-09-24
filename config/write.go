package config

import (
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path"
	"path/filepath"
	"sort"
	"strconv"
	"time"

	"github.com/spf13/cobra"
)

const (
	BackupsToKeep = 4
)

// backupFile backs up the file at path while appending the unix timestamp. Keeps at most BackupsToKeep backup files.
func backupFile(path string) {
	input, err := os.ReadFile(path)
	if errors.Is(err, fs.ErrNotExist) {
		return
	}

	cobra.CheckErr(err)

	err = os.WriteFile(path+"-"+strconv.FormatInt(time.Now().Unix(), 10), input, 0o600)
	cobra.CheckErr(err)

	backups, err := filepath.Glob(path + "-*")
	cobra.CheckErr(err)

	sort.Strings(backups)

	// cleanup old backups while keeping the oldest backup
	// 1 2 3 4 5   we want to keep 1 and delete 2 (len 5, max 4)
	// 1 2 3 4 5 6   we want to keep 1 and delete 2 and 3 (len 6, max 4)
	// 1 2 3 4 5 6 7  we want to keep 1 and delete 2, 3 and 4 (len 7, max 4)
	if len(backups) > BackupsToKeep {
		for i := 1; i < len(backups)-(BackupsToKeep-1); i++ {
			cobra.CheckErr(os.Remove(backups[i]))
		}
	}
}

func WriteOSConfig(directory string, clouds, secure, public []byte) (err error) {
	defer func() {
		if err != nil {
			err = fmt.Errorf("failed to write configuration: %w", err)
		}
	}()

	if err := os.MkdirAll(directory, 0o755); err != nil {
		return err
	}

	if clouds != nil {
		f := path.Join(directory, "clouds.yaml")

		backupFile(f)

		err := os.WriteFile(f, clouds, 0o600)
		if err != nil {
			return err
		}
	}

	if secure != nil {
		f := path.Join(directory, "secure.yaml")

		backupFile(f)

		err := os.WriteFile(f, secure, 0o600)
		if err != nil {
			return err
		}
	}

	if public != nil {
		f := path.Join(directory, "clouds-public.yaml")

		backupFile(f)

		err := os.WriteFile(f, public, 0o600)
		if err != nil {
			return err
		}
	}

	return nil
}
