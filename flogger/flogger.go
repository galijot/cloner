package flogger

import (
	"os"
)

// Writes the provided txt to file at path, and creates file if doesn't exist.
func LogToFile(txt string, path string) error {
	var file *os.File
	var err error

	file, err = os.OpenFile(path, os.O_APPEND|os.O_WRONLY, os.ModeAppend)

	if err != nil {
		file, err = os.Create(path)

		if err != nil {
			return err
		}
	}

	defer file.Close()
	_, err = file.WriteString("\n" + txt)
	return err
}
