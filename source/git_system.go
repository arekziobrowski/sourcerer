package source

import (
	"bytes"
	"log"
	"os/exec"
	"path/filepath"
	"strings"

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

func (g *SystemGitDownloader) Get(src string) error {
	const remoteName = "origin"
	origin, hash, perr := extractOriginAndHash(src)
	if perr != nil {
		return perr
	}

	org, repoName := extractOrganizationAndRepo(origin)
	destinationDir := filepath.Join(g.workingDirectory, org, repoName)

	perr = prepare(destinationDir)
	if perr != nil {
		return errors.Wrapf(perr, "unable to prepare the directory tree for %q", destinationDir)
	}

	var err error
	defer cleanUpIfError(err, destinationDir)

	err = g.initialize(destinationDir)
	if err != nil {
		return errors.Wrapf(err, "failed to initialize the repository: %s", src)
	}

	err = g.remoteAdd(remoteName, origin, destinationDir)
	if err != nil {
		return errors.Wrapf(err, "failed to add remote for %s", origin)
	}

	err = g.fetch(remoteName, hash, destinationDir)
	if err != nil {
		return errors.Wrapf(err, "failed to fetch from remote for revision: %s", hash)
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
