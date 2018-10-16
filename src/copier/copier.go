package copier

import (
	"errors"
	"fmt"
	"io"
	"os"
	"strings"
)

func JoinPath(path ...string) string {
	return strings.Replace(strings.Join(path, "/"), "//", "/", -1)
}

func CopyDir(src, dest string) error {
	srcInfo, err := os.Stat(src)
	if err != nil {
		return err
	}

	err = os.MkdirAll(dest, srcInfo.Mode())
	if err != nil {
		return err
	}

	directory, err := os.Open(src)
	if err != nil {
		return err
	}

	objects, err := directory.Readdir(-1)
	if err != nil {
		return err
	}

	for _, obj := range objects {
		srcFileName := JoinPath(src, obj.Name())
		destFileName := JoinPath(dest, obj.Name())

		if srcFileName == destFileName {
			return errors.New(fmt.Sprintf("failed to copy dir. same name already exist: %s\n", srcFileName))
		}

		if obj.IsDir() {
			err = CopyDir(srcFileName, destFileName)
			if err != nil {
				return err
			}
		} else {
			err = CopyFile(srcFileName, destFileName)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func CopyFile(src, dest string) error {
	file, err := os.Open(src)
	if err != nil {
		return err
	}
	defer file.Close()

	destfile, err := os.Create(dest)
	if err != nil {
		return err
	}
	defer destfile.Close()

	_, err = io.Copy(destfile, file)
	if err == nil {
		srcInfo, err := os.Stat(src)
		if err != nil {
			err = os.Chmod(dest, srcInfo.Mode())
		}
	}
	return nil
}
