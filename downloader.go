package main

import (
	"context"
	"log"

	"github.com/arekziobrowski/sourcerer/model"
	"github.com/arekziobrowski/sourcerer/source"
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
	destinationDir           string
	withDependencies         bool
	strict                   bool
}

func New(srcs []*model.Source, dst string, withDependencies bool, strict bool) *service {
	return &service{
		sources:                  srcs,
		sourceDownloaderType:     GitDirect,
		dependencyDownloaderType: MavenSystem,
		destinationDir:           dst,
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
			// TODO: prepare directory tree
			downloader := s.selectSourceDownloader(s.destinationDir)
			err := downloader.Get(src)
			if err != nil {
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

func (s *service) selectSourceDownloader(wd string) Downloader {
	switch s.sourceDownloaderType {
	case GitDirect:
		return source.NewGitDownloader(wd)
	case GitSystem:
		return source.NewSystemGitDownloader(wd)
	default:
		return source.NewSystemGitDownloader(wd)
	}
}
