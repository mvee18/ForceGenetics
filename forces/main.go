package main

import (
	"ga/forces/flags"
	"ga/forces/islands"
	"ga/forces/models"
	trad "ga/forces/traditional"
)

var (
	immPGA = make(chan models.Migrant)
	immTGA = make(chan models.Migrant)
	immIGA = make(chan models.Migrant)
)

func main() {
	switch *flags.GAs {
	case "all":
		islands.RunIslands(immPGA, immTGA, immIGA)
	case "tga":
		trad.RunTGAOnly()
	}

}
