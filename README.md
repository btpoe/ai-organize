# File Organizer

A desktop application built with Go and Wails that automatically analyzes and organizes files into logical folders.

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
├── app.go               # Backend logic (file analysis & organization)
├── go.mod               # Go dependencies
├── wails.json           # Wails configuration
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

- **AnalyzeDirectory**: Walks through the specified directory, categorizes files by extension, and proposes moves
- **ExecuteMoves**: Executes approved file moves, creates necessary folders, and handles conflicts
- **categorizeFile**: Determines the appropriate category for each file based on its extension

### Frontend (Svelte)

- Clean, modern UI with gradient design
- Real-time file selection and editing
- Grouped view by category
- Detailed move information (source, destination, reason)
- Execution results with success/failure reporting

## Safety Features

- **Dry run first**: Always shows proposed changes before executing
- **Conflict handling**: Automatically renames files if destination exists
- **Error reporting**: Detailed feedback on any failed operations
- **No destructive operations**: Only moves files, never deletes
- **Skip hidden files**: Ignores files and folders starting with "."

## License

MIT