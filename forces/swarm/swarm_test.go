package swarm

import (
	"ga/forces/models"
	"reflect"
	"testing"
)

func TestCreateVelocity(t *testing.T) {
	p := Particle{
		Pbest:    nil,
		PbestObj: nil,
		Velocity: [3][]float64{},
		Organism: models.Organism{
			DNA: models.DNA{
				{1, 2, 3},
				{4, 5, 6, 7},
				{8, 9, 10, 11},
			},
			Fitness: 0.0,
			Path:    "",
		},
	}

	t.Run("testing create velocity", func(t *testing.T) {

		p.CreateVelocity()

		for i, v := range p.Velocity {
			if len(v) != len(p.Organism.DNA[i]) {
				t.Errorf("wrong dimensions, got %v, wanted %v\n", len(v), len(p.Organism.DNA[i]))
			}
		}
	})

	t.Run("test update Pbest", func(t *testing.T) {
		p.UpdatePBest()

		if !reflect.DeepEqual(p.Pbest, p.DNA) {
			t.Errorf("error updating pBest field")
		}

		if p.Fitness != *p.PbestObj {
			t.Errorf("error updating pBestObj field, wanted %v, got %v\n", p.Fitness, p.PbestObj)
		}

	})

}
