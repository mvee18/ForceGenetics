.PHONY: pseudo/pseudo
.PHONY: informed/informed 
.PHONY: swarm/swarm
.PHONY: models/models 

pseudo/pseudo: pseudo/*.go
	cd pseudo; go test . -v -run TestRunPseudoGA -args -f 1.0 -sp ../spectro -i ../spectro.in -fi ../forces.inp -pop 80 -pc 0.50 -o ../forces.out; rm -rf forceOrganisms*;

informed/informed: informed/*.go
	cd informed; go test . -v -run TestRunIGA -args -f 1.0 -sp ../spectro -i ../spectro.in -fi ../forces.inp -pop 80 -o ../forces.out -imut 0.02; 

swarm/swarm: swarm/*.go
	cd swarm; go test . -v -args -f 1.0 -sp ../spectro -i ../spectro.in -fi ../forces.inp -pop 4 -o ../forces.out; rm -rf forceOrganisms*;

models/models: models/*.go
	cd models; go test -v -bench=. -args -f 1.0 -sp ../spectro -i ../spectro.in -fi ../forces.inp -pop 4 -o ../forces.out; rm -rf forceOrganisms*;