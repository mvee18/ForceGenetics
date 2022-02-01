# GoForceGenetics

GoForceGenetics is a genetic algorithm package meant to generate force constants from harmonic frequencies for a number of future uses.

It follows the methods prescribed in *Practical Genetic Algorithms* by Haupt and Haupt (2004) but changed to fit the specific circumstances.

## Usage and Flags

You can run it with the default parameters using `./forces` though you probably don't want to do that. Instead, the following flags are useful to use:

| Flag | Description | Usage | Default | Note |
|------|-----------|-------------|---------|------|
| **MutationRate** | mutation rate as a decimal | `-mut` | 0.04 | from 0-1 |
| **PopSize** | population size | `-pop` | 600 | int |
| **NumAtoms** | number of atoms | `-n` | 6 | int |
| **PoolSize**: | top fraction of the previous generation that survives | `-pool` | 0.50 | from 0-1 |
| **FitnessLimit** | fitness criteria | `-f` | 1.0 | 0.0 is exact |
| **OutFile** |  name of output file | `-o` | forces.out | string path |
| **PathToSpectro** | path/to/spectro | `-sp` | ./spectro | string path |
| **PathToSpectroIn** | path/to/spectro.in | `-i` | ./spectro.in | string path |
| **ZeroChance** | chance mutation will set value to zero instead of adding/subtracting | `-z` | 0.02 | from 0-1 |
| **TournamentPool** | the number of organisms selected to compete to be parents | `-t` | 3 | should not exceed pool 
| **DerivativeLevel** | level of derivative | `-d` | 2 | 2,3,4 |
| **Fort15File** | path to the fort.15 file that will be used for 3rd derivatives | `-ft2` | ./fort.15 | string path
