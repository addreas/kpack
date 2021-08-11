package git

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"

	git2go "github.com/libgit2/git2go/v31"
	"github.com/pkg/errors"

	"github.com/pivotal/kpack/pkg/apis/build/v1alpha2"
)

const defaultRemote = "origin"

var discardLogger = log.New(ioutil.Discard, "", 0)

type remoteGitResolver struct {
}

func (*remoteGitResolver) Resolve(keychain GitKeychain, sourceConfig v1alpha2.SourceConfig) (v1alpha2.ResolvedSourceConfig, error) {
	dir, err := ioutil.TempDir("", "git-resolve")
	if err != nil {
		return v1alpha2.ResolvedSourceConfig{}, err
	}
	defer os.RemoveAll(dir)

	repository, err := git2go.InitRepository(dir, false)
	if err != nil {
		return v1alpha2.ResolvedSourceConfig{}, errors.Wrap(err, "initializing repo")
	}
	defer repository.Free()

	remote, err := repository.Remotes.CreateWithOptions(sourceConfig.Git.URL, &git2go.RemoteCreateOptions{
		Name:  defaultRemote,
		Flags: git2go.RemoteCreateSkipInsteadof,
	})
	if err != nil {
		return v1alpha2.ResolvedSourceConfig{}, errors.Wrap(err, "create remote")
	}
	defer remote.Free()

	err = remote.ConnectFetch(&git2go.RemoteCallbacks{
		CredentialsCallback:      keychainAsCredentialsCallback(keychain),
		CertificateCheckCallback: certificateCheckCallback(discardLogger),
	}, nil, nil)
	if err != nil {
		return v1alpha2.ResolvedSourceConfig{
			Git: &v1alpha2.ResolvedGitSource{
				URL:      sourceConfig.Git.URL,
				Revision: sourceConfig.Git.Revision,
				Type:     v1alpha2.Unknown,
				SubPath:  sourceConfig.SubPath,
			},
		}, nil
	}

	references, err := remote.Ls()
	if err != nil {
		return v1alpha2.ResolvedSourceConfig{}, errors.Wrap(err, "remote ls")
	}

	for _, ref := range references {
		for _, format := range refRevParseRules {
			if fmt.Sprintf(format, sourceConfig.Git.Revision) == ref.Name {
				return v1alpha2.ResolvedSourceConfig{
					Git: &v1alpha2.ResolvedGitSource{
						URL:      sourceConfig.Git.URL,
						Revision: ref.Id.String(),
						Type:     sourceType(ref),
						SubPath:  sourceConfig.SubPath,
					},
				}, nil
			}
		}
	}

	return v1alpha2.ResolvedSourceConfig{
		Git: &v1alpha2.ResolvedGitSource{
			URL:      sourceConfig.Git.URL,
			Revision: sourceConfig.Git.Revision,
			Type:     v1alpha2.Commit,
			SubPath:  sourceConfig.SubPath,
		},
	}, nil
}

func sourceType(reference git2go.RemoteHead) v1alpha2.GitSourceKind {
	switch {
	case strings.HasPrefix(reference.Name, "refs/heads"):
		return v1alpha2.Branch
	case strings.HasPrefix(reference.Name, "refs/tags"):
		return v1alpha2.Tag
	default:
		return v1alpha2.Unknown
	}
}

var refRevParseRules = []string{
	"refs/%s",
	"refs/tags/%s",
	"refs/heads/%s",
	"refs/remotes/%s",
	"refs/remotes/%s/HEAD",
}
