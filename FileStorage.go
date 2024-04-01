package main

import (
	"fmt"
	"os"
)

// implement the IStorage interface
type FileStorage struct {
	//add filename and file as fields
	filename string
	file     *os.File
}

// Open method which accepts the filename and opens the file
func (d *FileStorage) Open(filename string) error {
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	d.file = file
	d.filename = filename
	return nil
}

func (d *FileStorage) Close() {
	d.file.Close()
}

// Save method to save the data in the Storage
func (d *FileStorage) Save(url string, md5Hash string) error {
	jsonString := fmt.Sprintf("{\"url\": \"%s\", \"md5\": \"%s\"}\n", url, md5Hash)
	d.file.WriteString(jsonString)

	return nil
}
