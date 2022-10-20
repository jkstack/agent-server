package file

import (
	"crypto/md5"
	"io"
	"os"

	"github.com/jkstack/anet"
	"github.com/jkstack/jkframe/compress"
)

const blockSize = 4096

func fillFile(f *os.File, size uint64) error {
	left := size
	dummy := make([]byte, blockSize)
	for left > 0 {
		if left >= blockSize {
			_, err := f.Write(dummy)
			if err != nil {
				return err
			}
			left -= blockSize
			continue
		}
		dummy = make([]byte, left)
		n, err := f.Write(dummy)
		if err != nil {
			return err
		}
		left -= uint64(n)
	}
	return nil
}

func writeFile(f *os.File, data *anet.DownloadData) (int, error) {
	_, err := f.Seek(int64(data.Offset), io.SeekStart)
	if err != nil {
		return 0, err
	}
	dec, err := compress.Decompress(data.Data)
	if err != nil {
		return 0, err
	}
	return f.Write(dec)
}

func md5From(r io.Reader) ([md5.Size]byte, error) {
	var ret [md5.Size]byte
	enc := md5.New()
	_, err := io.Copy(enc, r)
	if err != nil {
		return ret, err
	}
	copy(ret[:], enc.Sum(nil))
	return ret, nil
}
