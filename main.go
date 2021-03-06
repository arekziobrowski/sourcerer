package main

import (
	"flag"
	"os"
	"path/filepath"

	log "github.com/sirupsen/logrus"
)

const (
	git       = "git"
	gitSystem = "git-system"
)

var input = flag.String("input", "", "input file name")
var withDependencies = flag.Bool("with_dependencies", false, "download dependencies from Maven along with the sources")
var destination = flag.String("dst", "", "directory to which the sources will be downloaded")
var strict = flag.Bool("strict", false, "use strict mode")
var sourceDownloader = flag.String("source_downloader", gitSystem, "source downloader mode to use [git, git-system]")

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

	sourceDownloaderType := getSourceDownloaderType(*sourceDownloader)

	downloader := New(sources, *destination, sourceDownloaderType, *withDependencies, *strict)

	err := downloader.GetSources()
	if err != nil {
		log.Errorf("Error while downloading sources: %v", err)
		os.Exit(1)
	}
}

func getSourceDownloaderType(s string) SourceDownloaderType {
	switch s {
	case git:
		return GitDirect
	case gitSystem:
		return GitSystem
	default:
		return GitSystem
	}
}
