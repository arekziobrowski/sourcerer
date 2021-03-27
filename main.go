package main

import (
	"flag"
	"fmt"
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

		git@github.com:go-git/go-git.git ef33fff761a2fabb7f0daf0c1779d2dfac1056da

	*/

	sources := ReadList("./example_list.txt")
	fmt.Println(sources)
	downloader := New("/Users/arek/test", false, false)

	err := downloader.GetSources(sources)
	fmt.Println(err)
	/*git := source.NewSystemGitDownloader("/Users/arek/test")
	err := git.Get("git@github.com:go-git/go-billy.git d62fe84467b62d39c6082ecb621c15771e2de9d9")
	fmt.Println(err)*/
}
