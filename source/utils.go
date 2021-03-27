package source

import (
	"fmt"
	"os"
	"strings"

	"github.com/pkg/errors"
)

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

func prepare(path string) error {
	return os.MkdirAll(path, 0777)
}

func cleanUpIfError(err error, dir string) error {
	if err != nil {
		fmt.Println("removing", dir)
		return nil
		//return os.RemoveAll(dir)
	}
	return nil
}
