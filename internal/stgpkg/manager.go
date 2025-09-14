package stgpkg

import (
	"fmt"
	"os"
)

func Write(data []byte, path string) error {
	err := os.WriteFile(path, data, 0644)

	if err != nil {
		return fmt.Errorf("file write failed: %w", err)
	}

	return nil
}

func Read(path string) ([]byte, error) {
	data, err := os.ReadFile(path)

	if err != nil {
		return data, fmt.Errorf("file read failed: %w", err)
	}

	return data, nil
}
