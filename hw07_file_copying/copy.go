package main

import (
	"errors"
	"io"
	"log"
	"os"

	"github.com/cheggaaa/pb"
)

var (
	ErrUnsupportedFile       = errors.New("unsupported file")
	ErrOffsetExceedsFileSize = errors.New("offset exceeds file size")
	ErrFromFileIsEmpty       = errors.New("source file is empty")
)

func Copy(fromPath string, toPath string, offset, limit int64) error {
	fromFile, err := os.Open(fromPath)
	if err != nil {
		return err
	}

	defer fromFile.Close()

	fromInfo, err := fromFile.Stat()
	if err != nil {
		return err
	}

	switch {
	case !fromInfo.Mode().IsRegular():
		return ErrUnsupportedFile
	case fromInfo.Size() == 0:
		return ErrFromFileIsEmpty
	case fromInfo.Size() < offset:
		return ErrOffsetExceedsFileSize
	}

	toFile, err := os.Create(toPath)
	if err != nil {
		return err
	}

	defer toFile.Close()

	o, err := fromFile.Seek(offset, io.SeekStart)
	if err != nil || o != offset {
		return err
	}

	if limit == 0 || limit+offset > fromInfo.Size() {
		limit = fromInfo.Size() - offset
	}

	progress := pb.New64(limit).SetUnits(pb.U_BYTES).Start()
	proxyReader := progress.NewProxyReader(fromFile)

	_, err = io.CopyN(toFile, proxyReader, limit)
	if err != nil {
		log.Fatal(err)
	}

	progress.Finish()

	return nil
}
