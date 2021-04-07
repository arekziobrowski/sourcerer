package source

import (
	"bytes"
	"os/exec"
	"strings"

	"github.com/arekziobrowski/sourcerer/model"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
)

type SystemGitDownloader struct {
	workingDirectory string
}

func NewSystemGitDownloader(wd string) *SystemGitDownloader {
	return &SystemGitDownloader{
		workingDirectory: wd,
	}
}

func (g *SystemGitDownloader) Get(src *model.Source) error {
	const remoteName = "origin"
	log.Infof("Downloading %s-%s", src.Origin, src.Hash)

	err := g.initialize()
	if err != nil {
		return errors.Wrapf(err, "failed to initialize the repository: %s", src)
	}

	err = g.remoteAdd(remoteName, src.Origin)
	if err != nil {
		return errors.Wrapf(err, "failed to add remote for %s", src.Origin)
	}

	err = g.fetch(remoteName, src.Hash)
	if err != nil {
		return errors.Wrapf(err, "failed to fetch from remote for revision: %s", src.Hash)
	}

	err = g.reset()
	if err != nil {
		return errors.Wrap(err, "failed to reset to FETCH_HEAD")
	}

	log.Infof("Finished downloading for: %s-%s", src.Origin, src.Hash)

	return nil
}

func (g *SystemGitDownloader) initialize() error {
	return run(g.workingDirectory, "git", "init")
}

func (g *SystemGitDownloader) remoteAdd(originName, remote string) error {
	return run(g.workingDirectory, "git", "remote", "add", originName, remote)
}

func (g *SystemGitDownloader) fetch(originName, hash string) error {
	return run(g.workingDirectory, "git", "fetch", originName, hash, "--depth=1")
}

func (g *SystemGitDownloader) reset() error {
	return run(g.workingDirectory, "git", "reset", "FETCH_HEAD", "--hard")
}

func run(wd, command string, args ...string) error {
	cmd := exec.Command(command, args...)
	cmd.Dir = wd
	var stderr bytes.Buffer
	cmd.Stderr = &stderr
	err := cmd.Run()
	if err != nil {
		log.Printf("Error occured when running %q: %s", command+" "+strings.Join(args, " "), stderr.String())
		return err
	}
	return nil
}
