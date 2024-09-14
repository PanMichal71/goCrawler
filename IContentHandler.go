package main

type IContentHandler interface {
	HandleContent(url string, html string) error
}
