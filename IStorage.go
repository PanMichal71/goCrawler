package main

type IStorage interface {
	Write(data []byte) error
	Open(filename string) error
	Close()
}
