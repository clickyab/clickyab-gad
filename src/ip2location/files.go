package ip2location

import (
	"assert"
	"os"
	"path/filepath"

	"github.com/fzerorubigd/expand"
)

var (
	fp string
)

type fileMock struct {
	f *os.File
}

func newFileMock() (*fileMock, error) {
	f, err := os.Open(fp)
	if err != nil {
		return nil, err
	}

	return &fileMock{
		f: f,
	}, nil
}

func (fm *fileMock) ReadAt(b []byte, off int64) (n int, err error) {
	return fm.f.ReadAt(b, off)
}

func init() {
	pwd, err := expand.Pwd()
	assert.Nil(err)
	fp = filepath.Join(pwd, "IP-COUNTRY-REGION-CITY-ISP.BIN")
	assert.Nil(open())
}
