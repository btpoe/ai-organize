# Claude AI Instructions for Delphi File Organizer

> **IMPORTANT**: When working with AI agents/tasks on this codebase, always include this file in the context/prompt to provide complete understanding of the application architecture, current capabilities, and limitations.

## Project Context

- **Project Name**: Delphi File Organizer
- **Purpose**: Desktop application that organizes files on a computer using intelligent categorization
- **Technology Stack**: Go (backend), Wails (framework), Svelte (frontend)
- **Issue Tracking**: This project is NOT linked to Linear or any issue tracking system. Task management is handled independently.

## Current Implementation

**Organization Strategy**: Multi-factor intelligent categorization
- Files are analyzed using **content hash (SHA-256)**, **MIME type**, **file extension**, and **directory context**
- **Duplicate detection**: Files with identical content hashes are grouped into a "Duplicates" folder
- **MIME-based categorization**: More accurate than extension-only (handles files with wrong/missing extensions)
- **Directory relationship awareness**: Files in the same directory with a dominant type (>50%) are kept together
- **Hierarchical organization**: Related files in the same directory are organized into subcategories (e.g., `Images/ProjectPhotos`)
- Files already in correct category folder are not moved

**Key Features**:
- ✅ Content hash analysis (SHA-256) for duplicate detection
- ✅ MIME type detection using `net/http` package
- ✅ Directory context-aware categorization
- ✅ Preserves relationships between files in the same directory
- ✅ Metadata structure ready for future extensions (EXIF, ID3, etc.)

**Current Limitations**:
- No EXIF/ID3/metadata extraction (structure in place, extraction not implemented)
- No ML-based content similarity analysis
- No date-based organization (could be added using existing metadata framework)

## Algorithm Details

The organization algorithm follows this decision tree:

1. **Duplicate Detection** (Priority 1):
   - Calculate SHA-256 hash of file content
   - If hash matches another file → Category: `Duplicates`

2. **MIME Type Validation** (Priority 2):
   - Detect MIME type from first 512 bytes
   - If MIME category differs from extension category → Use MIME category
   - Handles misnamed files (e.g., `.jpg` file that's actually a PDF)

3. **Directory Context Analysis** (Priority 3):
   - Analyze all files in the same directory
   - Calculate dominant type (must be >50% of files)
   - If file matches dominant type AND ≥3 related files → Category: `Type/ParentDirName`
   - Example: Photos in `/VacationPhotos` → `Images/VacationPhotos`

4. **Extension-Based Fallback** (Priority 4):
   - Use traditional extension-based categorization
   - Map file extension to predefined categories

## Code Architecture

The codebase has been organized into logical files for maintainability:

### Core Files

- **app.go**: Main application struct and directory analysis
  - `AnalyzeDirectory()` - Walks directory tree and proposes file moves
  - `SelectDirectory()` - Directory picker placeholder

- **types.go**: All data structures
  - `FileInfo` - File metadata with hash, MIME type, and metadata fields
  - `FileMove` - Proposed move operation
  - `AnalysisResult` - Analysis output with all proposed moves
  - `MoveResult` - Execution results with success/failure counts
  - `DirectoryContext` - Directory relationship data

- **hash.go**: Content hash calculation
  - `calculateFileHash()` - SHA-256 content hashing

- **mime.go**: MIME type detection and categorization
  - `detectMimeType()` - Reads first 512 bytes for type detection
  - `categorizeByMimeType()` - Maps MIME types to categories

- **categorize.go**: File categorization logic
  - `categorizeFile()` - Extension-based categorization (fallback)
  - `categorizeFileAdvanced()` - Multi-factor decision tree

- **context.go**: Directory context analysis
  - `buildDirectoryContext()` - Groups files and calculates dominant types

- **moves.go**: File move execution
  - `ExecuteMoves()` - Creates directories and moves files
  - `getUniqueFilePath()` - Conflict resolution with numeric suffixes

### Data Flow

```
User Input → AnalyzeDirectory() → buildDirectoryContext()
                ↓
            categorizeFileAdvanced() → categorizeFile()
                                    → categorizeByMimeType()
                ↓
            FileMove proposals
                ↓
            User Selection → ExecuteMoves() → MoveResult
```

## File Categories

The application recognizes these categories:

- **Images**: 11 extensions (.jpg, .png, .gif, etc.)
- **Documents**: 15 extensions (.pdf, .docx, .txt, etc.)
- **Videos**: 10 extensions (.mp4, .mov, .mkv, etc.)
- **Audio**: 8 extensions (.mp3, .wav, .flac, etc.)
- **Archives**: 8 extensions (.zip, .rar, .7z, etc.)
- **Code**: 23 extensions (.go, .py, .js, etc.)
- **Applications**: 6 extensions (.exe, .app, .deb, etc.)
- **Duplicates**: Files with matching content hashes
- **Other**: Catch-all for unrecognized types

## Potential Enhancements

Future enhancements could include:

1. **Metadata Extraction** (structure ready):
   - Parse EXIF data from images (creation date, camera, location)
   - Extract ID3 tags from audio files (artist, album, genre)
   - Read PDF metadata (author, creation date, title)
   - Parse video metadata (resolution, codec, duration)
   - ✅ Struct already has `Metadata map[string]interface{}`

2. **Date-Based Organization**:
   - Organize photos by date taken (EXIF DateTimeOriginal)
   - Group documents by creation/modification date
   - Automatic archival of old files

3. **Content Similarity**:
   - ML-based image similarity clustering
   - Text document topic analysis
   - Perceptual hashing for near-duplicate detection

4. **Advanced Duplicate Handling**:
   - Keep highest quality version
   - Merge metadata from duplicates
   - Smart deduplication based on file age

5. **Custom Rules Engine**:
   - User-defined categorization rules
   - Regex-based file matching
   - Conditional logic for complex scenarios

## Development Guidelines

### When Adding New Features

1. **Maintain Separation of Concerns**: Add new functionality in the appropriate file
2. **Update Types**: If adding new data, update `types.go` first
3. **Follow Existing Patterns**: Use similar structure to existing categorization logic
4. **Document Changes**: Update both README.md (user-facing) and CLAUDE.md (AI-facing)
5. **Test Edge Cases**: Consider empty directories, missing files, permission errors

### Common Tasks

- **Add new file extension**: Update the appropriate map in `categorize.go`
- **Add new category**: Update both `categorizeFile()` and `categorizeByMimeType()`
- **Modify organization logic**: Update `categorizeFileAdvanced()` in `categorize.go`
- **Add metadata extraction**: Extend the FileInfo population in `app.go` `AnalyzeDirectory()`
- **Change duplicate detection**: Modify logic in `categorizeFileAdvanced()` and `buildDirectoryContext()`

## Testing Considerations

- **Large Directories**: Algorithm should handle 1000+ files efficiently
- **Duplicate Detection**: Hash calculation can be slow for large files
- **MIME Detection**: Limited to first 512 bytes, may misidentify some files
- **Permission Errors**: Handle read/write permission failures gracefully
- **Cross-Platform**: File paths work on macOS, Windows, and Linux

## Safety Features

- **Dry run first**: Always shows proposed changes before executing
- **Conflict handling**: Automatically renames files if destination exists (adds numeric suffix)
- **No destructive operations**: Only moves files, never deletes
- **Skip hidden files**: Ignores files and folders starting with "." to avoid system files
- **Permission safety**: Creates directories with 0755 permissions
