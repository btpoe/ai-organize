package main

// FileInfo represents information about a file
type FileInfo struct {
	Path         string                 `json:"path"`
	Name         string                 `json:"name"`
	Size         int64                  `json:"size"`
	IsDir        bool                   `json:"isDir"`
	Extension    string                 `json:"extension"`
	ModifiedTime string                 `json:"modifiedTime"`
	ContentHash  string                 `json:"contentHash"`
	MimeType     string                 `json:"mimeType"`
	Metadata     map[string]interface{} `json:"metadata,omitempty"`
	ParentDir    string                 `json:"parentDir"`
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
	Success        int      `json:"success"`
	Failed         int      `json:"failed"`
	FailedFiles    []string `json:"failedFiles"`
	CreatedFolders []string `json:"createdFolders"`
}

// DirectoryContext holds information about files in the same directory
type DirectoryContext struct {
	FilesByDir    map[string][]FileInfo
	FilesByHash   map[string][]FileInfo
	DominantTypes map[string]string // directory -> dominant file type
}
