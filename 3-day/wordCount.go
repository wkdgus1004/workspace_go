package main

import (
	//"bufio"
	"fmt"
	"io/ioutil"
	"log"
	"regexp"
	"sort"
)

type wordStruct struct {
	word  string
	count int
}

func main() {
	filePath := "test.txt"
	threadNum := 2
	bookData, err := ioutil.ReadFile(filePath)
	if err != nil {
		log.Fatal(err)
	}

	words := filterRegex(string(bookData))
	subWords := sliceWords(threadNum, words)

	mapChannels := make(chan map[string]int, threadNum)
	flagMerge := make(chan bool)

	for i := range subWords {
		go func(i int) {
			var wordMap map[string]int
			defer func() {
				fmt.Println("Thread ", i, " is End!")
				mapChannels <- wordMap
			}()
			wordMap = countWords(subWords[i])
		}(i)
	}

	var myMapBuffer = make([]map[string]int, 0, threadNum)
	var myWordMap map[string]int

	go func() {
		defer func() {
			myWordMap = mergeMap(myMapBuffer)
			flagMerge <- false
		}()
		for i := 0; i < threadNum; i++ {
			select {
			case buff := <-mapChannels:
				myMapBuffer = append(myMapBuffer, buff)
			}
		}

	}()

	<-flagMerge
	wst := make([]wordStruct, 0)
	for key, val := range myWordMap {
		wst = append(wst, wordStruct{word: key, count: val})
	}

	sort.Slice(wst, func(i, j int) bool {
		return wst[i].word < wst[j].word
	})

	fmt.Println(wst)

}

func mergeMap(subCounter []map[string]int) map[string]int {
	// fmt.Println(subCounter)
	counter := map[string]int{}
	for _, subValue := range subCounter {
		for k, v := range subValue {
			counter[k] = counter[k] + v
		}
	}
	return counter
}

func sliceWords(threadNum int, words []string) [][]string {
	splitSize := len(words) / threadNum
	subWords := make([][]string, threadNum)
	for i := 0; i < threadNum; i++ {
		if i+1 == threadNum {
			subWords[i] = words
		}
		subWords[i] = words[:splitSize]
		words = words[splitSize:]
	}
	return subWords
}

func countWords(words []string) map[string]int {
	counter := make(map[string]int)
	for _, word := range words {
		counter[word]++
	}
	return counter
}

func filterRegex(bookData string) []string {
	re, _ := regexp.Compile(`[a-zA-Z]*[a-zA-Z]`)
	return re.FindAllString(string(bookData), -1)
}
