package utils

import (
	"fmt"
	"os"
	"path/filepath"
)

const (
	recommandedDirMode  = 0755
	recommandedFileMode = 0644
)

// WriteFile will create directory and write content to file
func WriteFile(path string, data []byte) error {
	dir := filepath.Dir(path)
	err := EnsureDir(dir)
	if err != nil {
		return err
	}
	return os.WriteFile(path, data, recommandedFileMode)
}

func EnsureDir(dir string) error {
	info, err := os.Stat(dir)
	if err != nil {
		if os.IsNotExist(err) {
			err = os.MkdirAll(dir, recommandedDirMode)
			if err != nil {
				return err
			}
			return nil
		} else {
			return err
		}
	}

	if !info.IsDir() {
		return fmt.Errorf("bad path, %q is not a directory", dir)
	}

	return nil
}

func EnsureFile(path string) error {
	info, err := os.Stat(path)
	if err != nil {
		return err
	}

	if info.IsDir() {
		return fmt.Errorf("bad path, %q is a directory", path)
	}

	return nil
}
