package utils

import (
	"bytes"
	"encoding/gob"
)

// InterfaceToByte save interface into byte
func InterfaceToByte(in interface{}) ([]byte, error) {
	buf := &bytes.Buffer{}

	enc := gob.NewEncoder(buf)
	err := enc.Encode(in)
	if err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

// ByteToInterface return object from byte
func ByteToInterface(b []byte, out interface{}) error {
	buf := bytes.NewBuffer(b)

	dnc := gob.NewDecoder(buf)
	return dnc.Decode(out)
}
