package main

import "flag"

var (
	MutationRate = flag.Float64("mut", 0.0004, "mutation rate as a decimal")
	PopSize      = flag.Int("pop", 600, "population size")
	NumAtoms     = flag.Int("n", 6, "number of atoms")
	PoolSize     = flag.Int("pool", 35, "size of the the pool")
	FitnessLimit = flag.Float64("f", 10.0, "fitness criteria")
)

func init() {
	flag.Parse()
}
