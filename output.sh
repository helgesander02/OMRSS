#!/bin/bash

for i in $(seq 1 4); do
  go run main.go --topology_name typical_complex --osaco_K 5
  go run main.go --topology_name typical_complex --osaco_K 20
done

#go run main.go --topology_name typical_simple

for i in $(seq 1 4); do
  go run main.go --topology_name layered_ring --osaco_K 5
  go run main.go --topology_name layered_ring --osaco_K 20
done

for i in $(seq 1 4); do
  go run main.go --topology_name ring --osaco_K 5
  go run main.go --topology_name ring --osaco_K 20 
done


