package islands

import (
	"fmt"
	"ga/forces/informed"
	"ga/forces/models"
	"ga/forces/pseudo"
	trad "ga/forces/traditional"
	"reflect"
)

// Not exactly sure about the methods they share in common.
// Might want to add a TPopulation for the TGA.
type Island interface {
}

// We will gather than on the mig channel, then disperse them across the imm channels for each one.
func RunIslands(immPGA, immTGA, immIGA chan models.Migrant) {
	go pseudo.RunPGA(immPGA)

	go trad.RunTGA(immTGA)

	go informed.RunInformedGA(immIGA)

	migrantPool := make([]models.Migrant, 0)

	for {
		select {
		case p := <-immPGA:
			// fmt.Printf("bias of pga: %v\n", p.Bias)
			migrantPool = append(migrantPool, p)
			MigrationProtocol(p, immPGA, &migrantPool)

		case t := <-immTGA:
			// fmt.Printf("bias of tga: %v\n", t.Bias)
			migrantPool = append(migrantPool, t)
			MigrationProtocol(t, immTGA, &migrantPool)

		case i := <-immIGA:
			// fmt.Printf("bias of iga: %v\n", i.Bias)
			migrantPool = append(migrantPool, i)
			MigrationProtocol(i, immIGA, &migrantPool)
		}
	}
}

/*
func SendBestMigrant(o models.Migrant, mig chan<- models.Migrant, pool []models.Migrant) {
	// First, we need to check if the pool has more than one member.
	if o.Bias >= 0.50 {
		if len(pool) > 1 {
			bestIndex, bestHD := 0, 0.0
			for i, v := range pool {
				hd := models.CalculateHD(o.Org, v.Org)
				if hd > bestHD {
					bestIndex = i
					bestHD = hd
				}
			}

			mig <- models.Migrant{
				Org:  pool[bestIndex].Org,
				Bias: 0.0,
			}

			RemoveIndex(pool, bestIndex)
		}

	} else {
		mig <- models.Migrant{
			Org:  o.Org,
			Bias: 0.0,
		}

	}
}
*/

func MigrationProtocol(o models.Migrant, mig chan<- models.Migrant, pool *[]models.Migrant) {
	// First, we need to check if the pool has more than one member.

	fmt.Println("len of pool in: ", len(*pool))
	for len(*pool) != 0 {
		if len(*pool) > 1 {
			if o.Bias >= 0.50 {
				bestIndex, bestHD := 0, 0.0
				for i, v := range *pool {
					hd := models.CalculateHD(o.Org, v.Org)
					if hd > bestHD {
						bestIndex = i
						bestHD = hd
					}
				}

				mig <- models.Migrant{
					Org:            (*pool)[bestIndex].Org,
					Attractiveness: (*pool)[bestIndex].Attractiveness,
				}

				fmt.Println("sent via bias, pool len: ", len(*pool))
				*pool = RemoveIndex(*pool, bestIndex)

				return

			} else {
				bestAIndex, bestA := 0, 0.0
				for i, v := range *pool {
					if v.Attractiveness > bestA {
						bestAIndex = i
						bestA = v.Attractiveness
					}
				}

				mig <- models.Migrant{
					Org:            (*pool)[bestAIndex].Org,
					Attractiveness: (*pool)[bestAIndex].Attractiveness,
				}

				fmt.Println("sent via attractiveness, pool len: ", len(*pool))
				*pool = RemoveIndex(*pool, bestAIndex)

				return
			}

		} else {
			mig <- models.Migrant{
				Org:            o.Org,
				Attractiveness: o.Attractiveness,
			}

			fmt.Println("sent to self, pool len: ", len(*pool))

			for i, v := range *pool {
				if reflect.DeepEqual(o, v) {
					*pool = RemoveIndex(*pool, i)
					return
				}
			}
		}

	}

}

func RemoveIndex(s []models.Migrant, index int) []models.Migrant {
	return append(s[:index], s[index+1:]...)
}
