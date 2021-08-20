package main

import (
	"flag"
	"fmt"
)

var (
	MutationRate    = flag.Float64("mut", 0.04, "mutation rate as a decimal")
	PopSize         = flag.Int("pop", 600, "population size")
	NumAtoms        = flag.Int("n", 6, "number of atoms")
	PoolSize        = flag.Float64("pool", 0.50, "fraction size of the the previous generation that survives")
	FitnessLimit    = flag.Float64("f", 10.0, "fitness criteria")
	OutFile         = flag.String("o", "forces.out", "name of output file")
	PathToSpectro   = flag.String("sp", "./spectro", "path/to/spectro")
	PathToSpectroIn = flag.String("i", "./spectro.in", "path/to/spectro.in")
	ZeroChance      = flag.Float64("z", 0.02, "chance mutation will set value to zero instead of adding/subtracting")
)

func init() {
	flag.Parse()
	fmt.Printf("The pop size is %d, the pool size is %f\n", *PopSize, *PoolSize)
}
