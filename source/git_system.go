package source

import (
	"bytes"
	"log"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/arekziobrowski/sourcerer/model"
	"github.com/pkg/errors"
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

	destinationDir := filepath.Join(g.workingDirectory, src.Organization, src.Repository)

	err := g.initialize(destinationDir)
	if err != nil {
		return errors.Wrapf(err, "failed to initialize the repository: %s", src)
	}

	err = g.remoteAdd(remoteName, src.Origin, destinationDir)
	if err != nil {
		return errors.Wrapf(err, "failed to add remote for %s", src.Origin)
	}

	err = g.fetch(remoteName, src.Hash, destinationDir)
	if err != nil {
		return errors.Wrapf(err, "failed to fetch from remote for revision: %s", src.Hash)
	}

	err = g.reset(destinationDir)
	if err != nil {
		return errors.Wrap(err, "failed to reset to FETCH_HEAD")
	}

	return nil
}

func (g *SystemGitDownloader) initialize(wd string) error {
	return run(wd, "git", "init")
}

func (g *SystemGitDownloader) remoteAdd(originName, remote, wd string) error {
	return run(wd, "git", "remote", "add", originName, remote)
}

func (g *SystemGitDownloader) fetch(originName, hash, wd string) error {
	return run(wd, "git", "fetch", originName, hash, "--depth=1")
}

func (g *SystemGitDownloader) reset(wd string) error {
	return run(wd, "git", "reset", "FETCH_HEAD", "--hard")
}

func run(wd, command string, args ...string) error {
	cmd := exec.Command(command, args...)
	cmd.Dir = wd
	var stderr bytes.Buffer
	cmd.Stderr = &stderr
	err := cmd.Run()
	if err != nil {
		log.Printf("Error occured when running %q: %s", command+" "+strings.Join(args, " "), stderr)
		return err
	}
	return nil
}
