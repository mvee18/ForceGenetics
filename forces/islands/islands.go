package islands

type Island interface {
	GeneratePopulation()
	Migrate()
	Crossover()
}
