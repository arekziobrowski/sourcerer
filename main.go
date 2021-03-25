package main

import (
	"flag"
	"fmt"

	"github.com/arekziobrowski/sourcerer/source"
)

var input = flag.String("input", "", "input file name")
var withDependencies = flag.Bool("with_dependencies", false, "download dependencies from Maven along with the sources")
var destination = flag.String("dst", "downloaded-sources", "directory to which the sources will be downloaded")
var strict = flag.Bool("strict", false, "use strict mode")

func main() {
	//flag.Parse()
	/*
		1. git init
		2. git remote add origin git@github.com:go-git/go-git.git
		3. git fetch origin ef33fff761a2fabb7f0daf0c1779d2dfac1056da --depth=1
		4. git reset --hard FETCH_HEAD
	*/

	git := source.NewGitDownloader("/Users/arek/test")
	err := git.Get("git@github.com:go-git/go-git.git ef33fff761a2fabb7f0daf0c1779d2dfac1056da")
	fmt.Println(err)
}
