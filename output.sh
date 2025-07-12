#!/bin/bash

for i in $(seq 1 20); do
    go run main.go  --topology_name layered_ring --input_tsn 3 --input_avb 7 --test_case 5 --osaco_timeout 200    # Number of input streams: 10
    go run main.go  --topology_name layered_ring --input_tsn 6 --input_avb 14 --test_case 5 --osaco_timeout 200   # Number of input streams: 20
    go run main.go  --topology_name layered_ring --input_tsn 9 --input_avb 21 --test_case 5 --osaco_timeout 200   # Number of input streams: 30
    go run main.go  --topology_name layered_ring --input_tsn 12 --input_avb 28 --test_case 5 --osaco_timeout 200  # Number of input streams: 40
    go run main.go  --topology_name layered_ring --input_tsn 15 --input_avb 35 --test_case 5 --osaco_timeout 200  # Number of input streams: 50
    go run main.go  --topology_name layered_ring --input_tsn 18 --input_avb 42 --test_case 5 --osaco_timeout 200  # Number of input streams: 60
    go run main.go  --topology_name layered_ring --input_tsn 21 --input_avb 49 --test_case 5 --osaco_timeout 200  # Number of input streams: 70
    go run main.go  --topology_name layered_ring --input_tsn 24 --input_avb 56 --test_case 5 --osaco_timeout 200  # Number of input streams: 80
    go run main.go  --topology_name layered_ring --input_tsn 27 --input_avb 63 --test_case 5 --osaco_timeout 200  # Number of input streams: 90
    go run main.go  --topology_name layered_ring --input_tsn 30 --input_avb 70 --test_case 5 --osaco_timeout 200  # Number of input streams: 100
done

for i in $(seq 1 20); do
    go run main.go  --topology_name typical_complex --input_tsn 3 --input_avb 7 --test_case 5 --osaco_timeout 200    # Number of input streams: 10
    go run main.go  --topology_name typical_complex --input_tsn 6 --input_avb 14 --test_case 5 --osaco_timeout 200   # Number of input streams: 20
    go run main.go  --topology_name typical_complex --input_tsn 9 --input_avb 21 --test_case 5 --osaco_timeout 200   # Number of input streams: 30
    go run main.go  --topology_name typical_complex --input_tsn 12 --input_avb 28 --test_case 5 --osaco_timeout 200  # Number of input streams: 40
    go run main.go  --topology_name typical_complex --input_tsn 15 --input_avb 35 --test_case 5 --osaco_timeout 200  # Number of input streams: 50
    go run main.go  --topology_name typical_complex --input_tsn 18 --input_avb 42 --test_case 5 --osaco_timeout 200  # Number of input streams: 60
    go run main.go  --topology_name typical_complex --input_tsn 21 --input_avb 49 --test_case 5 --osaco_timeout 200  # Number of input streams: 70
    go run main.go  --topology_name typical_complex --input_tsn 24 --input_avb 56 --test_case 5 --osaco_timeout 200  # Number of input streams: 80
    go run main.go  --topology_name typical_complex --input_tsn 27 --input_avb 63 --test_case 5 --osaco_timeout 200  # Number of input streams: 90
    go run main.go  --topology_name typical_complex --input_tsn 30 --input_avb 70 --test_case 5 --osaco_timeout 200  # Number of input streams: 100
done

for i in $(seq 1 20); do
    go run main.go  --topology_name industrial --input_tsn 3 --input_avb 7 --test_case 5 --osaco_timeout 200    # Number of input streams: 10
    go run main.go  --topology_name industrial --input_tsn 6 --input_avb 14 --test_case 5 --osaco_timeout 200   # Number of input streams: 20
    go run main.go  --topology_name industrial --input_tsn 9 --input_avb 21 --test_case 5 --osaco_timeout 200   # Number of input streams: 30
    go run main.go  --topology_name industrial --input_tsn 12 --input_avb 28 --test_case 5 --osaco_timeout 200  # Number of input streams: 40
    go run main.go  --topology_name industrial --input_tsn 15 --input_avb 35 --test_case 5 --osaco_timeout 200  # Number of input streams: 50
    go run main.go  --topology_name industrial --input_tsn 18 --input_avb 42 --test_case 5 --osaco_timeout 200  # Number of input streams: 60
    go run main.go  --topology_name industrial --input_tsn 21 --input_avb 49 --test_case 5 --osaco_timeout 200  # Number of input streams: 70
    go run main.go  --topology_name industrial --input_tsn 24 --input_avb 56 --test_case 5 --osaco_timeout 200  # Number of input streams: 80
    go run main.go  --topology_name industrial --input_tsn 27 --input_avb 63 --test_case 5 --osaco_timeout 200  # Number of input streams: 90
    go run main.go  --topology_name industrial --input_tsn 30 --input_avb 70 --test_case 5 --osaco_timeout 200  # Number of input streams: 100
done

for i in $(seq 1 20); do
    go run main.go  --topology_name ring --input_tsn 3 --input_avb 7 --test_case 5 --osaco_timeout 200    # Number of input streams: 10
    go run main.go  --topology_name ring --input_tsn 6 --input_avb 14 --test_case 5 --osaco_timeout 200   # Number of input streams: 20
    go run main.go  --topology_name ring --input_tsn 9 --input_avb 21 --test_case 5 --osaco_timeout 200   # Number of input streams: 30
    go run main.go  --topology_name ring --input_tsn 12 --input_avb 28 --test_case 5 --osaco_timeout 200  # Number of input streams: 40
    go run main.go  --topology_name ring --input_tsn 15 --input_avb 35 --test_case 5 --osaco_timeout 200  # Number of input streams: 50
    go run main.go  --topology_name ring --input_tsn 18 --input_avb 42 --test_case 5 --osaco_timeout 200  # Number of input streams: 60
    go run main.go  --topology_name ring --input_tsn 21 --input_avb 49 --test_case 5 --osaco_timeout 200  # Number of input streams: 70
    go run main.go  --topology_name ring --input_tsn 24 --input_avb 56 --test_case 5 --osaco_timeout 200  # Number of input streams: 80
    go run main.go  --topology_name ring --input_tsn 27 --input_avb 63 --test_case 5 --osaco_timeout 200  # Number of input streams: 90
    go run main.go  --topology_name ring --input_tsn 30 --input_avb 70 --test_case 5 --osaco_timeout 200  # Number of input streams: 100
done
