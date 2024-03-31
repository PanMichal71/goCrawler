package main

// prepare IDatabase inteface which will have the method to save the data and accept two arguments: url and md5 hash
type IDatabase interface {
	Save(url string, hash string) error
	Open(filename string) error
	Close()
}
