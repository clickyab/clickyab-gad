package ip2location

import (
	"assert"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/fzerorubigd/expand"
)

var (
	fp string
)

type fileMock struct {
	data []byte
	ln   int
}

func newFileMock() (*fileMock, error) {
	f, err := os.Open(fp)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	data, err := ioutil.ReadAll(f)
	if err != nil {
		return nil, err
	}

	return &fileMock{
		data: data,
		ln:   len(data),
	}, nil
}

func (fm *fileMock) ReadAt(b []byte, off int64) (n int, err error) {
	lb := len(b)
	avail := int64(fm.ln) - off
	if avail < int64(lb) {
		return 0, io.EOF
	}

	copy(b, fm.data[off:off+int64(lb)])
	return lb, nil
}

func init() {
	pwd, err := expand.Pwd()
	assert.Nil(err)
	fp = filepath.Join(pwd, "IP-COUNTRY-REGION-CITY-ISP.BIN")
	assert.Nil(open())
}
