package main

import (
	"fmt"
	"math"
	"sync"

	"github.com/ghostec/circuiter/crt"
	"lukechampine.com/frand"
)

// func main() {
// 	circuit := crt.NewCircuit(2, 1)
// 	circuit.AddPart(crt.OR, [2]int{0, 0}, [2]int{1, 0})
// 	circuit.SetOutputInput(0, 2, 0)
// 	circuit.Save("haha")
// }

func main() {
	// samples := []Sample{
	// 	// 2 bit adder
	// 	// a0 a1 b0 b1 ci // s0 s1 c0
	// 	{crt.Bits(0, 0, 0, 0, 0), crt.Bits(0, 0, 0)},
	// 	{crt.Bits(1, 0, 0, 0, 0), crt.Bits(1, 0, 0)},
	// 	{crt.Bits(0, 1, 0, 0, 0), crt.Bits(0, 1, 0)},
	// 	{crt.Bits(1, 0, 0, 0, 0), crt.Bits(0, 1, 0)},
	// 	{crt.Bits(0, 1, 1, 1, 0), crt.Bits(0, 0, 1)},
	// 	{crt.Bits(0, 1, 0, 0, 1), crt.Bits(1, 1, 0)},
	// 	{crt.Bits(1, 0, 0, 0, 1), crt.Bits(1, 1, 0)},
	// 	{crt.Bits(1, 1, 1, 1, 1), crt.Bits(1, 1, 1)},
	// }

	// samples := []Sample{
	// 	{crt.Bits(0, 0, 0), crt.Bits(0)},
	// 	{crt.Bits(0, 0, 1), crt.Bits(0)},
	// 	{crt.Bits(0, 1, 0), crt.Bits(0)},
	// 	{crt.Bits(0, 1, 1), crt.Bits(0)},
	// 	{crt.Bits(1, 0, 0), crt.Bits(1)},
	// 	{crt.Bits(1, 0, 1), crt.Bits(0)},
	// 	{crt.Bits(1, 1, 0), crt.Bits(0)},
	// 	{crt.Bits(1, 1, 1), crt.Bits(0)},
	// }

	// samples := []Sample{
	// 	{crt.Bits(0, 0, 0), crt.Bits(0)},
	// 	{crt.Bits(0, 0, 1), crt.Bits(0)},
	// 	{crt.Bits(0, 1, 0), crt.Bits(0)},
	// 	{crt.Bits(0, 1, 1), crt.Bits(0)},
	// 	{crt.Bits(1, 0, 0), crt.Bits(0)},
	// 	{crt.Bits(1, 0, 1), crt.Bits(1)},
	// 	{crt.Bits(1, 1, 0), crt.Bits(1)},
	// 	{crt.Bits(1, 1, 1), crt.Bits(1)},
	// }

	// samples := SumUint8Samples2Bits()
	samples := SumUint8Samples4Bits()
	PrintSamples(samples)

	factory := &CircuitPartFactory{}
	for _, part := range []crt.CircuitPart{crt.OR, crt.AND, crt.XOR, crt.NOT, crt.Buffer} {
		factory.parts = append(factory.parts, part)
	}

	// TODO: parts with more than one output??
	// how to attach to another later part input

	algo := &MTGeneticAlgorithm{Threads: 20, MaxSurvivors: 103, CircuitPartFactory: factory}
	algo.Execute(100000000, samples, nil)
}

// func NewCircuitPart(samples []Sample)

func SumUint8Samples8Bits() []Sample {
	var samples []Sample

	for i := uint8(0); i < uint8(255); i++ {
		for j := uint8(0); j < uint8(255); j++ {
			samples = append(samples, Sample{
				Input:  append(Uint8ToBits8Bits(i), Uint8ToBits8Bits(j)...),
				Output: Uint8ToBits8Bits(i + j),
			})
		}
	}

	return samples
}

func SumUint8Samples4Bits() []Sample {
	var samples []Sample

	for i := uint8(0); i < uint8(15); i++ {
		for j := uint8(0); j < uint8(15); j++ {
			samples = append(samples, Sample{
				Input:  append(Uint8ToBits4Bits(i), Uint8ToBits4Bits(j)...),
				Output: Uint8ToBits4Bits(i + j),
			})
		}
	}

	return samples
}

