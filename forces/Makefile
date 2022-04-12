.PHONY: pseudo/pseudo
.PHONY: informed/informed 
.PHONY: swarm/swarm
.PHONY: models/models 
.PHONY: traditional/traditional 
.PHONY: islands/islands 
.PHONY: guess/guess 
.PHONY: quadratic/quadratic 

pseudo/pseudo: pseudo/*.go
	cd pseudo; go test . -v -run TestRunPseudoGA -args -f 1.0 -sp ../spectro -i ../spectro.in -fi ../forces.inp -pop 80 -pc 0.50 -o ../forces.out; rm -rf forceOrganisms*;

informed/informed: informed/*.go
	cd informed; go test . -v -race -run TestRunIGA -args -f 1.0 -sp ../spectro -i ../spectro.in -fi ../forces.inp -pop 80 -o ../forces.out -imut 1.0; 

swarm/swarm: swarm/*.go
	cd swarm; go test . -v -args -f 1.0 -sp ../spectro -i ../spectro.in -fi ../forces.inp -pop 4 -o ../forces.out; rm -rf forceOrganisms*;

models/models: models/*.go
	cd models; go test . -v -args -f 1.0 -sp ../spectro -i ../spectro.in -fi ../forces.inp -pop 4 -o ../forces.out; rm -rf forceOrganisms*;

traditional/traditional: traditional/*.go
	cd traditional; go test . -v -args -f 1.0 -sp ../spectro -i ../spectro.in -fi ../forces.inp -pop 4 -o ../forces.out; rm -rf forceOrganisms*;

islands/islands: islands/*.go
	cd islands; go test . -v -args -f 1.0 -sp ../spectro -i ../spectro.in -fi ../forces.inp -pop 20 -o ../forces.out; rm -rf forceOrganisms*;

guess/guess: guess/*.go
	cd guess; go test . -v -args -f 1.0 -sp ../spectro -i ../spectro.in -fi ../forces.inp -pop 20 -o ../forces.out; rm -rf forceOrganisms*;

quadratic/quadratic: quadratic/*.go
	cd quadratic; go test . -v -run TestLinearInterpolation -args -f 1.0 -sp ../spectro -i ../spectro.in -fi ../forces.inp -pop 20 -o ../forces.out; rm -rf forceOrganisms*;
