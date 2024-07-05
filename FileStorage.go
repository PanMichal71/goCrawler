package main

import (
	"fmt"
	"os"
)

// implement the IStorage interface
type FileStorage struct {
	//add filename and file as fields
	directory string
	filename  string
	file      *os.File
	//number of records written
	recordsWritten int
}

// NewFileStorage method to create a new FileStorage object, must accept a directory path as an argument
func NewFileStorage(directory string) *FileStorage {
	if _, err := os.Stat(directory); os.IsNotExist(err) {
		os.Mkdir(directory, os.ModePerm)
	}

	if _, err := os.Stat(directory); os.IsNotExist(err) {
		panic("Failed to create directory")
	}

	return &FileStorage{directory: directory}
}

// Open method which accepts the filename and opens the file
func (d *FileStorage) Open(filename string) error {

	filename = d.directory + "/" + filename
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	d.file = file
	d.filename = filename
	d.file.WriteString("[")
	d.recordsWritten = 0
	return nil
}

func (d *FileStorage) Close() {
	d.file.WriteString("]")
	d.file.Close()
}

// Save method to save the data in the Storage
func (d *FileStorage) Save(url string, md5Hash string) error {
	if d.recordsWritten > 0 {
		d.file.WriteString(",\n")
	}
	d.recordsWritten++
	jsonString := fmt.Sprintf("{\"url\": \"%s\", \"md5\": \"%s\"}", url, md5Hash)
	d.file.WriteString(jsonString)

	return nil
}
