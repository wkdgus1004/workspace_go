package main

import (
	. "fmt"
	"regexp"
)

func main() {
	myString := ".dflwef;pq!"
	re, _ := regexp.Compile(`[a-zA-Z]*[a-zA-Z]`)
	temp := re.FindAllString(string(myString), -1)
	Println(temp)
}
