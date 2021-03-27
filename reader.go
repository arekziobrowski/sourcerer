package main

import (
	"io/ioutil"
	"strings"
)

func ReadList(filename string) []string {
	fileBB, err := ioutil.ReadFile(filename)
	if err != nil {
		panic(err)
	}
	contents := string(fileBB)
	return strings.Split(contents, "\n")
}
