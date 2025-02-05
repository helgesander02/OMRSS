#!/bin/bash

# FIRST: T(5, 10, 15) and P(0.4, 0.5, 0.6, 0.7, 0.8, 0.9), Next: Cost(set1, set2, set3, set4) and timeout(200, 300, 400)
for i in $(seq 1 20); do
    go run main.go  --topology_name layered_ring --input_tsn 3 --input_avb 7 --test_case 5    # Number of input streams: 10
    go run main.go  --topology_name layered_ring --input_tsn 6 --input_avb 14 --test_case 5   # Number of input streams: 20
    go run main.go  --topology_name layered_ring --input_tsn 9 --input_avb 21 --test_case 5   # Number of input streams: 30
    go run main.go  --topology_name layered_ring --input_tsn 12 --input_avb 28 --test_case 5  # Number of input streams: 40
    go run main.go  --topology_name layered_ring --input_tsn 15 --input_avb 35 --test_case 5  # Number of input streams: 50
    go run main.go  --topology_name layered_ring --input_tsn 18 --input_avb 42 --test_case 5  # Number of input streams: 60
done