func SumUint8Samples2Bits() []Sample {
	var samples []Sample

	for i := uint8(0); i < uint8(3); i++ {
		for j := uint8(0); j < uint8(3); j++ {
			samples = append(samples, Sample{
				Input:  append(Uint8ToBits2Bits(i), Uint8ToBits2Bits(j)...),
				Output: Uint8ToBits2Bits(i + j),
			})
		}
	}

	return samples
}

func Uint8ToBits8Bits(i uint8) (bits []crt.Bit) {
	sbits := fmt.Sprintf("%08b", i)
	sbits = sbits[len(sbits)-8:]

	for _, val := range sbits {
		bit := crt.BitZero
		if val == '1' {
			bit = crt.BitOne
		}
		bits = append(bits, bit)
	}

	return
}

func Uint8ToBits4Bits(i uint8) (bits []crt.Bit) {
	sbits := fmt.Sprintf("%04b", i)
	sbits = sbits[len(sbits)-4:]

	for _, val := range sbits {
		bit := crt.BitZero
		if val == '1' {
			bit = crt.BitOne
		}
		bits = append(bits, bit)
	}

	return
}

func Uint8ToBits2Bits(i uint8) (bits []crt.Bit) {
	sbits := fmt.Sprintf("%02b", i)
	sbits = sbits[len(sbits)-2:]

	for _, val := range sbits {
		bit := crt.BitZero
		if val == '1' {
			bit = crt.BitOne
		}
		bits = append(bits, bit)
	}

	return
}

func PrintSamples(samples []Sample) {
	for _, s := range samples {
		for _, i := range s.Input {
			fmt.Printf("%d", i)
		}
		fmt.Printf("\t")
		for _, o := range s.Output {
			fmt.Printf("%d", o)
		}
		println()
	}
}

type MTGeneticAlgorithm struct {
	Threads            int
	Population         []*Individual
	MaxSurvivors       int
	CircuitPartFactory *CircuitPartFactory
}

func (algo *MTGeneticAlgorithm) Execute(epochs int, samples []Sample, stop <-chan bool) {
	var algos []*GeneticAlgorithm
	for i := 0; i < algo.Threads; i++ {
		algos = append(algos, &GeneticAlgorithm{Population: algo.Population, MaxSurvivors: algo.MaxSurvivors, CircuitPartFactory: algo.CircuitPartFactory})
	}

	wg := sync.WaitGroup{}
	var mostFit *Individual
	var mostFitFitness float64
	var lastFitness float64

	for i := 0; i < epochs; i++ {
		select {
		case <-stop:
			return
		default:
			var population []*Individual

			wg.Add(algo.Threads)
			mu := sync.Mutex{}
			ft := FitnessTracker{}

			for i := 0; i < algo.Threads; i++ {
				subalgo := algos[i]
				go func() {
					result, sft := subalgo.Epoch(algo.Population, samples)

					rMostFit := sft.MostFit()
					rMostFitFitness := sft.Get(rMostFit)

					mu.Lock()
					defer mu.Unlock()

					if rMostFitFitness.Float64() > mostFitFitness {
						mostFit = rMostFit
						mostFitFitness = float64(rMostFitFitness)
					}

					ft.Merge(sft)
					population = append(population, result...)

					wg.Done()
				}()
			}

			wg.Wait()

			ft.Build(population, samples)
			algo.Population = algos[0].SurvivalOfTheFittest(population, ft, samples)

			if mostFitFitness == 1.0 {
				fmt.Printf("%#v\n", mostFit.circuit.NParts())
				return
			}

			if i%500 == 0 || mostFitFitness != lastFitness {
				lastFitness = mostFitFitness
				fmt.Printf("Epoch: %d\tMost fit fitness: %f\n", i, mostFitFitness)
			}
		}
	}
}

type GeneticAlgorithm struct {
	Population         []*Individual
	MaxSurvivors       int
	CircuitPartFactory *CircuitPartFactory
}

