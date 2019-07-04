package cloner

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"strings"
)

// Clones the directory at scrPath to dstPath.
func Clone(srcPath, dstPath string) {
	err := filepath.Walk(srcPath, func(path string, info os.FileInfo, err error) error {

		// path of the item without `srcPath`
		internalPath := strings.Replace(path, srcPath, "", -1)

		/* if it's the path for te `dstPath` - happens only once, in the first itteration */
		if internalPath == "" {
			return nil
		}

		srcItemPath := pathByAppendingPath(srcPath, internalPath)
		dstItemPath := pathByAppendingPath(dstPath, internalPath)

		if !fileExistsAtPath(dstItemPath) {
			if info.IsDir() {
				return copyDir(srcItemPath, dstItemPath)
			} else {
				return copyFile(srcItemPath, dstItemPath)
			}
		}

		return nil
	})

	if err != nil {
		fmt.Println(err)
	}
}

/* Returns boolean indicating whether a file exists at provided path. */
func fileExistsAtPath(path string) bool {
	if _, err := os.Stat(path); !os.IsNotExist(err) {
		return true
	}
	return false
}

/*
Appens string p2 to string p1, while appending `/` to end of p1 (if isn't already),
and removing `/` from start of p2 (if exists).
*/
func pathByAppendingPath(p1, p2 string) string {
	pathDelimiter := "/"

	if !strings.HasSuffix(p1, pathDelimiter) {
		p1 = p1 + pathDelimiter
	}

	if strings.HasPrefix(p2, pathDelimiter) {
		strings.TrimPrefix(p2, pathDelimiter)
	}

	return p1 + p2
}

/* COPYING */

// copyDir copies a whole directory recursively
func copyDir(src string, dst string) error {
	var err error
	var fds []os.FileInfo
	var srcInfo os.FileInfo

	if srcInfo, err = os.Stat(src); err != nil {
		return err
	}

	if err = os.MkdirAll(dst, srcInfo.Mode()); err != nil {
		return err
	}

	if fds, err = ioutil.ReadDir(src); err != nil {
		return err
	}
	for _, fd := range fds {
		srcfp := path.Join(src, fd.Name())
		dstfp := path.Join(dst, fd.Name())

		if fd.IsDir() {
			if err = copyDir(srcfp, dstfp); err != nil {
				fmt.Println(err)
			}
		} else {
			if err = copyFile(srcfp, dstfp); err != nil {
				fmt.Println(err)
			}
		}
	}
	return nil
}

// copyFile copies a single file from src to dst
func copyFile(src, dst string) error {
	var err error
	var srcfd *os.File
	var dstfd *os.File
	var srcInfo os.FileInfo

	if srcfd, err = os.Open(src); err != nil {
		return err
	}
	defer srcfd.Close()

	if dstfd, err = os.Create(dst); err != nil {
		return err
	}
	defer dstfd.Close()

	if _, err = io.Copy(dstfd, srcfd); err != nil {
		return err
	}
	if srcInfo, err = os.Stat(src); err != nil {
		return err
	}
	return os.Chmod(dst, srcInfo.Mode())
}
