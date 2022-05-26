function run_job () {
	cd "$1"
	# cp ../../forces .
	../clear.sh
	# ./forces -ig /home/mvee/research/ForceGenetics/forces/testfiles/h2o/4th -pop 100 -mut "$2" -z 0.0 -dom15 1.0 -dom30 2.0 -dom40 3.0 -ga tga -sp ../spectro -fi ../forces.inp -i ../spectro.in -pool 0.30
}

run_job 000001 0.000001 > 000001.out &
run_job 00001 0.00001 > 00001.out &
run_job 0001 0.0001 > 0001.out &
run_job 001 0.001 > 001.out & 
run_job 005 0.005 > 005.out & 
run_job 010 0.010 > 010.out & 
run_job 015 0.015 > 015.out & 
run_job 020 0.020 > 020.out & 
run_job 050 0.050 > 050.out & 
run_job 100  0.10 > 100.out & 

wait
