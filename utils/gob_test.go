package utils

import (
	"testing"

	"github.com/sirupsen/logrus"
)

type testStruct struct {
	ID    int64
	Test  string
	Boogh bool
}

func TestByteToInterface(t *testing.T) {
	data := &testStruct{
		ID:    100,
		Test:  <-ID,
		Boogh: true,
	}

	x, err := InterfaceToByte(data)
	if err != nil {
		logrus.Warn(err)
		t.Fail()
		return
	}

	res := &testStruct{}
	err = ByteToInterface(x, res)
	if err != nil {
		logrus.Warn(err)
		t.Fail()
		return
	}

	if res.ID != data.ID || res.Test != data.Test || res.Boogh != data.Boogh {
		logrus.Error("not same")
		t.Fail()
		return
	}
}
