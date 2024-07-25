package utils

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
)

func Edit(tempDir string, src *string) ([]byte, error) {
	editor := os.Getenv("EDITOR")
	if editor == "" {
		editor = "vim"
	}

	editPath, err := ensureEditFile(tempDir, src)
	if err != nil {
		return nil, err
	}

	cmd := exec.Command(editor, editPath)
	cmd.Stdout = os.Stderr
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin

	err = cmd.Run()
	if err != nil {
		return nil, fmt.Errorf("run edit command (using editor %q): %w", editor, err)
	}

	data, err := os.ReadFile(editPath)
	if err != nil {
		return nil, fmt.Errorf("read edit temp file: %w", err)
	}

	err = os.Remove(editPath)
	if err != nil {
		return nil, fmt.Errorf("remove edit temp file: %w", err)
	}

	return data, nil
}

func ensureEditFile(tempDir string, src *string) (string, error) {
	var data []byte
	if src != nil {
		var err error
		data, err = os.ReadFile(*src)
		if err != nil {
			return "", err
		}
	}

	editPath := filepath.Join(tempDir, "edit.yaml")
	err := WriteFile(editPath, data)
	if err != nil {
		return "", fmt.Errorf("write edit temp file: %w", err)
	}

	return editPath, nil
}
