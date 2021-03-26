package source

import (
	"bytes"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/pkg/errors"
)

type GitDownloader struct {
	workingDirectory string
}

func NewGitDownloader(wd string) *GitDownloader {
	return &GitDownloader{
		workingDirectory: wd,
	}
}

func (g *GitDownloader) Get(src string) error {
	const remoteName = "origin"
	origin, hash, perr := extractOriginAndHash(src)
	if perr != nil {
		return perr
	}

	org, repoName := extractOrganizationAndRepo(origin)
	destinationDir := filepath.Join(g.workingDirectory, org, repoName)

	perr = g.prepare(destinationDir)
	if perr != nil {
		return errors.Wrapf(perr, "unable to prepare the directory tree for %q", destinationDir)
	}

	var err error
	defer g.cleanUpIfError(err, destinationDir)

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

func (g *GitDownloader) prepare(path string) error {
	return os.MkdirAll(path, 0777)
}

func (g *GitDownloader) cleanUpIfError(err error, dir string) error {
	if err != nil {
		fmt.Println("removing", dir)
		return nil
		//return os.RemoveAll(dir)
	}
	return nil
}

func (g *GitDownloader) initialize(wd string) error {
	return run(wd, "git", "init")
}

func (g *GitDownloader) remoteAdd(originName, remote, wd string) error {
	return run(wd, "git", "remote", "add", originName, remote)
}

func (g *GitDownloader) fetch(originName, hash, wd string) error {
	return run(wd, "git", "fetch", originName, hash, "--depth=1")
}

func (g *GitDownloader) reset(wd string) error {
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

func extractOriginAndHash(src string) (string, string, error) {
	split := strings.Split(src, " ")
	if len(split) != 2 {
		return "", "", errors.Errorf("invalid origin and hash definition: %s", src)
	}
	return split[0], split[1], nil
}

func extractOrganizationAndRepo(origin string) (string, string) {
	organizationAndRepo := strings.TrimSuffix(strings.Split(origin, ":")[1], ".git")
	split := strings.Split(organizationAndRepo, "/")
	return split[0], split[1]
}
