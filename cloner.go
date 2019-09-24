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

// File represents a system file (either a folder or an actual file) which is supposed to be cloned.
type File struct {
	info     os.FileInfo
	src, dst string
}

// CloneOptions which will configure a Clone
type CloneOptions interface {
	IncludeHidden() bool
}

func log(s string) {
	fmt.Println(s)
	flogger.Log(s)
}

// Clone clones the directory at scrPath to dstPath using provided options,
// writes the cloned results into 'cloner.txt' file in the root of the srcPath.
func Clone(srcPath, dstPath string, options CloneOptions) error {

	if err := validatePaths(srcPath, dstPath); err != nil {
		return err
	}

	/* preparing log file, into which we'll write which items were cloned */
	logFile := "cloner.txt"
	flogger.PrepareOnPath(filepath.Join(srcPath, logFile))
	defer flogger.Resign()

	log("Preparing files...")
	files, err := Diff(srcPath, dstPath, options)
	if err != nil {
		return err
	}

	totalCount := len(files)

	log(fmt.Sprintf("Cloning %d items.\n", totalCount))
	for i, file := range files {
		fmt.Printf("\rOn %d/%d", i+1, totalCount)

		if err = cper.Cp(file.info, file.src, file.dst); err != nil {
			return err
		}

		flogger.Log(strings.TrimLeft(file.src, srcPath))
	}
	log("\nCompleted successfully!")

	return err
}

// Diff returns a difference between src and dst paths, or error,
// if any occurred while itterating throught directory at src path
func Diff(src, dst string, options CloneOptions) ([]File, error) {
	var files []File

	/* Represents last folder that is supposed to be cloned.

	When we find that a folder does not exist in the dst path, we save it for cloning.
	Then, if that colder contains some files (which is usually the case),
	we'll check if they exists as well - we're sure they don't exists,
	because a folder that's holding them does not.

	So, to reduce the number of items, we'll save a path of that folder,
	and return every time we find a file whose path contains path of that folder.
	*/
	var lastNonExistingFolder string

	/* walk through folder at srcPath */
	err := filepath.Walk(src, func(path string, info os.FileInfo, err error) error {

		/* happens only once, at fist itteration */
		if src == path {
			return nil
		}

		/* first, read a comment above declaration of lastNonExistingFolder

		here we're simply checking if last cloned folder is saved, and if file's path contains
		a path to that folder, i.e, if it's in that folder. */
		if lastNonExistingFolder != "" {
			if strings.Contains(path, lastNonExistingFolder) {
				return nil
			}
			lastNonExistingFolder = ""
		}

		if err != nil {
			return filepath.SkipDir
		}

		if isFileHidden(info) && !options.IncludeHidden() {
			return nil
		}

		// path of the item without src
		internalPath := strings.Replace(path, src, "", -1)

		dstItemPath := filepath.Join(dst, internalPath)
		if fileExistsAtPath(dstItemPath) {
			return nil
		}

		if info.IsDir() {
			lastNonExistingFolder = path
		}

		item := File{info: info, src: path, dst: dstItemPath}
		files = append(files, item)

		return err
	})

	if err == nil {
		return files, nil
	}
	return nil, err
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

/*
Checks if the provided file is hidden.
Currently only works for unix systems, where hidden files begins with "."
*/
func isFileHidden(file os.FileInfo) bool {
	return strings.HasPrefix(file.Name(), ".")
}
