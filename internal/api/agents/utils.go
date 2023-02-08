package agents

import (
	"bytes"
	"crypto/md5"
	"errors"
	"hash"
	"io"
	"os"
	"path"

	"github.com/jkstack/jkframe/compress"
)

var errInvalidOffset = errors.New("invalid offset")

type cacheFile struct {
	f        *os.File
	enc      hash.Hash
	offset   uint64
	w        io.Writer
	checksum [md5.Size]byte
}

func newCacheFile(cacheDir string, checksum [md5.Size]byte) (*cacheFile, error) {
	tmp := path.Join(cacheDir, "download")
	os.MkdirAll(tmp, 0755)
	f, err := os.CreateTemp(tmp, "log")
	if err != nil {
		return nil, err
	}
	defer f.Close()
	if err != nil {
		f.Close()
		os.Remove(f.Name())
		return nil, err
	}
	f.Close()
	f, err = os.OpenFile(f.Name(), os.O_RDWR|os.O_TRUNC, 0644)
	if err != nil {
		f.Close()
		os.Remove(f.Name())
		return nil, err
	}
	enc := md5.New()
	return &cacheFile{
		f:        f,
		enc:      enc,
		w:        io.MultiWriter(f, enc),
		checksum: checksum,
	}, nil
}

func (f *cacheFile) Close() error {
	return f.f.Close()
}

func (f *cacheFile) Write(offset uint64, data string) (int, error) {
	if offset != f.offset {
		return 0, errInvalidOffset
	}
	dt, err := compress.Decompress(data)
	if err != nil {
		return 0, err
	}
	return f.w.Write(dt)
}

func (f *cacheFile) Name() string {
	return f.f.Name()
}

func (f *cacheFile) Remove() {
	f.Close()
	os.Remove(f.f.Name())
}

func (f *cacheFile) check() bool {
	return bytes.Equal(f.enc.Sum(nil), f.checksum[:])
}