func (algo *GeneticAlgorithm) Execute(epochs int, samples []Sample, stop <-chan bool) {
	var ft FitnessTracker
	lastFitness := -1.0

	for i := 0; i < epochs; i++ {
		select {
		case <-stop:
			return
		default:
			algo.Population, ft = algo.Epoch(algo.Population, samples)

			mostFit := ft.MostFit()
			mostFitFitness := ft.Get(mostFit)

			if i%1000 == 0 || mostFitFitness.Float64() != lastFitness {
				lastFitness = mostFitFitness.Float64()
				fmt.Printf("Epoch: %d\tMost fit fitness: %f\n", i, ft.Get(mostFit))
			}
		}
	}
}

func (algo *GeneticAlgorithm) Epoch(original []*Individual, samples []Sample) ([]*Individual, FitnessTracker) {
	population := make([]*Individual, len(original))
	copy(population, original)

	for i := 0; i < algo.MaxSurvivors; i++ {
		population = append(population, RandomIndividual(algo.CircuitPartFactory, len(samples[0].Input), len(samples[0].Output)))
	}

	descendants := make([]*Individual, len(population))
	copy(descendants, population)

	ft := FitnessTracker{}
	for _, individual := range population {
		ft.Set(individual, individual.Fitness(samples))
	}

	for _, individual := range descendants {
		if ft.Get(individual).Float64() < frand.Float64() {
			continue
		}

		var shuffled []*Individual
		copy(shuffled, algo.Population)
		frand.Shuffle(len(shuffled), func(i, j int) { shuffled[i], shuffled[j] = shuffled[j], shuffled[i] })

		for _, other := range shuffled {
			if ft.Get(other).Float64() < frand.Float64() {
				continue
			}

			population = append(population, individual.Mate(other, ft)...)
		}
	}

	for _, individual := range population {
		if frand.Float64() <= 0.5 {
			mutation := individual.Mutate(algo.CircuitPartFactory)
			ft.Set(mutation, mutation.Fitness(samples))
			population = append(population, mutation)
		}
	}

	return algo.SurvivalOfTheFittest(population, ft, samples), ft
}

func (algo *GeneticAlgorithm) SurvivalOfTheFittest(population []*Individual, ft FitnessTracker, samples []Sample) []*Individual {
	mostFit := ft.MostFit()

	survivorsS := IndividualsSet{}
	survivorsS.Add(mostFit)

	for i := 0; i < algo.MaxSurvivors; i++ {
		a, b := population[frand.Intn(len(population))], population[frand.Intn(len(population))]

		if a == b || a == mostFit || b == mostFit {
			continue
		}

		switch {
		case a.circuit.NParts() > 200 && b.circuit.NParts() > 200:
			continue
		case a.circuit.NParts() <= 200 && b.circuit.NParts() > 200:
			survivorsS.Add(a)
			continue
		case a.circuit.NParts() > 200 && b.circuit.NParts() <= 200:
			survivorsS.Add(b)
			continue
		}

		af := ft.Get(a)
		if af == FitnessInvalid {
			af = a.Fitness(samples)
			ft.Set(a, af)
		}

		bf := ft.Get(b)
		if bf == FitnessInvalid {
			bf = b.Fitness(samples)
			ft.Set(b, bf)
		}

		switch {
		case af < bf:
			survivorsS.Add(b)
		case af > bf:
			survivorsS.Add(a)
		default:
			if frand.Intn(2) == 0 {
				survivorsS.Add(b)
			} else {
				survivorsS.Add(a)
			}
		}
	}

	return survivorsS.Slice()
}

type IndividualsSet map[*Individual]bool

func (s IndividualsSet) Add(i *Individual) {
	s[i] = true
}

func (s IndividualsSet) Remove(i *Individual) {
	delete(s, i)
}

func (s IndividualsSet) Contains(i *Individual) bool {
	_, ok := s[i]
	return ok
}

func (s IndividualsSet) Slice() (slice []*Individual) {
	for i := range s {
		slice = append(slice, i)
	}
	return
}

type FitnessTracker map[*Individual]Fitness

func (ft FitnessTracker) Build(population []*Individual, samples []Sample) {
	for _, individual := range population {
		ft[individual] = individual.Fitness(samples)
	}
}

func (ft FitnessTracker) Get(i *Individual) Fitness {
	f, ok := ft[i]
	if !ok {
		return FitnessInvalid
	}
	return f
}

func (ft FitnessTracker) Set(i *Individual, f Fitness) {
	ft[i] = f
}

