package main

import (
	"os"
	"strings"
)

type FileStorage struct {
	directory string
	filename  string
	file      *os.File
}

func NewFileStorage(directory string) *FileStorage {
	if _, err := os.Stat(directory); os.IsNotExist(err) {
		os.Mkdir(directory, os.ModePerm)
	}

	if _, err := os.Stat(directory); os.IsNotExist(err) {
		panic("Failed to create directory")
	}

	return &FileStorage{directory: directory}
}

func (d *FileStorage) Open(filename string) error {

	filename = d.directory + "/" + filename
	// get parent directory of filename
	parentDir := filename[:len(filename)-len(filename[strings.LastIndex(filename, "/"):])]

	if err := os.MkdirAll(parentDir, os.ModePerm); err != nil {
		return err
	}

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
	d.file = nil
}

func (d *FileStorage) Write(bytes []byte) error {
	if d.file == nil {
		panic("File not open")
	}

	_, err := d.file.Write(bytes)

	return err
}
