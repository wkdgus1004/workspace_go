package main

import (
	. "fmt"
	"io/ioutil"
	"log"
)

func main() {
	filePath := "test.txt"
	bookData, err := ioutil.ReadFile(filePath)
	if err != nil {
		log.Fatal(err)
	}
	Println(string(bookData))
}
