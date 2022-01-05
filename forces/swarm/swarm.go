package swarm

import (
	"ga/forces/models"
	"ga/forces/selection"
	"math/rand"
	"time"
)

var r1 *rand.Rand

func init() {
	r := rand.NewSource(time.Now().UnixNano())
	r1 = rand.New(r)
}

// The best DNA ever found, and the best fitness ever found.
var GBest models.DNA
var GbestObj *float64

type Swarm []Particle

type Particle struct {
	// The personal best vector DNA and fitness function.
	Pbest    models.DNA
	PbestObj *float64
	Velocity [3][]float64
	models.Organism
}

func GenerateSwarm() Swarm {
	var s Swarm

	organisms := selection.CreatePopulation()

	for i := 0; i < len(organisms); i++ {
		p := Particle{
			Velocity: [3][]float64{},
			Organism: organisms[i],
		}

		s = append(s, p)
	}

	return s
}

func (p *Particle) CreateVelocity() {
	var v [3][]float64

	// This should allocate the second dimension to be the same length as
	// the DNA... a 1:1 mapping for velocity on the same index.
	// Thus, the corresponding DNA at [i][j] is V[i][j].
	for i, val := range p.DNA {
		v[i] = make([]float64, len(val))
	}

	p.Velocity = v

	for j, k := range p.Velocity {
		for l := range k {
			p.Velocity[j][l] = (r1.NormFloat64() * 0.1)
		}
	}
}

func (p *Particle) UpdatePBest() {
	if p.PbestObj == nil {
		p.Pbest = p.DNA
		p.PbestObj = &p.Fitness
	} else {
		// i.e., if the current personal best is worse than the new
		// one...
		if *p.PbestObj > p.Fitness {
			p.PbestObj = &p.Fitness
			p.Pbest = p.DNA
		}
	}
}
