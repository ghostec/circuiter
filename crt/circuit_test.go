package crt_test

import (
	"fmt"
	"testing"

	"github.com/ghostec/circuiter/crt"
)

func Test(t *testing.T) {
	circuit := crt.NewCircuit(2, 1)
	circuit.AddPart(crt.OR, 0, 1)
	circuit.SetOutputInput(0, 2)
	fmt.Printf("%#v\n", circuit.Evaluate([]crt.Bit{crt.BitOne, crt.BitZero}))
}
