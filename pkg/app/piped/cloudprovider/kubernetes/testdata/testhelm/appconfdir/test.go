package main

import (
	"fmt"
	"os"
)

// resolveSymlink resolves symbolic link.
func resolveSymlink(path string) (string, error) {
	fmt.Println("path:", path)
	lstat, err := os.Lstat(path)
	if err != nil {
		if _, ok := err.(*os.PathError); ok {
			return path, nil
		}
		return "", err
	}

	if lstat.Mode()&os.ModeSymlink == os.ModeSymlink {
		resolved, err := os.Readlink(path)
		if err != nil {
			return "", err
		}

		return resolveSymlink(resolved)
	}

	return path, nil
}

func main() {
	resolveSymlink("./invalid-symlink")
}
