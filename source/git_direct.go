package source

import (
	"fmt"
	"os/user"
	"path/filepath"

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

	perr = prepare(destinationDir)
	if perr != nil {
		return errors.Wrapf(perr, "unable to prepare the directory tree for %q", destinationDir)
	}

	fs := osfs.New(destinationDir)
	dot, perr := fs.Chroot(".git")
	if perr != nil {
		return errors.Wrap(perr, "cannot create a .git directory")
	}
	storage := filesystem.NewStorage(dot, cache.NewObjectLRUDefault())

	/*_, _ = git.Clone(storage, fs, &git.CloneOptions{
		URL:  "git@github.com:go-git/go-billy.git",
		Auth: getSshKeyAuth(),
	})
	fmt.Println()*/
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

	fmt.Println(remote)

	sshAuth := getSshKeyAuth()

	// get default branch https://stackoverflow.com/questions/28666357/git-how-to-get-default-branch then use as master here
	// check git config --get uploadpack.allowReachableSHA1InWant https://github.com/src-d/go-git/issues/628
	refSpec := config.RefSpec(fmt.Sprintf("%v:%v", hash, "refs/remotes/origin/master"))
	err = remote.Fetch(&git.FetchOptions{
		RemoteName: remoteName,
		Depth:      1,
		RefSpecs: []config.RefSpec{
			refSpec,
		},
		Auth: sshAuth,
	})
	if err != nil {
		return errors.Wrapf(err, "failed to invoke 'git fetch %s %s --depth=1'", remoteName, hash)
	}

	workTree, err := repo.Worktree()
	if err != nil {
		return err
	}

	refs, _ := repo.References()
	refs.ForEach(func(ref *plumbing.Reference) error {
		if ref.Type() == plumbing.HashReference {
			fmt.Println(ref)
		}

		return nil
	})

	h, err := repo.ResolveRevision(plumbing.Revision(hash))
	if err != nil {
		return errors.Wrap(err, "cannot resolve revision")
	}

	err = createHeadRef(fs, "master")
	if err != nil {
		return err
	}

	err = workTree.Reset(&git.ResetOptions{
		Commit: *h,
		Mode:   git.HardReset,
	})

	if err != nil {
		return errors.Wrap(err, "failed to invoke 'git reset --hard FETCH_HEAD'")
	}
	return nil
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
