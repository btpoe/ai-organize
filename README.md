# Delphi File Organizer

A desktop application built with Go and Wails that automatically analyzes and organizes files on a computer into logical folders based on file type and metadata.

## Features

- **Automatic File Analysis**: Scans directories and categorizes files by type
- **Smart Organization**: Groups files into categories (Images, Documents, Videos, Audio, Code, Archives, Applications, etc.)
- **Dry Run Preview**: Shows proposed changes before executing
- **Interactive UI**: Review, edit, and approve file moves through an intuitive interface
- **Bulk Operations**: Select multiple files or entire categories at once
- **Safe Execution**: Creates necessary folders and handles file conflicts automatically

## Categories

The application organizes files into these categories:

- **Images**: jpg, jpeg, png, gif, bmp, svg, webp, ico, heic, raw, tiff
- **Documents**: pdf, doc, docx, txt, rtf, odt, pages, tex, md, csv, xls, xlsx, ppt, pptx, key
- **Videos**: mp4, avi, mov, mkv, wmv, flv, webm, m4v, mpg, mpeg
- **Audio**: mp3, wav, flac, aac, m4a, wma, ogg, opus
- **Archives**: zip, rar, 7z, tar, gz, bz2, xz, dmg
- **Code**: go, py, js, ts, java, c, cpp, h, cs, rb, php, swift, rs, kt, scala, sh, html, css, json, xml, yaml, yml, sql
- **Applications**: exe, app, deb, rpm, apk, msi
- **Other**: Files that don't match any category

## Prerequisites

- Go 1.21 or later
- Node.js 18 or later
- Wails v2 CLI
- macOS: Full Disk Access or permission to access folders like Downloads, Documents, Desktop

## Installation

1. Install Wails CLI:
```bash
go install github.com/wailsapp/wails/v2/cmd/wails@latest
```

2. Install frontend dependencies:
```bash
cd frontend
npm install
cd ..
```

3. Download Go dependencies:
```bash
go mod download
```

## Development

Run in development mode with hot reload:

```bash
wails dev
```

## Building

Build for your current platform:

```bash
wails build
```

The compiled application will be in the `build/bin` directory.

### macOS Permissions

On macOS, the application requires file access permissions. When you first run the app and try to access protected folders (Downloads, Documents, Desktop, etc.), macOS will prompt you to grant permission. Alternatively, you can manually grant Full Disk Access:

1. Open **System Settings** > **Privacy & Security** > **Full Disk Access**
2. Click the **+** button
3. Navigate to and select the **File Organizer.app**
4. Restart the application

### Build for specific platforms:

```bash
# macOS
wails build -platform darwin/universal

# Windows
wails build -platform windows/amd64

# Linux
wails build -platform linux/amd64
```

## Usage

1. Launch the application
2. Enter the path to the directory you want to organize
3. Click "Analyze Directory" to scan and categorize files
4. Review the proposed file moves:
   - Check/uncheck individual files or entire categories
   - Edit destination paths if needed
5. Click "Execute Moves" to perform the organization
6. Review the results showing successful and failed operations

## Project Structure

```
.
├── main.go              # Application entry point
├── app.go               # App struct and AnalyzeDirectory method
├── types.go             # Data structures (FileInfo, FileMove, etc.)
├── hash.go              # Content hash calculation (SHA-256)
├── mime.go              # MIME type detection and categorization
├── categorize.go        # File categorization logic
├── context.go           # Directory context analysis
├── moves.go             # File move execution and conflict resolution
├── go.mod               # Go dependencies
├── wails.json           # Wails configuration
├── CLAUDE.md            # AI development guidelines and technical docs
├── frontend/
│   ├── src/
│   │   ├── App.svelte   # Main UI component
│   │   └── main.js      # Frontend entry point
│   ├── index.html       # HTML template
│   ├── vite.config.js   # Vite configuration
│   └── package.json     # Frontend dependencies
└── README.md
```

## How It Works

### Backend (Go)

The application uses a Go backend with the following key components:

#### Core Data Structures

