package utils

import (
	"bytes"
	"compress/zlib"
	"encoding/gob"

	"github.com/sirupsen/logrus"
)

// InterfaceToByte save interface into byte
func InterfaceToByte(in interface{}) ([]byte, error) {
	buf := &bytes.Buffer{}

	comp := zlib.NewWriter(buf)
	enc := gob.NewEncoder(comp)
	err := enc.Encode(in)
	if err != nil {
		return nil, err
	}
	if err := comp.Close(); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

// ByteToInterface return object from byte
func ByteToInterface(b []byte, out interface{}) error {
	buf := bytes.NewBuffer(b)

	decomp, err := zlib.NewReader(buf)
	if err != nil {
		return err
	}
	defer func() {
		if err := decomp.Close(); err != nil {
			logrus.Debug("couldn't close zlib file", err.Error())
		}
	}()
	dnc := gob.NewDecoder(decomp)
	return dnc.Decode(out)
}