func (ft FitnessTracker) MostFit() *Individual {
	var mostFit *Individual
	mostFitFitness := Fitness(-1.0)

	for individual, fitness := range ft {
		if fitness.Float64() <= mostFitFitness.Float64() {
			continue
		}

		mostFit = individual
		mostFitFitness = fitness
	}

	return mostFit
}

func (ft FitnessTracker) Merge(other FitnessTracker) {
	for k, v := range other {
		ft[k] = v
	}
}

type Fitness float64

const (
	FitnessInvalid = Fitness(-1.0)
)

func (f Fitness) Float64() float64 {
	return float64(f)
}

type Individual struct {
	circuit *crt.Circuit
}

func (individual *Individual) Mate(other *Individual, ft FitnessTracker) (children []*Individual) {
	mate := func() *Individual {
		a, b := individual.Clone(), other.Clone()

		var aOutputs, bOutputs []int

		anf, bnf := ft.Get(individual).Float64(), ft.Get(other).Float64()
		anf, bnf = anf/(anf+bnf), bnf/(anf+bnf)

		for i := 0; i < individual.circuit.NOutputs(); i++ {
			if anf > frand.Float64() {
				aOutputs = append(aOutputs, i)
				continue
			}
			bOutputs = append(bOutputs, i)
		}

		buildForOutputs := func(outputsIdxs []int, circuit, child *crt.Circuit) {
			var tmpset IntSet
			var tmpstack IntStack

			// build paths that lead to all outputsIdxs
			outputs := circuit.GetOutputs()
			for _, oidx := range outputsIdxs {
				for _, itupl := range outputs[oidx].Inputs {
					if tmpset.Contains(itupl[0]) {
						continue
					}
					tmpset.Add(itupl[0])
					tmpstack.Add(itupl[0])
				}
			}

			childParts := child.GetParts()

			circuitParts := circuit.GetParts()

			for !tmpstack.Empty() {
				partIdx := tmpstack.Pop()
				part := circuitParts[partIdx].Part
				inputs := make([][2]int, part.NInputs())

				for i := 0; i < part.NInputs(); i++ {
					itupl := circuitParts[partIdx].Inputs[i]

					if itupl[0] < circuit.NInputs() {
						inputs[i] = itupl
						continue
					}

					inputs[i][0] = len(childParts) - part.NInputs() + i // is the order correct? or should be reverse?
					inputs[i][1] = itupl[1]
				}

				cpwc := crt.CircuitPartWithConnections{Part: part, Inputs: inputs}
				childParts = append(childParts, cpwc)
				child.AddPart(part, inputs...)
			}

			for _, oidx := range outputsIdxs {
				inputIdx := frand.Intn(len(childParts))
				outputIdx := frand.Intn(childParts[inputIdx].Part.NOutputs())
				child.SetOutputInput(oidx, inputIdx, outputIdx)
			}
		}

		child := crt.NewCircuit(a.circuit.NInputs(), a.circuit.NOutputs())
		buildForOutputs(aOutputs, a.circuit, child)
		buildForOutputs(bOutputs, b.circuit, child)

		return &Individual{circuit: child}
	}

	for i := 0; i < 16; i++ {
		children = append(children, mate())
	}

	return
}

