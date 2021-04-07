package main

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"sync"

	"github.com/arekziobrowski/sourcerer/dependency"
	"github.com/arekziobrowski/sourcerer/model"
	"github.com/arekziobrowski/sourcerer/source"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"golang.org/x/sync/errgroup"
)

type SourceDownloader interface {
	Get(src *model.Source) error
}

type DependencyDownloader interface {
	Get() error
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

func New(srcs []*model.Source, dir string, srcDownloaderType SourceDownloaderType, withDependencies bool, strict bool) *service {
	return &service{
		sources:                  srcs,
		sourceDownloaderType:     srcDownloaderType,
		dependencyDownloaderType: MavenSystem,
		rootDir:                  dir,
		withDependencies:         withDependencies,
		strict:                   strict,
	}
}

func (s *service) GetSources() error {
	var mutex sync.Mutex
	eg, _ := errgroup.WithContext(context.Background())
	for _, src := range s.sources {
		src := src
		eg.Go(func() error {
			wd := filepath.Join(s.rootDir, src.Organization, src.Repository+"-"+src.Hash)

			// We need to sync the preparation of directory tree, because the directory tree is nested
			// and two goroutines may try to create the same parent dir.
			mutex.Lock()
			if err := prepareDirectoryTree(wd); err != nil {
				return err
			}
			mutex.Unlock()

			downloader := s.createSourceDownloader(wd)
			if err := downloader.Get(src); err != nil {
				err = errors.Wrapf(err, "error while parsing: %s", fmt.Sprintf("%s@%s", src.Origin, src.Hash))
				if s.strict {
					return err
				}
				log.Errorf("Error occured: %v", err)
			}

			if s.withDependencies {
				dependencyDownloader := s.createDependencyDownloader(wd)
				if err := dependencyDownloader.Get(); err != nil {
					log.Errorf("Error occured: %v", err)
				}
			}
			return nil
		})
	}
	if err := eg.Wait(); err != nil {
		return err
	}
	return nil
}

func (s *service) createSourceDownloader(wd string) SourceDownloader {
	switch s.sourceDownloaderType {
	case GitDirect:
		return source.NewGitDownloader(wd)
	case GitSystem:
		return source.NewSystemGitDownloader(wd)
	default:
		return source.NewSystemGitDownloader(wd)
	}
}

func (s *service) createDependencyDownloader(wd string) DependencyDownloader {
	switch s.dependencyDownloaderType {
	case MavenSystem:
		return dependency.NewSystemMavenDownloader(wd)
	default:
		return dependency.NewSystemMavenDownloader(wd)
	}
}

func prepareDirectoryTree(path string) error {
	log.Infof("Creating directory: %s", path)
	return os.MkdirAll(path, 0777)
}
