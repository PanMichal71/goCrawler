package main

import (
	"fmt"
	"net/url"
	"strings"
)

type PathExclusionFilter struct {
	exclusionPaths map[string]struct{}
}

func NewPathExclusionFilter(exclusionPaths []string) PathExclusionFilter {
	exclusionMap := make(map[string]struct{})
	for _, path := range exclusionPaths {
		exclusionMap[strings.ToLower(path)] = struct{}{}
	}
	return PathExclusionFilter{exclusionPaths: exclusionMap}
}

func (l PathExclusionFilter) FilterLink(link string) bool {
	urlLink, err := url.Parse(link)
	if err != nil {
		fmt.Println("Error parsing URL:", err)
		return true
	}

	normalizedPath := strings.ToLower(urlLink.Path)
	pathSegments := strings.Split(normalizedPath, "/")
	pathSegments = pathSegments[1:] // Remove the first empty segment

	for windowSize := 1; windowSize <= len(pathSegments); windowSize++ {
		for i := 0; i <= len(pathSegments)-windowSize; i++ {
			subPath := strings.Join(pathSegments[i:i+windowSize], "/")
			if _, exists := l.exclusionPaths[subPath]; exists {
				return true
			}
		}
	}
	return false
}
