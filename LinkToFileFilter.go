package main

import (
	"net/url"
	"strings"
)

type LinkToFileFilter struct{}

// FilterLink checks if the link ends with a common file extension.
func (l LinkToFileFilter) FilterLink(link string) bool {
	parsedURL, err := url.Parse(strings.ToLower(link))
	if err != nil {
		return false // Unable to parse URL, cannot determine if it points to a file resource.
	}

	extension := getExtension(parsedURL.Path)
	_, present := fileExtensions[extension]
	return present
}

func getExtension(resourcePath string) string {
	dotIndex := strings.LastIndex(resourcePath, ".")
	if dotIndex < 0 || dotIndex == len(resourcePath)-1 {
		return "" // No extension or resourcePath ends with "."
	}
	return resourcePath[dotIndex+1:]
}

var fileExtensions = map[string]bool{
	"pdf": true, "doc": true, "docx": true, "xls": true, "xlsx": true,
	"ppt": true, "pptx": true, "jpg": true, "png": true,
	"gif": true, "zip": true, "rar": true, "7z": true, "mp3": true,
	"mp4": true, "avi": true, "mov": true, "wmv": true, "flv": true,
	"swf": true, "exe": true, "msi": true, "apk": true, "dmg": true,
	"iso": true, "torrent": true, "jar": true, "svg": true,
	"eps": true, "ai": true, "psd": true, "ttf": true, "otf": true,
	"woff": true, "woff2": true, "eot": true, "ico": true, "bmp": true,
	"tif": true, "tiff": true, "svgz": true, "webp": true, "json": true,
	"xml": true, "csv": true, "txt": true, "rtf": true,
	"odt": true, "ods": true, "odp": true, "odg": true, "odf": true,
	"epub": true, "mobi": true, "azw": true, "azw3": true, "fb2": true,
	"djvu": true, "djv": true, "chm": true, "pdb": true, "xps": true,
	"cbr": true, "cbz": true, "cb7": true, "cbt": true, "cba": true,
	"cbw": true, "lit": true, "prc": true, "pml": true,
	"rb": true, "c": true, "tar": true, "gz": true,
}
