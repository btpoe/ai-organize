package main

import (
	"context"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
)

// App struct
type App struct {
	ctx context.Context
}

// NewApp creates a new App application struct
func NewApp() *App {
	return &App{}
}

// startup is called when the app starts
func (a *App) startup(ctx context.Context) {
	a.ctx = ctx
}

// AnalyzeDirectory analyzes a directory and proposes file organization
func (a *App) AnalyzeDirectory(dirPath string) AnalysisResult {
	result := AnalysisResult{
		ProposedMoves: []FileMove{},
	}

	// Validate directory exists
	if _, err := os.Stat(dirPath); os.IsNotExist(err) {
		result.Error = "Directory does not exist"
		return result
	}

	var files []FileInfo

	// Walk through directory
	err := filepath.WalkDir(dirPath, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		// Skip the root directory itself
		if path == dirPath {
			return nil
		}

		// Skip hidden files and directories
		if strings.HasPrefix(d.Name(), ".") {
			if d.IsDir() {
				return filepath.SkipDir
			}
			return nil
		}

		if !d.IsDir() {
			info, err := d.Info()
			if err != nil {
				return nil
			}

			fileInfo := FileInfo{
				Path:         path,
				Name:         d.Name(),
				Size:         info.Size(),
				IsDir:        false,
				Extension:    strings.ToLower(filepath.Ext(d.Name())),
				ModifiedTime: info.ModTime().Format("2006-01-02 15:04:05"),
				ParentDir:    filepath.Base(filepath.Dir(path)),
				Metadata:     make(map[string]interface{}),
			}

			// Calculate content hash
			if hash, err := calculateFileHash(path); err == nil {
				fileInfo.ContentHash = hash
			}

			// Detect MIME type
			if mimeType, err := detectMimeType(path); err == nil {
				fileInfo.MimeType = mimeType
			}

			files = append(files, fileInfo)
		}

		return nil
	})

	if err != nil {
		result.Error = fmt.Sprintf("Error walking directory: %v", err)
		return result
	}

	result.TotalFiles = len(files)

	// Build directory context map for smarter categorization
	dirContext := buildDirectoryContext(files)

	// Organize files by category
	for _, file := range files {
		category, reason := categorizeFileAdvanced(file, dirContext)

		// Only propose a move if the file isn't already in the right place
		currentDir := filepath.Base(filepath.Dir(file.Path))
		if currentDir != category {
			destPath := filepath.Join(dirPath, category, file.Name)

			result.ProposedMoves = append(result.ProposedMoves, FileMove{
				SourcePath:      file.Path,
				DestinationPath: destPath,
				FileName:        file.Name,
				Reason:          reason,
				Category:        category,
			})
		}
	}

	return result
}

// SelectDirectory opens a directory picker dialog
func (a *App) SelectDirectory() string {
	// This is a placeholder - Wails runtime provides native dialog
	// In real implementation, use runtime.OpenDirectoryDialog
	return ""
}
