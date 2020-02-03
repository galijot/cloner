package cloner

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/galijot/cloner/cper"
	"github.com/galijot/cloner/differ"
	"github.com/galijot/cloner/flogger"
)

func log(s ...interface{}) {
	fmt.Println(s...)
	flogger.Log(fmt.Sprintf("%v", s))
}

// Clone clones the directory at scrPath to dstPath,
// writes the cloned results into 'cloner.txt' file in the root of the srcPath.
func Clone(srcPath, dstPath string) error {

	if err := validatePaths(srcPath, dstPath); err != nil {
		return err
	}

	/* preparing log file, into which we'll write which items were cloned */
	logFile := "cloner.txt"
	flogger.PrepareOnPath(filepath.Join(srcPath, logFile))
	defer flogger.Resign()

	log("Preparing files...")
	result, err := differ.Diff(srcPath, dstPath)
	if err != nil {
		return err
	}
	files := result.SrcToDstItems

	totalSize := differ.SizeOfItems(files)
	fmt.Println(totalSize)
	totalCount := len(files)

	log(fmt.Sprintf("Cloning %d items.\n", totalCount))
	for i, file := range files {
		fmt.Printf("\rOn %d/%d", i+1, totalCount)

		if err = cper.Cp(file.Info, file.Src, file.Dst); err != nil {
			return err
		}

		flogger.Log(strings.TrimLeft(file.Src, srcPath))
	}
	log("\nCompleted successfully!")

	return err
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

	if !isDirectory(src) {
		return errors.New("Source path must point to a folder, not to a file")
	}

	if !isDirectory(dst) {
		return errors.New("Destination path must point to a folder, not to a file")
	}

	return nil
}

func isDirectory(path string) bool {
	fileInfo, err := os.Stat(path)
	if err != nil {
		return false
	}
	return fileInfo.IsDir()
}

/* Returns boolean indicating whether a file exists at provided path. */
func fileExistsAtPath(path string) bool {
	if _, err := os.Stat(path); !os.IsNotExist(err) {
		return true
	}
	return false
}
