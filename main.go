package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"
)

type DownloadOrganizer struct {
	FileExtensions map[string][]string
	DownloadsPath  string
}

var (
	downloads_path = filepath.Join(os.Getenv("HOME"), "Downloads")
	prependDate    = flag.Bool("prependDate", false, "Prepend the date to the filename")
	fileExtensions = map[string][]string{
		"audio":      {".mp3", ".wav", ".flac", ".m4a"},
		"video":      {".mp4", ".mkv", ".flv", ".avi"},
		"image":      {".jpg", ".jpeg", ".png", ".gif", ".bmp", ".svg", ".webp", ".ico", ".tiff", ".psd", ".svgz"},
		"document":   {".pdf", ".doc", ".docx", ".xls", ".xlsx", ".ppt", ".pptx", ".txt", ".rtf"},
		"archive":    {".dmg", ".zip", ".rar", ".tar", ".gz", ".7z"},
		"executable": {".exe", ".msi"},
		"programming": {
			".c", ".h", ".cpp", ".hpp", ".cc", ".cxx", ".h", ".hxx", ".cs", ".java", ".py", ".pyc",
			".pyd", ".pyo", ".pyw", ".pyz", ".js", ".html", ".htm", ".css", ".php", ".rb", ".swift",
			".go", ".rs", ".kt", ".kts", ".ts", ".tsx", ".pl", ".pm", ".sh", ".bash", ".ps1",
			".psm1", ".m", ".r", ".rdata", ".rds", ".lua", ".scala", ".groovy", ".dart", ".hs",
			".lhs", ".lisp", ".lsp", ".pl", ".f", ".f90", ".f95", ".asm", ".m", ".swift", ".dart",
			".jl", ".erl", ".hrl", ".cr", ".ex", ".exs", ".ml", ".mli", ".fs", ".fsx", ".fsi",
			".fsproj", ".vhd", ".vhdl", ".v", ".vlog",
		},
	}
)

// NewDownloadOrganizer creates a new instance of DownloadOrganizer.
//
// It takes a string parameter `downloadsPath` which represents the path to the downloads directory.
// It returns a pointer to DownloadOrganizer.
func NewDownloadOrganizer(downloadsPath string) *DownloadOrganizer {

	return &DownloadOrganizer{
		FileExtensions: fileExtensions,
		DownloadsPath:  downloadsPath,
	}
}

func main() {
	flag.Parse()
	organizer := NewDownloadOrganizer(downloads_path)
	organizer.OrganizeDownloads(*prependDate)
}

// OrganizeDownloads organizes the downloaded files by moving them to the appropriate directory.
//
// prependDate specifies whether or not to prepend the date to the file name.
// This function does not return anything.
func (o *DownloadOrganizer) OrganizeDownloads(prependDate bool) {
	files, err := os.ReadDir(o.DownloadsPath)
	if err != nil {
		fmt.Println("Error reading directory:", err)
		return
	}

	for _, file := range files {
		if !file.IsDir() {
			err := o.MoveFile(filepath.Join(o.DownloadsPath, file.Name()), file.Name(), prependDate)
			if err != nil {
				fmt.Println("Error moving file:", err)
				return
			}
		}
	}
}

// MoveFile moves a file to a specified category based on its file extension.
//
// Parameters:
// - filePath: the path of the file to be moved.
// - fileName: the name of the file to be moved.
// - prependDate: a flag indicating whether to prepend the current date to the file name.
//
// Returns:
// - error: an error if any occurred during the file moving process, otherwise nil.
func (o *DownloadOrganizer) MoveFile(filePath, fileName string, prependDate bool) error {
	fileExtension := strings.ToLower(filepath.Ext(fileName))

	for category, extensions := range o.FileExtensions {
		for _, extension := range extensions {
			if fileExtension == extension {
				categoryPath := filepath.Join(o.DownloadsPath, category)
				err := os.MkdirAll(categoryPath, os.ModePerm)
				if err != nil {
					return err
				}

				newFileName := fileName
				if prependDate {
					// Generate the current time and format it to string
					currentTime := time.Now().Format("2006-01-02")
					// Add the current time to the beginning of the file name
					newFileName = fmt.Sprintf("%s_%s", currentTime, fileName)
				}

				err = os.Rename(filePath, filepath.Join(categoryPath, newFileName))
				if err != nil {
					return err
				}
				return nil
			}
		}
	}
	return nil
}
