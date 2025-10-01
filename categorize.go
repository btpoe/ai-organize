package main

import (
	"fmt"
	"path/filepath"
)

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

// categorizeFileAdvanced determines category using multiple factors
func categorizeFileAdvanced(file FileInfo, dirContext DirectoryContext) (category string, reason string) {
	// First, get the basic category from extension
	baseCategory, baseReason := categorizeFile(file)

	// Check if file is a duplicate (same content hash)
	if file.ContentHash != "" {
		if dupes, exists := dirContext.FilesByHash[file.ContentHash]; exists && len(dupes) > 1 {
			// This is a duplicate file
			reason = fmt.Sprintf("Duplicate file (hash: %s...)", file.ContentHash[:8])
			category = "Duplicates"
			return
		}
	}

	// Check MIME type for better accuracy
	if file.MimeType != "" {
		if mimeCategory := categorizeByMimeType(file.MimeType); mimeCategory != "" {
			// MIME type provides more accurate categorization
			if mimeCategory != baseCategory {
				reason = fmt.Sprintf("%s (MIME: %s)", baseReason, file.MimeType)
				category = mimeCategory
				return
			}
		}
	}

	// Check if files in the same directory are strongly related
	dir := filepath.Dir(file.Path)
	if dominantType, exists := dirContext.DominantTypes[dir]; exists {
		// If this file's type matches the dominant type in its directory, keep it there
		if baseCategory == dominantType {
			filesInDir := dirContext.FilesByDir[dir]
			if len(filesInDir) >= 3 { // At least 3 related files
				reason = fmt.Sprintf("Related to %d other %s files in same directory", len(filesInDir)-1, dominantType)
				category = fmt.Sprintf("%s/%s", dominantType, file.ParentDir)
				return
			}
		}
	}

	// Default to base categorization
	return baseCategory, baseReason
}
