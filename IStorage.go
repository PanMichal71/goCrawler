package main

type IStorage interface {
	Save(url string, hash string) error
	Open(filename string) error
	Close()
}
