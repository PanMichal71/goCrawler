package main

// LinkFilter is an interface that requires any implementing type to have a FilterLink method.
type LinkFilter interface {
    FilterLink(link string) bool
}
