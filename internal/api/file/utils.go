package file

import (
	"bytes"
	"compress/gzip"
	"crypto/md5"
	"encoding/base64"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"

	"github.com/jkstack/anet"
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
	dec, err := decodeData(data.Data)
	if err != nil {
		return 0, err
	}
	return f.Write(dec)
}

func encodeData(data []byte) string {
	var str string
	if strings.Contains(http.DetectContentType(data), "text/plain") {
		var buf bytes.Buffer
		w := gzip.NewWriter(&buf)
		if _, err := w.Write(data); err == nil {
			w.Close()
			str = "$1$" + base64.StdEncoding.EncodeToString(buf.Bytes())
		}
	}
	if len(str) == 0 {
		str = "$0$" + base64.StdEncoding.EncodeToString(data)
	}
	return str
}

func decodeData(str string) ([]byte, error) {
	switch {
	case strings.HasPrefix(str, "$0$"):
		str := strings.TrimPrefix(str, "$0$")
		return base64.StdEncoding.DecodeString(str)
	case strings.HasPrefix(str, "$1$"):
		str := strings.TrimPrefix(str, "$1$")
		b64 := base64.NewDecoder(base64.StdEncoding, strings.NewReader(str))
		r, err := gzip.NewReader(b64)
		if err != nil {
			return nil, err
		}
		var buf bytes.Buffer
		_, err = io.Copy(&buf, r)
		if err != nil {
			return nil, err
		}
		return buf.Bytes(), nil
	default:
		return nil, fmt.Errorf("invalid data: %s", str)
	}
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
