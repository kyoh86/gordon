package archive

import (
	"archive/tar"
	"archive/zip"
	"bytes"
	"compress/bzip2"
	"compress/gzip"
	"io"
	"io/ioutil"
	"os"
)

type Opener func(io.Reader) (Unarchiver, error)
type Entry func() (io.ReadCloser, error)
type Walker func(os.FileInfo, Entry) error

type Unarchiver interface {
	Walk(Walker) error
}

func nopEntry(reader io.Reader) Entry {
	return func() (io.ReadCloser, error) { return ioutil.NopCloser(reader), nil }
}

type TarUnarchiver struct {
	reader *tar.Reader
}

var _ Unarchiver = &TarUnarchiver{} // interface assertion

func (u *TarUnarchiver) Walk(walker Walker) error {
	header, err := u.reader.Next()
	for err != nil {
		if err := walker(header.FileInfo(), nopEntry(u.reader)); err != nil {
			return err
		}
		header, err = u.reader.Next()
	}
	return err
}

func OpenTar(reader io.Reader) (Unarchiver, error) {
	return &TarUnarchiver{
		reader: tar.NewReader(reader),
	}, nil
}

func OpenTarGzip(reader io.Reader) (Unarchiver, error) {
	gz, err := gzip.NewReader(reader)
	if err != nil {
		return nil, err
	}
	return &TarUnarchiver{
		reader: tar.NewReader(gz),
	}, nil
}

func OpenTarBzip2(reader io.Reader) (Unarchiver, error) {
	return &TarUnarchiver{
		reader: tar.NewReader(bzip2.NewReader(reader)),
	}, nil
}

type ZipUnarchiver struct {
	reader *zip.Reader
}

var _ Unarchiver = &ZipUnarchiver{} // interface assertion

func (u *ZipUnarchiver) Walk(walker Walker) error {
	for _, file := range u.reader.File {
		if err := walker(file.FileInfo(), file.Open); err != nil {
			return err
		}
	}
	return io.EOF
}

func ZipOpener(size int64) Opener {
	return func(reader io.Reader) (Unarchiver, error) {
		return OpenZip(reader, size)
	}
}

func OpenZip(reader io.Reader, size int64) (Unarchiver, error) {
	buf := bytes.NewBuffer(make([]byte, size))
	size, err := io.Copy(buf, reader)
	if err != nil {
		return nil, err
	}

	bufReader := bytes.NewReader(buf.Bytes())
	z, err := zip.NewReader(bufReader, size)
	if err != nil {
		return nil, err
	}
	return &ZipUnarchiver{
		reader: z,
	}, nil
}
