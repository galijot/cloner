package differ

import (
	"os"
	"path/filepath"
	"strings"
)

// Result represents a result from differ's Diff function.
type Result struct {

	// Items that exists in src but don't in dst.
	SrcToDstItems []File

	// Items that exists in dst but don't in src.
	DstToSrcItems []File
}

// Diff returns a difference between directories at src and dst paths, or error,
// if any occurred while itterating throught directories.
func Diff(src, dst string) (*Result, error) {

	// Retrieving the files that exists in src but don't in dst folder & returning if error occurred.
	srcToDst, err1 := _diff(src, dst)
	if err1 != nil {
		return nil, err1
	}

	// Now retrieving the files that exists in dst but don't in src folder & returning if error occurred.
	dstToSrc, err2 := _diff(dst, src)
	if err2 != nil {
		return nil, err2
	}

	// Preparing a Diff result.
	result := Result{srcToDst, dstToSrc}
	return &result, nil
}

func _diff(src, dst string) ([]File, error) {
	var files []File

	/* Represents last folder that is supposed to be cloned.

	When we find that a folder does not exist in the dst path, we save it for cloning.
	Then, if that folder contains some files (which is usually the case),
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

		// path of the item without src
		internalPath := strings.Replace(path, src, "", -1)

		dstItemPath := filepath.Join(dst, internalPath)
		if fileExistsAtPath(dstItemPath) {
			return nil
		}

		if info.IsDir() {
			lastNonExistingFolder = path
		}

		file := File{info, path, dstItemPath}
		files = append(files, file)

		return err
	})

	if err == nil {
		return files, nil
	}
	return nil, err
}

/* Returns boolean indicating whether a file exists at provided path. */
func fileExistsAtPath(path string) bool {
	if _, err := os.Stat(path); !os.IsNotExist(err) {
		return true
	}
	return false
}
