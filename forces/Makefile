.PHONY: pseudo/pseudo
.PHONY: informed/informed 
.PHONY: swarm/swarm

pseudo/pseudo: pseudo/*.go
	cd pseudo; go test . -v -run TestRunPseudoGA -args -f 1.0 -sp ../spectro -i ../spectro.in -fi ../forces.inp -pop 80 -pc 0.50 -o ../forces.out; rm -rf forceOrganisms*;

informed/informed: informed/*.go
	cd informed; go test . -v -args -f 1.0 -sp ../spectro -i ../spectro.in -fi ../forces.inp -pop 4 -o ../forces.out; rm -rf forceOrganisms*;

swarm/swarm: swarm/*.go
	cd swarm; go test . -v -args -f 1.0 -sp ../spectro -i ../spectro.in -fi ../forces.inp -pop 4 -o ../forces.out; rm -rf forceOrganisms*;