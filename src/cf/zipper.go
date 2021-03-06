package cf

import (
	"archive/zip"
	"errors"
	"fileutils"
	"os"
	"path/filepath"
)

type Zipper interface {
	Zip(dirToZip string, targetFile *os.File) (err error)
}

type ApplicationZipper struct{}

var doNotZipExtensions = []string{".zip", ".war", ".jar"}

func (zipper ApplicationZipper) Zip(dirOrZipFile string, targetFile *os.File) (err error) {
	if shouldNotZip(filepath.Ext(dirOrZipFile)) {
		err = fileutils.CopyPathToWriter(dirOrZipFile, targetFile)
	} else {
		err = writeZipFile(dirOrZipFile, targetFile)
	}
	targetFile.Seek(0, os.SEEK_SET)
	return
}

func shouldNotZip(extension string) (result bool) {
	for _, ext := range doNotZipExtensions {
		if ext == extension {
			return true
		}
	}
	return
}

func writeZipFile(dir string, targetFile *os.File) (err error) {
	isEmpty, err := fileutils.IsDirEmpty(dir)
	if err != nil {
		return
	}
	if isEmpty {
		err = errors.New("Directory is empty")
		return
	}

	writer := zip.NewWriter(targetFile)
	defer writer.Close()

	err = walkAppFiles(dir, func(fileName string, fullPath string) (err error) {
		zipFilePart, err := writer.Create(fileName)
		err = fileutils.CopyPathToWriter(fullPath, zipFilePart)
		return
	})

	return
}
