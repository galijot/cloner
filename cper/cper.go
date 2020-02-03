package cper

import (
	"io"
	"io/ioutil"
	"os"
	"path"
)

// Cp copies a file form src to dst path
func Cp(file os.FileInfo, src, dst string) error {
	if file.IsDir() {
		return Dir(src, dst)
	}
	return File(src, dst)
}

/*
cper (copy paster) is a package for file/directory copying.
*/

// Dir copies a whole directory recursively, and return an error, if any occurs.
func Dir(src string, dst string) error {
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
			if err = Dir(srcfp, dstfp); err != nil {
				return err
			}
		} else {
			if err = File(srcfp, dstfp); err != nil {
				return err
			}
		}
	}
	return nil
}

// File copies a single file from src to dst, and returns an error, if any occurs.
func File(src, dst string) error {
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
