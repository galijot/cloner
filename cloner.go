package cloner

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/galijot/cloner/cper"
	"github.com/galijot/cloner/flogger"
)

type CloneOptions interface {
	IncludeHidden() bool
}

/*
Clones the directory at scrPath to dstPath using provided options,
& writes the cloned results into 'cloner.txt' file in the root of the srcPath.
*/
func Clone(srcPath, dstPath string, options CloneOptions) error {

	err := validatePaths(srcPath, dstPath)
	if err != nil {
		return err
	}

	logFile := "cloner.txt"
	flogger.PrepareOnPath(filepath.Join(srcPath, logFile))
	defer flogger.Resign()

	return filepath.Walk(srcPath, func(path string, info os.FileInfo, err error) error {

		if err != nil {
			return err
		}
		if !options.IncludeHidden() && isFileHidden(info) {
			return nil
		}

		// path of the item without `srcPath`
		internalPath := strings.Replace(path, srcPath, "", -1)

		/* if it's the path for te `dstPath` - happens only once, in the first itteration */
		if internalPath == "" {
			return nil
		}

		dstItemPath := filepath.Join(dstPath, internalPath)
		if fileExistsAtPath(dstItemPath) {
			return nil
		}

		srcItemPath := filepath.Join(srcPath, internalPath)

		if info.IsDir() {
			err = cper.Dir(srcItemPath, dstItemPath)

			if err == nil {
				log := fmt.Sprintf("🗂 Cloned folder %v to %v", info.Name(), dstItemPath)
				err = flogger.Log(log)
			}
		} else {
			err = cper.File(srcItemPath, dstItemPath)

			if err == nil {
				log := fmt.Sprintf("📄 Cloned item %v to %v\n", info.Name(), dstItemPath)
				err = flogger.Log(log)
			}
		}

		return err
	})
}

// checks if paths are valid, and if they're not returns an error describing the issue.
func validatePaths(src, dst string) error {
	if src == "" {
		return errors.New("Can't clone from empty source path")
	}

	if dst == "" {
		return errors.New("Can't clone to empty destination path")
	}

	if !fileExistsAtPath(src) {
		return errors.New("Source path does not exist")
	}

	if !fileExistsAtPath(dst) {
		return errors.New("Destination path does not exist")
	}

	return nil
}

/* Returns boolean indicating whether a file exists at provided path. */
func fileExistsAtPath(path string) bool {
	if _, err := os.Stat(path); !os.IsNotExist(err) {
		return true
	}
	return false
}

/*
Checks if the provided file is hidden.
Currently only works for unix systems, where hidden files begins with "."
*/
func isFileHidden(file os.FileInfo) bool {
	return strings.HasPrefix(file.Name(), ".")
}
