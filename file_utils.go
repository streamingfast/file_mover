package file_mover

import (
	"fmt"
	"io"
	"os"
)

func moveFile(sourcePath, destPath string) error {
	err := copyFile(sourcePath, destPath)
	if err != nil {
		return fmt.Errorf("copying file: %w", err)
	}

	err = deleteFile(sourcePath)
	if err != nil {
		return fmt.Errorf("removing file: %w", err)
	}
	return nil
}

func copyFile(sourcePath, destPath string) error {
	inputFile, err := os.Open(sourcePath)
	if err != nil {
		return fmt.Errorf("open source file: %s", err)
	}
	defer inputFile.Close()

	outputFile, err := os.Create(destPath)
	if err != nil {
		return fmt.Errorf("open dest file: %s", err)
	}
	defer outputFile.Close()

	_, err = io.Copy(outputFile, inputFile)
	if err != nil {
		return fmt.Errorf("writing to output: %s", err)
	}

	return nil
}

func deleteFile(sourcePath string) error {
	// The copy was successful, so now delete the original file
	err := os.Remove(sourcePath)
	if err != nil {
		return fmt.Errorf("removing file: %s", err)
	}
	return nil
}
