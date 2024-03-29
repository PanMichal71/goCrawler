package main

import (
	"net/url"
	"strings"
)

// LinkToFileFilter implements the LinkFilter interface to filter links that lead to files.
type LinkToFileFilter struct{}

// FilterLink checks if the link ends with a common file extension.
func (l LinkToFileFilter) FilterLink(link string) bool {
	parsedURL, err := url.Parse(strings.ToLower(link))
	if err != nil {
		return false // Unable to parse URL, cannot determine if it points to a file resource.
	}

	// Check if the path ends with any of the file extensions.
	for _, ext := range fileExtensions {
		if strings.HasSuffix(parsedURL.Path, ext) {
			return true
		}
	}

	return false
}

var fileExtensions = []string{
	".pdf", ".doc", ".docx", ".xls", ".xlsx", ".ppt", ".pptx", ".jpg", ".png",
	".gif", ".zip", ".rar", ".7z", ".mp3", ".mp4", ".avi", ".mov", ".wmv", ".flv",
	".swf", ".exe", ".msi", ".apk", ".dmg", ".iso", ".torrent", ".jar", ".svg",
	".eps", ".ai", ".psd", ".ttf", ".otf", ".woff", ".woff2", ".eot", ".ico", ".bmp",
	".tif", ".tiff", ".svgz", ".webp", ".json", ".xml", ".csv", ".txt", ".rtf",
	".odt", ".ods", ".odp", ".odg", ".odf", ".epub", ".mobi", ".azw", ".azw3", ".fb2",
	".djvu", ".djv", ".chm", ".pdb", ".xps", ".cbr", ".cbz", ".cb7", ".cbt", ".cba",
	".cbw", ".lit", ".prc", ".pdb", ".pml", ".rb", ".c", ".tar", ".gz",
}

func isKnownExtension(ext string) bool {
	for _, knownExt := range fileExtensions {
		if ext == knownExt {
			return true
		}
	}
	return false
}
