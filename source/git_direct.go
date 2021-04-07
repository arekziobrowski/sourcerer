package source

import (
	"fmt"
	"os/user"
	"strings"

	"github.com/arekziobrowski/sourcerer/model"
	"github.com/go-git/go-billy/v5"
	"github.com/go-git/go-billy/v5/osfs"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/config"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/cache"
	"github.com/go-git/go-git/v5/plumbing/transport"
	"github.com/go-git/go-git/v5/plumbing/transport/ssh"
	"github.com/go-git/go-git/v5/storage/filesystem"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
)

type GitDownloader struct {
	workingDirectory string
}

func NewGitDownloader(wd string) *GitDownloader {
	return &GitDownloader{
		workingDirectory: wd,
	}
}

func (g *GitDownloader) Get(src *model.Source) error {
	const remoteName = "origin"
	log.Infof("Downloading %s-%s", src.Origin, src.Hash)

	fs := osfs.New(g.workingDirectory)
	dot, err := fs.Chroot(".git")
	if err != nil {
		return errors.Wrap(err, "cannot create a .git directory")
	}
	storage := filesystem.NewStorage(dot, cache.NewObjectLRUDefault())

	repo, err := git.Init(storage, fs)
	if err != nil {
		return errors.Wrap(err, "failed to init repo")
	}

	remote, err := repo.CreateRemote(&config.RemoteConfig{
		Name: remoteName,
		URLs: []string{src.Origin},
	})
	if err != nil {
		return errors.Wrapf(err, "failed to invoke 'git remote add %s %s'", remoteName, src.Origin)
	}

	sshAuth := getSshKeyAuth()

	branch, err := getDefaultBranchName(remote)
	if err != nil {
		return err
	}

	refSpec := config.RefSpec(fmt.Sprintf("%v:%v", src.Hash, strings.Join([]string{"refs/remotes", remoteName, branch}, "/")))
	err = remote.Fetch(&git.FetchOptions{
		RemoteName: remoteName,
		Depth:      1,
		RefSpecs: []config.RefSpec{
			refSpec,
		},
		Auth: sshAuth,
	})
	if err != nil {
		return errors.Wrapf(err, "failed to invoke 'git fetch %s %s --depth=1'", remoteName, src.Hash)
	}

	workTree, err := repo.Worktree()
	if err != nil {
		return err
	}

	h, err := repo.ResolveRevision(plumbing.Revision(src.Hash))
	if err != nil {
		return errors.Wrap(err, "cannot resolve revision")
	}

	err = createHeadRef(fs, branch)
	if err != nil {
		return err
	}

	err = workTree.Reset(&git.ResetOptions{
		Commit: *h,
		Mode:   git.HardReset,
	})

	if err != nil {
		return errors.Wrap(err, "failed to invoke 'git reset --hard SHA1'")
	}
	log.Infof("Finished downloading for: %s-%s", src.Origin, src.Hash)
	return nil
}

func getDefaultBranchName(remote *git.Remote) (string, error) {
	refs, err := remote.List(&git.ListOptions{Auth: getSshKeyAuth()})
	if err != nil {
		return "", errors.Wrap(err, "cannot invoke ls-remote")
	}
	for _, ref := range refs {
		if ref.Name() == "HEAD" {
			target := strings.Split(ref.Target().String(), "/")
			return target[len(target)-1], nil
		}
	}
	return "", errors.New("cannot find a HEAD reference in remote references list")
}

func createHeadRef(fs billy.Filesystem, refName string) error {
	heads, err := fs.Chroot(".git/refs/heads")
	if err != nil {
		return errors.Wrap(err, "cannot chroot for .git/refs/heads")
	}
	_, err = heads.Create(refName)
	if err != nil {
		return errors.Wrapf(err, "cannot create a refname: %s", refName)
	}
	return nil
}

func getSshKeyAuth() transport.AuthMethod {
	usr, _ := user.Current()
	auth, _ := ssh.NewPublicKeysFromFile("git", usr.HomeDir+"/.ssh/id_rsa", "")
	return auth
}
