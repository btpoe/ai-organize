package main

import (
	"io"
	"net/http"
	"os"
	"strings"
)

// detectMimeType detects the MIME type of a file
func detectMimeType(filePath string) (string, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return "", err
	}
	defer file.Close()

	// Read first 512 bytes for MIME detection
	buffer := make([]byte, 512)
	n, err := file.Read(buffer)
	if err != nil && err != io.EOF {
		return "", err
	}

	// Detect content type
	mimeType := http.DetectContentType(buffer[:n])
	return mimeType, nil
}

// categorizeByMimeType provides MIME-based categorization
func categorizeByMimeType(mimeType string) string {
	switch {
	case strings.HasPrefix(mimeType, "image/"):
		return "Images"
	case strings.HasPrefix(mimeType, "video/"):
		return "Videos"
	case strings.HasPrefix(mimeType, "audio/"):
		return "Audio"
	case strings.HasPrefix(mimeType, "text/"):
		return "Documents"
	case mimeType == "application/pdf":
		return "Documents"
	case strings.Contains(mimeType, "word") || strings.Contains(mimeType, "document"):
		return "Documents"
	case strings.Contains(mimeType, "sheet") || strings.Contains(mimeType, "excel"):
		return "Documents"
	case strings.Contains(mimeType, "presentation") || strings.Contains(mimeType, "powerpoint"):
		return "Documents"
	case strings.Contains(mimeType, "zip") || strings.Contains(mimeType, "compressed"):
		return "Archives"
	case strings.Contains(mimeType, "executable") || strings.Contains(mimeType, "application/x-"):
		return "Applications"
	}
	return ""
}
