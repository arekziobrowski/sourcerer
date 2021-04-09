package main

import (
	"io/ioutil"
	"log"
	"strings"

	"github.com/arekziobrowski/sourcerer/model"
)

func ReadList(filename string) []*model.Source {
	fileBB, err := ioutil.ReadFile(filename)
	if err != nil {
		log.Fatalf("Error while reading file: %v", err)
	}
	contents := string(fileBB)
	lines := strings.Split(contents, "\n")
	lines = removeDuplicates(lines)
	out := make([]*model.Source, 0, len(lines))
	for _, line := range lines {
		src, err := model.ToSource(line)
		if err != nil {
			log.Fatalf("Error while converting input: %v", err)
		}
		out = append(out, src)
	}
	return out
}

func removeDuplicates(list []string) []string {
	set := make(map[string]struct{}, 0)
	for _, s := range list {
		set[s] = struct{}{}
	}

	var out []string
	for k := range set {
		out = append(out, k)
	}
	return out
}
