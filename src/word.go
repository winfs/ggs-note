package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"strings"
)

func main() {
	file, err := ioutil.ReadFile("word.txt")
	if err != nil {
		fmt.Println(err)
		return
	}
	words := strings.FieldsFunc(string(file), func(c rune) bool {
		return c == ',' || c == '、' || c == '，'
	})

	ioutil.WriteFile("word-new.txt", []byte(strings.Join(words, "\n")), os.ModePerm)
	fmt.Println("success")
}
