package main

import (
	"fmt"
	"os"
)

// implement the IDatabase interface
type FileDatabase struct {
	//add filename and file as fields
	filename string
	file     *os.File
}

// Open method which accepts the filename and opens the file
func (d *FileDatabase) Open(filename string) error {
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	d.file = file
	d.filename = filename
	return nil
}

func (d *FileDatabase) Close() {
	d.file.Close()
}

// Save method to save the data in the database
func (d *FileDatabase) Save(url string, md5Hash string) error {
	jsonString := fmt.Sprintf("{\"url\": \"%s\", \"md5\": \"%s\"}\n", url, md5Hash)
	d.file.WriteString(jsonString)

	return nil
}
