.PHONY: pseudo/pseudo
.PHONY: informed/informed 

pseudo/pseudo: pseudo/*.go
	cd pseudo; go test . -v -run TestRunPseudoGA -args -f 1.0 -sp ../spectro -i ../spectro -fi ../forces.inp -pop 800 -pc 1.0 -o ../forces.out; rm -rf forceOrganisms*;

informed/informed: informed/*.go
	cd informed; go test . -v -args -f 1.0 -sp ../spectro -i ../spectro -fi ../forces.inp -pop 4 -o ../forces.out; rm -rf forceOrganisms*;
