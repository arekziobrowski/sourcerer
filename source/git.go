package source

import (
	"fmt"
	"os/user"
	"strings"

	"github.com/go-git/go-billy/v5/memfs"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/config"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/transport"
	"github.com/go-git/go-git/v5/plumbing/transport/ssh"
	"github.com/go-git/go-git/v5/storage/memory"
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
	origin, hash, err := extractOriginAndHash(src)
	if err != nil {
		return err
	}
	storage := memory.NewStorage()
	fs := memfs.New()
	repo, err := git.Init(storage, fs)
	if err != nil {
		return errors.Wrap(err, "failed to init repo")
	}

	remote, err := repo.CreateRemote(&config.RemoteConfig{
		Name: remoteName,
		URLs: []string{origin},
	})
	if err != nil {
		return errors.Wrapf(err, "failed to invoke 'git remote add %s %s'", remoteName, origin)
	}

	sshAuth := getSshKeyAuth()
	err = remote.Fetch(&git.FetchOptions{
		RemoteName: remoteName + " " + hash, // TODO: add hash
		Auth:       sshAuth,
		Depth:      1,
	})
	if err != nil {
		return errors.Wrapf(err, "failed to invoke 'git fetch %s %s --depth=1'", remoteName, hash)
	}

	workTree, err := repo.Worktree()
	if err != nil {
		return err
	}

	/*fetchHead, err := repo.Reference(plumbing.ReferenceName("FETCH_HEAD"), true)
	if err != nil {
		return errors.Wrap(err, "failed to resolve FETCH_HEAD reference")
	}*/

	refs, _ := repo.References()
	refs.ForEach(func(ref *plumbing.Reference) error {
		if ref.Type() == plumbing.HashReference {
			fmt.Println(ref)
		}

		return nil
	})

	h, err := repo.ResolveRevision(plumbing.Revision(hash))
	fmt.Println(h)

	/*	head, err := repo.Head()
		if err != nil {
			return errors.Wrap(err, "failed to resolve HEAD reference")
		}*/

	err = workTree.Reset(&git.ResetOptions{
		Commit: plumbing.NewHash(hash),
		Mode:   git.HardReset,
	})

	if err != nil {
		return errors.Wrap(err, "failed to invoke 'git reset --hard FETCH_HEAD'")
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

func getSshKeyAuth() transport.AuthMethod {
	usr, _ := user.Current()
	auth, _ := ssh.NewPublicKeysFromFile("git", usr.HomeDir+"/.ssh/id_rsa", "")
	return auth
}
