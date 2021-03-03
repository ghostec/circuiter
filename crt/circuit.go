package crt

import (
	"encoding/json"
	"errors"
)

type Circuit struct {
	nInputs, nOutputs int
	parts, outputs    []CircuitPartWithConnections
}

func NewCircuit(nInputs, nOutputs int) *Circuit {
	var parts []CircuitPartWithConnections

	for i := 0; i < nInputs; i++ {
		parts = append(parts, CircuitPartWithConnections{Part: Buffer})
	}

	var outputs []CircuitPartWithConnections

	for i := 0; i < nOutputs; i++ {
		outputs = append(outputs, CircuitPartWithConnections{Part: Buffer})
	}

	return &Circuit{nInputs: nInputs, nOutputs: nOutputs, parts: parts, outputs: outputs}
}

// TODO: override MarshalJSON interface of Circuit
// and of Part (_OR, etc)
// then
// 1. try to solve sum2bits
// 2. try to use sum2bits circuit as a part
// 3. try to solve sum4bits
// 4. try new mate logic
type circuit struct {
	NInputs, NOutputs int
	Parts, Outputs    []CircuitPartWithConnections
}

func (c *Circuit) Save(path string) error {
	cc := circuit{
		NInputs:  c.nInputs,
		NOutputs: c.nOutputs,
		Parts:    c.parts,
		Outputs:  c.outputs,
	}
	b, err := json.Marshal(cc)
	if err != nil {
		return err
	}
	println(string(b))
	return nil
}

func Load(path string) (*Circuit, error) {
	return nil, nil
}

func (c *Circuit) Evaluate(input []Bit) []Bit {
	results := make([][]Bit, len(c.parts)+len(c.outputs))
	for i := range results[:len(c.parts)] {
		results[i] = make([]Bit, c.parts[i].Part.NOutputs())
	}
	for i := range results[len(c.parts):] {
		results[i] = make([]Bit, c.outputs[i].Part.NOutputs())
	}

	for i, bit := range input {
		results[i][0] = bit
	}

	for i := c.nInputs; i < len(c.parts); i++ {
		var partInput []Bit
		for _, idxs := range c.parts[i].Inputs {
			partInput = append(partInput, results[idxs[0]][idxs[1]])
		}

		results[i] = c.parts[i].Part.Value(partInput...)
	}

	for i := range c.outputs {
		results[len(c.parts)+i] = results[c.outputs[i].Inputs[0][0]]
	}

	flat := make([]Bit, len(c.outputs))
	for i := range flat {
		flat[i] = results[len(c.parts)+i][0]
	}

	return flat
}

func (c *Circuit) Clone() *Circuit {
	clone := NewCircuit(c.nInputs, c.nOutputs)

	clone.parts = make([]CircuitPartWithConnections, len(c.parts))
	for i, cp := range c.parts {
		clone.parts[i] = cp.Clone()
	}

	clone.outputs = make([]CircuitPartWithConnections, len(c.outputs))
	for i, co := range c.outputs {
		clone.outputs[i] = co.Clone()
	}

	return clone
}

func (c *Circuit) NInputs() int  { return c.nInputs }
func (c *Circuit) NOutputs() int { return c.nOutputs }

// oidx = output idx
// iidx = input idx
// ioidx = input->output idx
func (c *Circuit) SetOutputInput(oidx, iidx, ioidx int) {
	c.outputs[oidx].Inputs = [][2]int{{iidx, ioidx}}
}

func (c *Circuit) AddPart(part CircuitPart, inputs ...[2]int) error {
	for _, idxs := range inputs {
		if idxs[0] >= len(c.parts) {
			return errors.New("circuit: invalid part inputs")
		}
	}

	c.parts = append(c.parts, CircuitPartWithConnections{
		Part:   part,
		Inputs: inputs,
	})

	return nil
}

func (c *Circuit) GetParts() (dst []CircuitPartWithConnections) {
	dst = make([]CircuitPartWithConnections, len(c.parts))
	copy(dst, c.parts)
	return
}

func (c *Circuit) SetParts(parts []CircuitPartWithConnections) {
	c.parts = parts
}

func (c *Circuit) GetOutputs() (dst []CircuitPartWithConnections) {
	dst = make([]CircuitPartWithConnections, len(c.outputs))
	copy(dst, c.outputs)
	return
}

func (c *Circuit) SetOutputs(outputs []CircuitPartWithConnections) {
	c.outputs = outputs
}

func (c *Circuit) NParts() int {
	return len(c.parts)
}

type CircuitPartWithConnections struct {
	Part CircuitPart
	// inputs indexes
	// 0: gate index; must be less than the index of the part itself
	// 1: gate output index
	Inputs [][2]int
}

func (p CircuitPartWithConnections) Clone() CircuitPartWithConnections {
	inputs := make([][2]int, len(p.Inputs))
	copy(inputs, p.Inputs)

	return CircuitPartWithConnections{
		Part:   p.Part,
		Inputs: inputs,
	}
}

type CircuitPart interface {
	Value(...Bit) []Bit
	NInputs() int
	NOutputs() int
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

func (_Buffer) Value(inputs ...Bit) []Bit {
	if len(inputs) != 1 {
		panic("Buffer: must have 1 input")
	}

	return []Bit{inputs[0]}
}

func (_Buffer) NInputs() int  { return 1 }
func (_Buffer) NOutputs() int { return 1 }

func (_OR) Value(inputs ...Bit) []Bit {
	if len(inputs) != 2 {
		panic("OR: must have 2 inputs")
	}

	a, b := inputs[0], inputs[1]

	if a == BitZero && b == BitZero {
		return []Bit{BitZero}
	}

	return []Bit{BitOne}
}

func (_OR) NInputs() int  { return 2 }
func (_OR) NOutputs() int { return 1 }

func (_AND) Value(inputs ...Bit) []Bit {
	if len(inputs) != 2 {
		panic("AND: must have 2 inputs")
	}

	a, b := inputs[0], inputs[1]

	if a == BitOne && b == BitOne {
		return []Bit{BitOne}
	}

	return []Bit{BitZero}
}

func (_AND) NInputs() int  { return 2 }
func (_AND) NOutputs() int { return 1 }

func (_NOT) Value(inputs ...Bit) []Bit {
	if len(inputs) != 1 {
		panic("NOT: must have 1 input")
	}

	a := inputs[0]

	if a == BitOne {
		return []Bit{BitZero}
	}

	return []Bit{BitOne}
}

func (_NOT) NInputs() int  { return 1 }
func (_NOT) NOutputs() int { return 1 }

func (_XOR) Value(inputs ...Bit) []Bit {
	if len(inputs) != 2 {
		panic("XOR: must have 2 inputs")
	}

	a, b := inputs[0], inputs[1]

	if a == BitZero && b == BitZero {
		return []Bit{BitZero}
	}

	if a == BitOne && b == BitOne {
		return []Bit{BitZero}
	}

	return []Bit{BitOne}
}

func (_XOR) NInputs() int  { return 2 }
func (_XOR) NOutputs() int { return 1 }
