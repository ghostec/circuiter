package crt

import (
	"errors"
)

type Circuit struct {
	nInputs, nOutputs int
	parts, outputs    []CircuitPartWithInputs
}

func NewCircuit(nInputs, nOutputs int) *Circuit {
	var parts []CircuitPartWithInputs

	for i := 0; i < nInputs; i++ {
		parts = append(parts, CircuitPartWithInputs{Part: Buffer})
	}

	var outputs []CircuitPartWithInputs

	for i := 0; i < nOutputs; i++ {
		outputs = append(outputs, CircuitPartWithInputs{Part: Buffer})
	}

	return &Circuit{nInputs: nInputs, nOutputs: nOutputs, parts: parts, outputs: outputs}
}

func (c *Circuit) Evaluate(input []Bit) []Bit {
	results := make([]Bit, len(c.parts)+len(c.outputs))

	for i, bit := range input {
		results[i] = bit
	}

	for i := c.nInputs; i < len(c.parts); i++ {
		var partInput []Bit
		for _, iidx := range c.parts[i].Inputs {
			partInput = append(partInput, results[iidx])
		}

		results[i] = c.parts[i].Part.Value(partInput...)
	}

	for i := range c.outputs {
		results[len(c.parts)+i] = results[c.outputs[i].Inputs[0]]
	}

	return results[len(c.parts):]
}

func (c *Circuit) Clone() *Circuit {
	clone := NewCircuit(c.nInputs, c.nOutputs)

	clone.parts = make([]CircuitPartWithInputs, len(c.parts))
	for i, cp := range c.parts {
		clone.parts[i] = cp.Clone()
	}

	clone.outputs = make([]CircuitPartWithInputs, len(c.outputs))
	for i, co := range c.outputs {
		clone.outputs[i] = co.Clone()
	}

	return clone
}

func (c *Circuit) NInputs() int  { return c.nInputs }
func (c *Circuit) NOutputs() int { return c.nOutputs }

func (c *Circuit) SetOutputInput(oidx, iidx int) {
	c.outputs[oidx].Inputs = []int{iidx}
}

func (c *Circuit) AddPart(part CircuitPart, inputs ...int) error {
	for _, idx := range inputs {
		if idx >= len(c.parts) {
			return errors.New("circuit: invalid part inputs")
		}
	}

	c.parts = append(c.parts, CircuitPartWithInputs{
		Part:   part,
		Inputs: inputs,
	})

	return nil
}

func (c *Circuit) GetParts() (dst []CircuitPartWithInputs) {
	dst = make([]CircuitPartWithInputs, len(c.parts))
	copy(dst, c.parts)
	return
}

func (c *Circuit) SetParts(parts []CircuitPartWithInputs) {
	c.parts = parts
}

func (c *Circuit) GetOutputs() (dst []CircuitPartWithInputs) {
	dst = make([]CircuitPartWithInputs, len(c.outputs))
	copy(dst, c.outputs)
	return
}

func (c *Circuit) SetOutputs(outputs []CircuitPartWithInputs) {
	c.outputs = outputs
}

func (c *Circuit) NParts() int {
	return len(c.parts)
}

type CircuitPartWithInputs struct {
	Part CircuitPart
	// inputs indexes
	// must be less than the index of the part itself
	Inputs []int
}

func (p CircuitPartWithInputs) Clone() CircuitPartWithInputs {
	inputs := make([]int, len(p.Inputs))
	copy(inputs, p.Inputs)

	return CircuitPartWithInputs{
		Part:   p.Part,
		Inputs: inputs,
	}
}

type CircuitPart interface {
	Value(...Bit) Bit
	NInputs() int
}

type Bit uint8

const (
	BitZero Bit = 0
	BitOne  Bit = 1
)

func Bits(vals ...uint8) (ret []Bit) {
	for _, val := range vals {
		bit := BitZero
		if val > 0 {
			bit = BitOne
		}
		ret = append(ret, bit)
	}
	return
}

type Gate struct{}

type _OR Gate
type _AND Gate
type _NOT Gate
type _XOR Gate
type _Buffer Gate

var (
	OR     = _OR{}
	AND    = _AND{}
	NOT    = _NOT{}
	XOR    = _XOR{}
	Buffer = _Buffer{}
)

func (_Buffer) Value(inputs ...Bit) Bit {
	if len(inputs) != 1 {
		panic("Buffer: must have 1 input")
	}

	return inputs[0]
}

func (_Buffer) NInputs() int { return 1 }

func (_OR) Value(inputs ...Bit) Bit {
	if len(inputs) != 2 {
		panic("OR: must have 2 inputs")
	}

	a, b := inputs[0], inputs[1]

	if a == BitZero && b == BitZero {
		return BitZero
	}

	return BitOne
}

func (_OR) NInputs() int { return 2 }

func (_AND) Value(inputs ...Bit) Bit {
	if len(inputs) != 2 {
		panic("AND: must have 2 inputs")
	}

	a, b := inputs[0], inputs[1]

	if a == BitOne && b == BitOne {
		return BitOne
	}

	return BitZero
}

func (_AND) NInputs() int { return 2 }

func (_NOT) Value(inputs ...Bit) Bit {
	if len(inputs) != 1 {
		panic("NOT: must have 1 input")
	}

	a := inputs[0]

	if a == BitOne {
		return BitZero
	}

	return BitOne
}

func (_NOT) NInputs() int { return 1 }

func (_XOR) Value(inputs ...Bit) Bit {
	if len(inputs) != 2 {
		panic("XOR: must have 2 inputs")
	}

	a, b := inputs[0], inputs[1]

	if a == BitZero && b == BitZero {
		return BitZero
	}

	if a == BitOne && b == BitOne {
		return BitZero
	}

	return BitOne
}

func (_XOR) NInputs() int { return 2 }
