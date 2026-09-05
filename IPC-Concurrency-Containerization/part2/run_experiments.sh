#! /bin/bash 

for type in "cpu" "mixed"; do
  for g in 1 2 4 8 16 32 64; do
    for p in 1 2 $(nproc); do
       go run main.go -workload=$type -goroutines=$g -procs=$p >> results/results.txt
    done
  done
done

