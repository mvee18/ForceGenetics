package trad

import (
	"ga/forces/models"
	"reflect"
	"testing"
)

func TestAddMigrant(t *testing.T) {
	t.Run("test adding migrant to pool pointer", func(t *testing.T) {
		o1 := models.Organism{
			DNA:     models.DNA{},
			Fitness: 1.0,
			Path:    "",
		}

		o2 := models.Organism{
			DNA:     models.DNA{},
			Fitness: 10.0,
			Path:    "",
		}

		pool := []models.Organism{o1, o2}

		migOrg := models.Organism{
			DNA:     models.DNA{},
			Fitness: 5.0,
			Path:    "",
		}

		mig := models.Migrant{
			Org:  migOrg,
			Bias: 0.0,
		}

		migrants := make(chan models.Migrant)

		// Add migrant to channel.
		// Note: you cannot just do func(), you need go func().
		go func() { migrants <- mig }()

		models.AddImmigrant(&pool, migrants)

		if len(pool) != 2 {
			t.Errorf("wrong length of pool, wanted %v, got %v\n", 2, len(pool))
		}

		if !reflect.DeepEqual(pool[0], o1) {
			t.Errorf("wrong first pool member, wanted %v, got %v\n", o1, pool[0])
		}

		if !reflect.DeepEqual(pool[1], mig) {
			t.Errorf("wrong first pool member, wanted %v, got %v\n", mig, pool[1])
		}
	})
}
