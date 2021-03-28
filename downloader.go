package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/arekziobrowski/sourcerer/model"
	"github.com/arekziobrowski/sourcerer/source"
	"github.com/pkg/errors"
	"golang.org/x/sync/errgroup"
)

type Downloader interface {
	Get(src *model.Source) error
}

type SourceDownloaderType int

const (
	GitDirect SourceDownloaderType = 1 << iota
	GitSystem
)

type DependencyDownloaderType int

const (
	MavenSystem DependencyDownloaderType = 1 << iota
)

type service struct {
	sources                  []*model.Source
	sourceDownloaderType     SourceDownloaderType
	dependencyDownloaderType DependencyDownloaderType
	rootDir                  string
	withDependencies         bool
	strict                   bool
}

func New(srcs []*model.Source, dir string, withDependencies bool, strict bool) *service {
	return &service{
		sources:                  srcs,
		sourceDownloaderType:     GitDirect,
		dependencyDownloaderType: MavenSystem,
		rootDir:                  dir,
		withDependencies:         withDependencies,
		strict:                   strict,
	}
}

func (s *service) GetSources() error {
	eg, _ := errgroup.WithContext(context.Background())
	for _, src := range s.sources {
		src := src
		eg.Go(func() error {
			log.Println("Downloading", src)
			wd := filepath.Join(s.rootDir, src.Organization, src.Repository)
			if err := prepareDirectoryTree(wd); err != nil {
				return err
			}
			downloader := s.createSourceDownloader(wd)
			if err := downloader.Get(src); err != nil {
				cleanupErr := cleanUpDirectoryTree(wd)
				if cleanupErr != nil {
					return errors.Wrapf(err, "error occurred when cleaning up directory: %v", cleanupErr)
				}
				return err
			}
			return nil
		})
	}
	if err := eg.Wait(); err != nil {
		return err
	}
	return nil
}

func (s *service) createSourceDownloader(wd string) Downloader {
	switch s.sourceDownloaderType {
	case GitDirect:
		return source.NewGitDownloader(wd)
	case GitSystem:
		return source.NewSystemGitDownloader(wd)
	default:
		return source.NewSystemGitDownloader(wd)
	}
}

func prepareDirectoryTree(path string) error {
	return os.MkdirAll(path, 0777)
}

func cleanUpDirectoryTree(dir string) error {
	fmt.Println("Removing", dir)
	//return os.RemoveAll(dir)
	return nil
}
