package main

import (
	"flag"
	"os"
	"path/filepath"

	log "github.com/sirupsen/logrus"
)

var input = flag.String("input", "", "input file name")
var withDependencies = flag.Bool("with_dependencies", false, "download dependencies from Maven along with the sources")
var destination = flag.String("dst", "", "directory to which the sources will be downloaded")
var strict = flag.Bool("strict", false, "use strict mode")

func main() {
	flag.Parse()

	if *input == "" {
		log.Errorf("Input is missing, please use --input flag to provide the input")
		flag.Usage()
		os.Exit(1)
	}

	log.Infof("Reading sources from: %s", *input)
	sources := ReadList(*input)

	if *destination == "" {
		usr, _ := os.UserHomeDir()
		*destination = filepath.Join(usr, "downloaded-sources")
	}

	downloader := New(sources, *destination, *withDependencies, *strict)

	err := downloader.GetSources()
	if err != nil {
		log.Error(err)
		os.Exit(1)
	}
}
