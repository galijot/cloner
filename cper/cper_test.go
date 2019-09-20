package cper

import (
	"testing"

	"fmt"
	"os"
	"path/filepath"
	"strconv"

	"github.com/galijot/cloner/flogger"
)

/* name of a temp file which is created in */
const fileName = "test-file.txt"

func TestDir(t *testing.T) {
	var err error
	var ex string

	ex, err = os.Executable()
	if err != nil {
		t.Error(err)
	}

	pathOfSelf := filepath.Dir(ex)

	var selfInfo os.FileInfo
	if selfInfo, err = os.Stat(pathOfSelf); err != nil {
		t.Error(err)
	}

	srcDir := filepath.Join(pathOfSelf, "test-srcDir")
	dstDir := filepath.Join(pathOfSelf, "test-dstDir")

	if err = os.Mkdir(srcDir, selfInfo.Mode()); err != nil {
		t.Error(err)
	}
	fmt.Println(srcDir)
	writeToFileAtPath(srcDir)

	if err = Dir(srcDir, dstDir); err != nil {
		t.Error(err)
	}

	dstFile := filepath.Join(dstDir, fileName)
	if !fileExistsAtPath(dstFile) {
		t.Error("File not copied")
	}
}

func writeToFileAtPath(path string) {
	filePath := filepath.Join(path, fileName)
	flogger.PrepareOnPath(filePath)
	for i := 0; i < 50; i++ {
		flogger.Log(strconv.Itoa(i))
	}
	flogger.Resign()
}

/* Returns boolean indicating whether a file exists at provided path. */
func fileExistsAtPath(path string) bool {
	if _, err := os.Stat(path); !os.IsNotExist(err) {
		return true
	}
	return false
}
