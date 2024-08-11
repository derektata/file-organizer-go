package main

import (
	"mime"
	"path/filepath"
	"strings"
)

// categoryFromMimeType categorizes a file based on its MIME type.
func (o *FileOrganizer) categoryFromMimeType(filePath string) string {
	mimeType := guessMimeType(filePath)
	if mimeType == "" {
		return ""
	}
	mainType := strings.Split(mimeType, "/")[0]
	switch mainType {
	case "image":
		return "image"
	case "audio":
		return "audio"
	case "video":
		return "video"
	case "application":
		return "document" // General fallback for application types
	default:
		return "others"
	}
}

// guessMimeType guesses the MIME type of a file based on its extension.
func guessMimeType(filePath string) string {
	ext := strings.ToLower(filepath.Ext(filePath))
	return mime.TypeByExtension(ext)
}
