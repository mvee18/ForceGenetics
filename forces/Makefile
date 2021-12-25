.PHONY: pseudo/pseudo

pseudo/pseudo: pseudo/*.go
	cd pseudo; go test . -v -run TestRunPseudoGA -args -f 1.0 -sp ../spectro -i ../spectro -fi ../forces.inp -pop 800 -pc 1.0 -o ../forces.out; rm -rf forceOrganisms*;