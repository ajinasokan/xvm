package main

import (
	"fmt"
	"os"
	"path/filepath"
)

func createSymlink(dir string, self string, cmd string) {
	linkPath := filepath.Join(dir, cmd)
	err := os.Mkdir(dir, 0755)
	if err != nil && !os.IsExist(err) {
		exitWithError("unable create", dir, "dir", err)
	}
	os.Remove(linkPath)
	err = os.Symlink(self, linkPath)
	if err != nil {
		exitWithError("unable to create symlink for", cmd, err)
	}
	fmt.Println("symlink", self+"->"+cmd, "created at", linkPath)
}

func removeSymlink(dir string, cmd string) {
	linkPath := filepath.Join(dir, cmd)
	err := os.Remove(linkPath)
	if err != nil {
		exitWithError("unable to remove symlink", err)
	}
	fmt.Println("symlink", linkPath, "removed")
}
