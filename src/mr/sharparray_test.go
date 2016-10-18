package mr

import (
	"fmt"
	"testing"
)

func TestSharpArray(t *testing.T) {
	data := "#11#111#55#99#1#987#"

	var (
		pa SharpArray
		ps SharpArray = SharpArray{1, 11, 55, 99, 111, 987}
	)
	pa.Scan(data)

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
