package stgpkg

import (
	"fmt"
	"os"
)

var BOM = []byte{0xEF, 0xBB, 0xBF}

func Write(data []byte, path string, addBom bool) error {
	if addBom {
		data = append(BOM, data...)
	}

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

	if len(data) >= 3 && data[0] == BOM[0] && data[1] == BOM[1] && data[2] == BOM[2] {
		data = data[3:]
	}

	return data, nil
}
