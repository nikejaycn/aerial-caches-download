package main

import (
	"./pack"
)

func main() {
	// entriesresources.Download()
	entriesresources.HandleRetries("./temp/entries.json")
	// fmt.Printf("Hello, worldxxx\n")
}
