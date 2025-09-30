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

// FileInfo represents information about a file
type FileInfo struct {
	Path         string `json:"path"`
	Name         string `json:"name"`
	Size         int64  `json:"size"`
	IsDir        bool   `json:"isDir"`
	Extension    string `json:"extension"`
	ModifiedTime string `json:"modifiedTime"`
}

// FileMove represents a proposed file move operation
type FileMove struct {
	SourcePath      string `json:"sourcePath"`
	DestinationPath string `json:"destinationPath"`
	FileName        string `json:"fileName"`
	Reason          string `json:"reason"`
	Category        string `json:"category"`
}

// AnalysisResult represents the result of analyzing a directory
type AnalysisResult struct {
	TotalFiles    int        `json:"totalFiles"`
	ProposedMoves []FileMove `json:"proposedMoves"`
	Error         string     `json:"error,omitempty"`
}

// MoveResult represents the result of executing file moves
type MoveResult struct {
	Success       int      `json:"success"`
	Failed        int      `json:"failed"`
	FailedFiles   []string `json:"failedFiles"`
	CreatedFolders []string `json:"createdFolders"`
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

			files = append(files, FileInfo{
				Path:         path,
				Name:         d.Name(),
				Size:         info.Size(),
				IsDir:        false,
				Extension:    strings.ToLower(filepath.Ext(d.Name())),
				ModifiedTime: info.ModTime().Format("2006-01-02 15:04:05"),
			})
		}

		return nil
	})

	if err != nil {
		result.Error = fmt.Sprintf("Error walking directory: %v", err)
		return result
	}

	result.TotalFiles = len(files)

	// Organize files by category
	for _, file := range files {
		category, reason := categorizeFile(file)

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

// categorizeFile determines the category and reason for a file
func categorizeFile(file FileInfo) (category string, reason string) {
	ext := file.Extension

	// Image files
	imageExts := map[string]bool{
		".jpg": true, ".jpeg": true, ".png": true, ".gif": true,
		".bmp": true, ".svg": true, ".webp": true, ".ico": true,
		".heic": true, ".raw": true, ".tiff": true,
	}
	if imageExts[ext] {
		return "Images", fmt.Sprintf("Image file (%s)", ext)
	}

	// Document files
	docExts := map[string]bool{
		".pdf": true, ".doc": true, ".docx": true, ".txt": true,
		".rtf": true, ".odt": true, ".pages": true, ".tex": true,
		".md": true, ".csv": true, ".xls": true, ".xlsx": true,
		".ppt": true, ".pptx": true, ".key": true,
	}
	if docExts[ext] {
		return "Documents", fmt.Sprintf("Document file (%s)", ext)
	}

	// Video files
	videoExts := map[string]bool{
		".mp4": true, ".avi": true, ".mov": true, ".mkv": true,
		".wmv": true, ".flv": true, ".webm": true, ".m4v": true,
		".mpg": true, ".mpeg": true,
	}
	if videoExts[ext] {
		return "Videos", fmt.Sprintf("Video file (%s)", ext)
	}

	// Audio files
	audioExts := map[string]bool{
		".mp3": true, ".wav": true, ".flac": true, ".aac": true,
		".m4a": true, ".wma": true, ".ogg": true, ".opus": true,
	}
	if audioExts[ext] {
		return "Audio", fmt.Sprintf("Audio file (%s)", ext)
	}

	// Archive files
	archiveExts := map[string]bool{
		".zip": true, ".rar": true, ".7z": true, ".tar": true,
		".gz": true, ".bz2": true, ".xz": true, ".dmg": true,
	}
	if archiveExts[ext] {
		return "Archives", fmt.Sprintf("Archive file (%s)", ext)
	}

	// Code files
	codeExts := map[string]bool{
		".go": true, ".py": true, ".js": true, ".ts": true,
		".java": true, ".c": true, ".cpp": true, ".h": true,
		".cs": true, ".rb": true, ".php": true, ".swift": true,
		".rs": true, ".kt": true, ".scala": true, ".sh": true,
		".html": true, ".css": true, ".json": true, ".xml": true,
		".yaml": true, ".yml": true, ".sql": true,
	}
	if codeExts[ext] {
		return "Code", fmt.Sprintf("Code file (%s)", ext)
	}

	// Executable files
	execExts := map[string]bool{
		".exe": true, ".app": true, ".deb": true, ".rpm": true,
		".apk": true, ".msi": true,
	}
	if execExts[ext] {
		return "Applications", fmt.Sprintf("Application file (%s)", ext)
	}

	// Default category for unknown files
	return "Other", fmt.Sprintf("Unknown file type (%s)", ext)
}

// ExecuteMoves executes the approved file moves
func (a *App) ExecuteMoves(moves []FileMove) MoveResult {
	result := MoveResult{
		FailedFiles:   []string{},
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

// SelectDirectory opens a directory picker dialog
func (a *App) SelectDirectory() string {
	// This is a placeholder - Wails runtime provides native dialog
	// In real implementation, use runtime.OpenDirectoryDialog
	return ""
}