package source

import (
	"context"

	"golang.org/x/sync/errgroup"
)

type Downloader interface {
	Get(src string) error
}

type service struct {
	sourceDownloader     Downloader
	dependencyDownloader Downloader
	destinationDir       string
	withDependencies     bool
	strict               bool
}

func New(dst string, withDependencies bool, strict bool) *service {
	return &service{
		sourceDownloader:     NewGitDownloader(dst),
		dependencyDownloader: NewMavenDownloader(dst),
		destinationDir:       dst,
		withDependencies:     withDependencies,
		strict:               strict,
	}
}

func (s *service) GetSources(list []string) error {
	eg, _ := errgroup.WithContext(context.Background())
	for _, src := range list {
		eg.Go(func() error {
			err := s.sourceDownloader.Get(src)
			return err
		})
	}
	if err := eg.Wait(); err != nil {
		return err
	}
	return nil
}

func downloadDependency() {

}

func setupDirectory() {

}
