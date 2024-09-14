package main

import (
	"encoding/json"
	"fmt"
)

type PageVersion struct {
	Hash     string
	FilePath string
	Version  int
}

func ConstructFilePath(url string, version int) string {
	path := fmt.Sprintf("%s/v%d.html", NormalizeUrl(url), version)
	return path
}

func PageVersionsToJson(pageVersions []PageVersion) ([]byte, error) {
	return json.Marshal(pageVersions)
}

func PageVersionsFromJson(data []byte) ([]PageVersion, error) {
	var pageVersions []PageVersion
	err := json.Unmarshal(data, &pageVersions)

	return pageVersions, err
}
