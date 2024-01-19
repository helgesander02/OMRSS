#!/bin/bash

for i in $(seq 1 100); do
  go run main.go --topology_name typical_complex --osaco_K 5 --osaco_P 0.6
done

#go run main.go --topology_name typical_simple

#go run main.go --topology_name layered_ring --osaco_K 5 --osaco_P 0.8

#go run main.go --topology_name ring --osaco_K 5 --osaco_P 0.8


