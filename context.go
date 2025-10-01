package main

import (
	"path/filepath"
)

// buildDirectoryContext analyzes files to understand directory relationships
func buildDirectoryContext(files []FileInfo) DirectoryContext {
	ctx := DirectoryContext{
		FilesByDir:    make(map[string][]FileInfo),
		FilesByHash:   make(map[string][]FileInfo),
		DominantTypes: make(map[string]string),
	}

	// Group files by directory
	for _, file := range files {
		dir := filepath.Dir(file.Path)
		ctx.FilesByDir[dir] = append(ctx.FilesByDir[dir], file)

		// Group by content hash for duplicate detection
		if file.ContentHash != "" {
			ctx.FilesByHash[file.ContentHash] = append(ctx.FilesByHash[file.ContentHash], file)
		}
	}

	// Determine dominant file type in each directory
	for dir, dirFiles := range ctx.FilesByDir {
		typeCounts := make(map[string]int)
		for _, f := range dirFiles {
			baseCategory, _ := categorizeFile(f)
			typeCounts[baseCategory]++
		}

		// Find dominant type
		maxCount := 0
		dominantType := ""
		for fileType, count := range typeCounts {
			if count > maxCount {
				maxCount = count
				dominantType = fileType
			}
		}

		// Only set dominant type if it represents majority (>50%)
		if maxCount > len(dirFiles)/2 {
			ctx.DominantTypes[dir] = dominantType
		}
	}

	return ctx
}
