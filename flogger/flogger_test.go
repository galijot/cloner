package flogger

import (
	"os"
	"path"
	"strconv"
	"testing"
)

func TestLog(t *testing.T) {
	var dir string
	var err error

	dir, err = os.Getwd()

	if err != nil {
		t.Errorf("Failed to get wd with error: %v", err)
	}

	logFile := "cloner_test.txt"
	filePath := path.Join(dir, logFile)

	PrepareOnPath(filePath)

	for i := 0; i < 100000; i++ {
		txt := strconv.Itoa(i) + "\n"

		err = LogToFile(txt, filePath)

		if err != nil {
			t.Error(err)
		}
	}

	err = os.Remove(filePath)
	if err != nil {
		t.Error(err)
	}

	Resign()
}
