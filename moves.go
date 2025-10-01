package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// ExecuteMoves executes the approved file moves
func (a *App) ExecuteMoves(moves []FileMove) MoveResult {
	result := MoveResult{
		FailedFiles:    []string{},
		CreatedFolders: []string{},
	}

	// Track created directories to avoid duplicates
	createdDirs := make(map[string]bool)

	for _, move := range moves {
		// Create destination directory if it doesn't exist
		destDir := filepath.Dir(move.DestinationPath)
		if _, err := os.Stat(destDir); os.IsNotExist(err) {
			err := os.MkdirAll(destDir, 0755)
			if err != nil {
				result.Failed++
				result.FailedFiles = append(result.FailedFiles,
					fmt.Sprintf("%s: failed to create directory - %v", move.FileName, err))
				continue
			}
			if !createdDirs[destDir] {
				result.CreatedFolders = append(result.CreatedFolders, destDir)
				createdDirs[destDir] = true
			}
		}

		// Check if destination file already exists
		if _, err := os.Stat(move.DestinationPath); err == nil {
			// File exists, add number suffix
			move.DestinationPath = getUniqueFilePath(move.DestinationPath)
		}

		// Move the file
		err := os.Rename(move.SourcePath, move.DestinationPath)
		if err != nil {
			result.Failed++
			result.FailedFiles = append(result.FailedFiles,
				fmt.Sprintf("%s: %v", move.FileName, err))
		} else {
			result.Success++
		}
	}

	return result
}

// getUniqueFilePath generates a unique file path by adding a number suffix
func getUniqueFilePath(path string) string {
	ext := filepath.Ext(path)
	base := strings.TrimSuffix(path, ext)

	counter := 1
	newPath := fmt.Sprintf("%s_%d%s", base, counter, ext)

	for {
		if _, err := os.Stat(newPath); os.IsNotExist(err) {
			return newPath
		}
		counter++
		newPath = fmt.Sprintf("%s_%d%s", base, counter, ext)
	}
}
