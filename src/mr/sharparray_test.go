package mr

import (
	"fmt"
	"testing"

	"github.com/Sirupsen/logrus"
)

func TestSharpArray(t *testing.T) {
	data := "#11#111#55#99#1#987#"

	var (
		pa SharpArray
		ps = SharpArray{1, 11, 55, 99, 111, 987}
	)
	err := pa.Scan(data)
	if err != nil {
		logrus.Info(err)
	}

	fmt.Print(pa)

	if len(pa) != len(ps) {
		fmt.Print("S2s")

		t.Fail()
		return
	}
	fmt.Print(pa)
	for i := range pa {
		if pa[i] != ps[i] {
			fmt.Print("Ssww")

			t.Fail()
			return
		}
	}

	if !pa.Has(11) {
		fmt.Print("Ss")

		t.Fail()
		return
	}

	if pa.Has(44) {
		fmt.Print("Ssedes")

		t.Fail()
		return
	}
}
