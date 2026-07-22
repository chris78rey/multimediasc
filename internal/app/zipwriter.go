package app

import (
	"archive/zip"
	"io"
	"os"
)

type zipWriter struct {
	w *zip.Writer
	f *os.File
}

func newZipWriter(dst io.Writer) *zipWriter {
	return &zipWriter{w: zip.NewWriter(dst)}
}

func (z *zipWriter) addFile(name, srcPath string) error {
	src, err := os.Open(srcPath)
	if err != nil {
		return err
	}
	defer src.Close()
	h, err := z.w.Create(name)
	if err != nil {
		return err
	}
	_, err = io.Copy(h, src)
	return err
}

func (z *zipWriter) close() error {
	return z.w.Close()
}