func (individual *Individual) Mutate(factory *CircuitPartFactory) *Individual {
	clone := individual.Clone()

	for {
		chance := frand.Float64()
		switch {
		case chance < 0.33:
			var idx int
			for {
				idx = frand.Intn(clone.circuit.NParts())
				if idx >= clone.circuit.NInputs() {
					break
				}
			}

			parts := clone.circuit.GetParts()

			for {
				random := factory.RandomCircuitPart()

				if random == parts[idx].Part {
					continue
				}

				switch {
				case random.NInputs() == parts[idx].Part.NInputs():
					parts[idx].Part = random
				case random.NInputs() != parts[idx].Part.NInputs():
					var inputs [][2]int

					for len(inputs) < random.NInputs() {
						inputIdx := frand.Intn(idx)
						outputIdx := frand.Intn(parts[inputIdx].Part.NOutputs())
						inputs = append(inputs, [2]int{inputIdx, outputIdx})
					}

					parts[idx].Part = random
					parts[idx].Inputs = inputs
				}

				break
			}
			clone.circuit.SetParts(parts)
		case chance >= 0.33 && chance < 0.67:
			var idx int
			for {
				idx = frand.Intn(clone.circuit.NParts())
				if idx >= clone.circuit.NInputs() {
					break
				}
			}

			parts := clone.circuit.GetParts()
			input := frand.Intn(len(parts[idx].Inputs))
			inputIdx := frand.Intn(idx)
			outputIdx := frand.Intn(parts[inputIdx].Part.NOutputs())
			parts[idx].Inputs[input] = [2]int{inputIdx, outputIdx}
			clone.circuit.SetParts(parts)
		default:
			outputs := clone.circuit.GetOutputs()
			idx := frand.Intn(len(outputs))
			inputIdx := frand.Intn(clone.circuit.NParts())
			outputIdx := frand.Intn(clone.circuit.GetParts()[inputIdx].Part.NOutputs())
			outputs[idx].Inputs[0] = [2]int{inputIdx, outputIdx}
			clone.circuit.SetOutputs(outputs)
		}

		if frand.Float64() > 0.25 {
			break
		}
	}

	return clone
}

func (individual *Individual) Clone() *Individual {
	return &Individual{circuit: individual.circuit.Clone()}
}

func (individual *Individual) Fitness(samples []Sample) Fitness {
	if len(samples) == 0 {
		return Fitness(-1.0)
	}

	var fitnessRegular, fitnessCorrectBlocks float64

	for _, sample := range samples {
		fitnessRegular += individual.fitnessRegular(sample)
		fitnessCorrectBlocks += individual.fitnessCorrectBlocks(sample)
	}

	fitnessRegular /= float64(len(samples) * len(samples[0].Output))
	fitnessCorrectBlocks /= float64(len(samples) * len(samples[0].Output) * len(samples[0].Output))

	return Fitness(math.Sqrt(fitnessRegular * fitnessCorrectBlocks))
}

func (individual *Individual) fitnessRegular(sample Sample) float64 {
	output := individual.circuit.Evaluate(sample.Input)

	var fitness float64
	var correct int

	for i := range sample.Output {
		if sample.Output[i] != output[i] {
			fitness += float64(correct)
			correct = 0
			continue
		}
		correct += 1
	}

	fitness += float64(correct)

	return fitness
}

func (individual *Individual) fitnessCorrectBlocks(sample Sample) float64 {
	output := individual.circuit.Evaluate(sample.Input)

	var fitness float64
	var correct int

	for i := range sample.Output {
		if sample.Output[i] != output[i] {
			fitness += float64(correct * correct)
			correct = 0
			continue
		}
		correct += 1
	}

	fitness += float64(correct * correct)

	return fitness
}

type Sample struct {
	Input  []crt.Bit
	Output []crt.Bit
}

func RandomIndividual(factory *CircuitPartFactory, nInputs, nOutputs int) *Individual {
	newPartChance := 0.25

	circuit := crt.NewCircuit(nInputs, nOutputs)
	var parts []crt.CircuitPart

	for _, part := range circuit.GetParts() {
		parts = append(parts, part.Part)
	}

	for {
		if frand.Float64() > newPartChance && circuit.NParts() > nInputs {
			break
		}

		part := factory.RandomCircuitPart()
		var partInputs [][2]int

		for i := 0; i < part.NInputs(); i++ {
			inputIdx := frand.Intn(circuit.NParts())
			outputIdx := frand.Intn(parts[inputIdx].NOutputs())
			partInputs = append(partInputs, [2]int{inputIdx, outputIdx})
		}

		circuit.AddPart(part, partInputs...)
		parts = append(parts, part)
	}

	for i := 0; i < nOutputs; i++ {
		inputIdx := frand.Intn(circuit.NParts())
		outputIdx := frand.Intn(parts[inputIdx].NOutputs())
		circuit.SetOutputInput(i, inputIdx, outputIdx)
	}

	return &Individual{circuit: circuit}
}

type CircuitPartFactory struct {
	parts []crt.CircuitPart
}

func (f *CircuitPartFactory) RandomCircuitPart() crt.CircuitPart {
	return f.parts[frand.Intn(len(f.parts))]
}
