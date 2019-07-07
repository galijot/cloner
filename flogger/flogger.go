package flogger

import (
	"errors"
	"os"
)

/*
flogger (file logger) is a package that writes a text into a file.
*/

var file *os.File

func PrepareOnPath(path string) error {
	if file != nil {
		return errors.New("File already prepared")
	}
	var err error

	file, err = os.OpenFile(path, os.O_APPEND|os.O_WRONLY, os.ModeAppend)

	if err != nil {
		file, err = os.Create(path)
	}

	return err
}

// Writes the provided txt to file previously created with `PrepareOnPath`.
func Log(txt string) error {
	if file == nil {
		return errors.New("File is not prepared. Call `PrepareOnPath` first. Also, don't forget to `Resign`.")
	}
	_, err := file.WriteString("\n" + txt)
	return err
}

func Resign() {
	file.Close()
	file = nil
}
