package differ

import (
	"io/ioutil"
	"os"
	"path/filepath"
)

// File represents a file, or folder, that exists in one dir and doesn't in another.
type File struct {
	// Underlying fileInfo.
	Info os.FileInfo

	// Source path of the item.
	Src string

	// Destination path of the idem, contained in the destination folder.
	Dst string
}

// Size of the file (in bytes).
func (f File) Size() int64 {

	if f.Info.IsDir() {
		return sizeOfDir(f.Src)
	}
	return f.Info.Size()
}

func sizeOfDir(dir string) int64 {
	files, err := ioutil.ReadDir(dir)
	if err != nil {
		return 0
	}

	var totalSize int64

	for _, file := range files {
		if file.IsDir() {
			subDir := filepath.Join(dir, file.Name())
			totalSize += sizeOfDir(subDir)
		} else {
			totalSize += file.Size()
		}
	}

	// var kilobytes int64
	// kilobytes = (totalSize / 1024)

	// var megabytes float64
	// megabytes = float64(kilobytes / 1024) // cast to type float64

	// // i need to calculate the size for folder's sub-folders & files
	// fmt.Println(file.info.Name(), "-", file.info.Size(), "b ", totalSize, "b ", kilobytes, "Kb ", megabytes, "MB ")
	return totalSize
}

// SizeOfItems returns the size, of the provided items, in kilobytes.
func SizeOfItems(items []File) int64 {
	var size int64
	for _, item := range items {
		size += item.Size()
	}
	return size
}
