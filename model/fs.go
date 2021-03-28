package model

import (
	"fmt"
	"os"
)

func PrepareDirectoryTree(path string) error {
	return os.MkdirAll(path, 0777)
}

func CleanUpDirectoryTree(err *error, dir string) error {
	if err != nil {
		fmt.Println("removing", dir)
		return nil
		//return os.RemoveAll(dir)
	}
	return nil
}
