package main

import (
	"crypto/md5"
	"encoding/hex"
	"log"
)

type DifferenceTracker struct {
	database    IDatabase
	fileStorage IStorage
}

func NewDifferenceTracker(database IDatabase, fileStorage IStorage) *DifferenceTracker {
	return &DifferenceTracker{
		database:    database,
		fileStorage: fileStorage,
	}
}

func (diffTracker *DifferenceTracker) HandleContent(url string, htmlContent string) error {
	md5Hash := getMD5Hash(htmlContent)

	urlExists, err := diffTracker.database.Exists(url)
	if err != nil {
		handleError(err, "Error checking if URL exists in database, url="+url)
		return err
	}

	if urlExists {
		return diffTracker.updateExistingContent(url, md5Hash, htmlContent)
	} else {
		return diffTracker.storeNewContent(url, md5Hash, htmlContent)
	}
}

func (diffTracker *DifferenceTracker) updateExistingContent(url, md5Hash, htmlContent string) error {
	versionsBytes, err := diffTracker.database.Read(url)
	if err != nil {
		handleError(err, "Error reading versions from database for url="+url)
		return err
	}

	pageVersions, err := PageVersionsFromJson(versionsBytes)
	if err != nil {
		handleError(err, "Error parsing page versions from JSON, bytes="+string(versionsBytes))
		return err
	}

	latestPageVersion := pageVersions[len(pageVersions)-1]

	if latestPageVersion.Hash != md5Hash {
		newPageVersion := diffTracker.createPageVersion(url, latestPageVersion.Version+1, md5Hash)
		if err := diffTracker.writeHtmlToFileStorage(newPageVersion, htmlContent); err != nil {
			return err
		}

		pageVersions = append(pageVersions, newPageVersion)
		return diffTracker.storePageVersionsInDatabase(url, pageVersions)
	}

	return nil
}

func (diffTracker *DifferenceTracker) storeNewContent(url, md5Hash, htmlContent string) error {
	newPageVersion := diffTracker.createPageVersion(url, 1, md5Hash)
	if err := diffTracker.writeHtmlToFileStorage(newPageVersion, htmlContent); err != nil {
		return err
	}

	pageVersions := []PageVersion{newPageVersion}
	return diffTracker.storePageVersionsInDatabase(url, pageVersions)
}

func (diffTracker *DifferenceTracker) createPageVersion(url string, version int, md5Hash string) PageVersion {
	return PageVersion{
		Hash:     md5Hash,
		FilePath: ConstructFilePath(url, version),
		Version:  version,
	}
}

func (diffTracker *DifferenceTracker) writeHtmlToFileStorage(pageVersion PageVersion, htmlContent string) error {
	if err := diffTracker.fileStorage.Open(pageVersion.FilePath); err != nil {
		handleError(err, "Error opening file for writing, path="+pageVersion.FilePath)
		return err
	}
	defer diffTracker.fileStorage.Close()

	if err := diffTracker.fileStorage.Write([]byte(htmlContent)); err != nil {
		handleError(err, "Error writing HTML content to file, file="+pageVersion.FilePath)
		return err
	}

	return nil
}

func (diffTracker *DifferenceTracker) storePageVersionsInDatabase(url string, pageVersions []PageVersion) error {
	bytes, err := PageVersionsToJson(pageVersions)
	if err != nil {
		handleError(err, "Error serializing page versions to JSON")
		return err
	}

	if err := diffTracker.database.Store(url, bytes); err != nil {
		handleError(err, "Error storing page versions in database")
		return err
	}

	// Notify the system of a new page detected (if required)
	return nil
}

func getMD5Hash(text string) string {
	hash := md5.Sum([]byte(text))
	return hex.EncodeToString(hash[:])
}

func handleError(err error, message string) {
	log.Printf("ERROR: %s: %v", message, err)
}
