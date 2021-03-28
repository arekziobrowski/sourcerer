package model

import (
	"strings"

	"github.com/pkg/errors"
)

type Source struct {
	Origin       string
	Hash         string
	Organization string
	Repository   string
}

func ToSource(src string) (*Source, error) {
	origin, hash, err := extractOriginAndHash(src)
	if err != nil {
		return nil, errors.Wrap(err, "invalid input for origin and hash")
	}
	organization, repo, err := extractOrganizationAndRepo(origin)
	if err != nil {
		return nil, errors.Wrap(err, "invalid input for organization and name")
	}
	return &Source{
		Origin:       origin,
		Hash:         hash,
		Organization: organization,
		Repository:   repo,
	}, nil
}

func extractOriginAndHash(src string) (string, string, error) {
	split := strings.Split(src, " ")
	if len(split) != 2 {
		return "", "", errors.Errorf("invalid origin and hash definition: %s", src)
	}
	return split[0], split[1], nil
}

func extractOrganizationAndRepo(origin string) (string, string, error) {
	protocolSplit := strings.Split(origin, ":")
	if len(protocolSplit) != 2 {
		return "", "", errors.Errorf("invalid URL structure: %s", origin)
	}

	organizationAndRepo := strings.TrimSuffix(protocolSplit[1], ".git")
	split := strings.Split(organizationAndRepo, "/")
	if len(split) != 2 {
		return "", "", errors.Errorf("invalid organization and repository name structure: %s", organizationAndRepo)
	}
	return split[0], split[1], nil
}
