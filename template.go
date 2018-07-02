package main

import (
	"io/ioutil"
	"fmt"
)

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func initTemplate(dir string) error {

}

func main() {
	data, err := ioutil.ReadFile("index.html")
	check(err)
	fmt.Println(string(data))
}