- **FileInfo**: Represents file metadata including path, name, size, extension, modified time
- **FileMove**: Represents a proposed file move operation with source, destination, category, and reason
- **AnalysisResult**: Contains total file count and array of proposed moves
- **MoveResult**: Contains execution results with success/failure counts and created folders

#### Core Methods

- **AnalyzeDirectory(dirPath string) AnalysisResult**:
  - Walks through the specified directory using `filepath.WalkDir`
  - Skips hidden files/folders (starting with ".")
  - Categorizes each file by extension using `categorizeFile()`
  - Only proposes moves for files not already in the correct category folder
  - Returns analysis with all proposed moves

- **ExecuteMoves(moves []FileMove) MoveResult**:
  - Creates destination directories with `os.MkdirAll` (0755 permissions)
  - Checks for existing files and uses `getUniqueFilePath()` to avoid conflicts
  - Executes file moves using `os.Rename`
  - Tracks created folders and failed operations
  - Returns detailed results

- **categorizeFile(file FileInfo) (category, reason)**:
  - Determines category based on file extension
  - Uses predefined extension maps for each category
  - Returns category name and human-readable reason
  - Defaults to "Other" for unknown types

- **getUniqueFilePath(path string)**:
  - Generates unique file paths by appending numeric suffixes (_1, _2, etc.)
  - Prevents file overwriting during moves

#### File Categorization Logic

Files are categorized purely by extension into these predefined categories:
- Images: 11 extensions (.jpg, .png, .gif, etc.)
- Documents: 15 extensions (.pdf, .docx, .txt, etc.)
- Videos: 10 extensions (.mp4, .mov, .mkv, etc.)
- Audio: 8 extensions (.mp3, .wav, .flac, etc.)
- Archives: 8 extensions (.zip, .rar, .7z, etc.)
- Code: 23 extensions (.go, .py, .js, etc.)
- Applications: 6 extensions (.exe, .app, .deb, etc.)
- Other: Catch-all for unrecognized extensions

### Frontend (Svelte)

- Clean, modern UI with gradient design
- Real-time file selection and editing
- Grouped view by category
- Detailed move information (source, destination, reason)
- Execution results with success/failure reporting
- Calls Go backend methods via Wails bridge

### Application Flow

1. User inputs directory path
2. Frontend calls `AnalyzeDirectory(path)` via Wails bridge
3. Backend walks directory tree, skipping hidden files
4. Each file is categorized by extension
5. Proposed moves returned to frontend (only files not in correct folder)
6. User reviews/selects files to move in UI
7. Frontend calls `ExecuteMoves(selectedMoves)`
8. Backend creates folders, handles conflicts, executes moves
9. Results displayed showing success/failure for each operation

## Safety Features

- **Dry run first**: Always shows proposed changes before executing
- **Conflict handling**: Automatically renames files if destination exists (adds numeric suffix)
- **Error reporting**: Detailed feedback on any failed operations
- **No destructive operations**: Only moves files, never deletes
- **Skip hidden files**: Ignores files and folders starting with "." to avoid system files
- **Permission safety**: Creates directories with 0755 permissions

## Advanced Features

### Multi-Factor Categorization

The application uses an intelligent decision tree that considers:
- **Content hash (SHA-256)** for duplicate detection
- **MIME type** for accurate file type identification
- **Directory context** to preserve relationships between related files
- **File extension** as a fallback

### Smart Organization

- **Duplicate Detection**: Files with identical content are automatically grouped
- **Hierarchical Organization**: Related files stay together (e.g., `Images/ProjectPhotos`)
- **MIME-based Accuracy**: Handles files with wrong or missing extensions
- **Context-Aware**: Files in directories with a dominant type (>50%) are kept together

## Code Organization

- **app.go**: Main application struct and directory analysis
- **types.go**: All data structures (FileInfo, FileMove, etc.)
- **hash.go**: SHA-256 content hash calculation
- **mime.go**: MIME type detection and categorization
- **categorize.go**: File categorization logic
- **context.go**: Directory context analysis
- **moves.go**: File move execution and conflict resolution

For detailed technical documentation and AI development guidelines, see [CLAUDE.md](CLAUDE.md).

## License

MIT