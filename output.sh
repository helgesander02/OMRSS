#!/bin/bash

for i in $(seq 1 4); do
  # K=5
  go run main.go --topology_name ring --input_tsn 3 --input_avb 7 --osaco_K 5 --osaco_P 0.8 # Number of input streams: 10
  go run main.go --topology_name ring --input_tsn 6 --input_avb 14 --osaco_K 5 --osaco_P 0.8 # Number of input streams: 20
  go run main.go --topology_name ring --input_tsn 9 --input_avb 21 --osaco_K 5 --osaco_P 0.8 # Number of input streams: 30
  go run main.go --topology_name ring --input_tsn 12 --input_avb 28 --osaco_K 5 --osaco_P 0.8 # Number of input streams: 40
  go run main.go --topology_name ring --input_tsn 15 --input_avb 35 --osaco_K 5 --osaco_P 0.8 # Number of input streams: 50
  go run main.go --topology_name ring --input_tsn 18 --input_avb 42 --osaco_K 5 --osaco_P 0.8 # Number of input streams: 60

  go run main.go --topology_name ring --input_tsn 3 --input_avb 7 --osaco_K 5 --osaco_P 0.7 # Number of input streams: 10
  go run main.go --topology_name ring --input_tsn 6 --input_avb 14 --osaco_K 5 --osaco_P 0.7 # Number of input streams: 20
  go run main.go --topology_name ring --input_tsn 9 --input_avb 21 --osaco_K 5 --osaco_P 0.7 # Number of input streams: 30
  go run main.go --topology_name ring --input_tsn 12 --input_avb 28 --osaco_K 5 --osaco_P 0.7 # Number of input streams: 40
  go run main.go --topology_name ring --input_tsn 15 --input_avb 35 --osaco_K 5 --osaco_P 0.7 # Number of input streams: 50
  go run main.go --topology_name ring --input_tsn 18 --input_avb 42 --osaco_K 5 --osaco_P 0.7 # Number of input streams: 60

  go run main.go --topology_name ring --input_tsn 3 --input_avb 7 --osaco_K 5 --osaco_P 0.6 # Number of input streams: 10
  go run main.go --topology_name ring --input_tsn 6 --input_avb 14 --osaco_K 5 --osaco_P 0.6 # Number of input streams: 20
  go run main.go --topology_name ring --input_tsn 9 --input_avb 21 --osaco_K 5 --osaco_P 0.6 # Number of input streams: 30
  go run main.go --topology_name ring --input_tsn 12 --input_avb 28 --osaco_K 5 --osaco_P 0.6 # Number of input streams: 40
  go run main.go --topology_name ring --input_tsn 15 --input_avb 35 --osaco_K 5 --osaco_P 0.6 # Number of input streams: 50
  go run main.go --topology_name ring --input_tsn 18 --input_avb 42 --osaco_K 5 --osaco_P 0.6 # Number of input streams: 60

  go run main.go --topology_name ring --input_tsn 3 --input_avb 7 --osaco_K 5 --osaco_P 0.5 # Number of input streams: 10
  go run main.go --topology_name ring --input_tsn 6 --input_avb 14 --osaco_K 5 --osaco_P 0.5 # Number of input streams: 20
  go run main.go --topology_name ring --input_tsn 9 --input_avb 21 --osaco_K 5 --osaco_P 0.5 # Number of input streams: 30
  go run main.go --topology_name ring --input_tsn 12 --input_avb 28 --osaco_K 5 --osaco_P 0.5 # Number of input streams: 40
  go run main.go --topology_name ring --input_tsn 15 --input_avb 35 --osaco_K 5 --osaco_P 0.5 # Number of input streams: 50
  go run main.go --topology_name ring --input_tsn 18 --input_avb 42 --osaco_K 5 --osaco_P 0.5 # Number of input streams: 60

  # K=10
  go run main.go --topology_name ring --input_tsn 3 --input_avb 7 --osaco_K 10 --osaco_P 0.8 # Number of input streams: 10
  go run main.go --topology_name ring --input_tsn 6 --input_avb 14 --osaco_K 10 --osaco_P 0.8 # Number of input streams: 20
  go run main.go --topology_name ring --input_tsn 9 --input_avb 21 --osaco_K 10 --osaco_P 0.8 # Number of input streams: 30
  go run main.go --topology_name ring --input_tsn 12 --input_avb 28 --osaco_K 10 --osaco_P 0.8 # Number of input streams: 40
  go run main.go --topology_name ring --input_tsn 15 --input_avb 35 --osaco_K 10 --osaco_P 0.8 # Number of input streams: 50
  go run main.go --topology_name ring --input_tsn 18 --input_avb 42 --osaco_K 10 --osaco_P 0.8 # Number of input streams: 60

  go run main.go --topology_name ring --input_tsn 3 --input_avb 7 --osaco_K 10 --osaco_P 0.7 # Number of input streams: 10
  go run main.go --topology_name ring --input_tsn 6 --input_avb 14 --osaco_K 10 --osaco_P 0.7 # Number of input streams: 20
  go run main.go --topology_name ring --input_tsn 9 --input_avb 21 --osaco_K 10 --osaco_P 0.7 # Number of input streams: 30
  go run main.go --topology_name ring --input_tsn 12 --input_avb 28 --osaco_K 10 --osaco_P 0.7 # Number of input streams: 40
  go run main.go --topology_name ring --input_tsn 15 --input_avb 35 --osaco_K 10 --osaco_P 0.7 # Number of input streams: 50
  go run main.go --topology_name ring --input_tsn 18 --input_avb 42 --osaco_K 10 --osaco_P 0.7 # Number of input streams: 60

  go run main.go --topology_name ring --input_tsn 3 --input_avb 7 --osaco_K 10 --osaco_P 0.6 # Number of input streams: 10
  go run main.go --topology_name ring --input_tsn 6 --input_avb 14 --osaco_K 10 --osaco_P 0.6 # Number of input streams: 20
  go run main.go --topology_name ring --input_tsn 9 --input_avb 21 --osaco_K 10 --osaco_P 0.6 # Number of input streams: 30
  go run main.go --topology_name ring --input_tsn 12 --input_avb 28 --osaco_K 10 --osaco_P 0.6 # Number of input streams: 40
  go run main.go --topology_name ring --input_tsn 15 --input_avb 35 --osaco_K 10 --osaco_P 0.6 # Number of input streams: 50
  go run main.go --topology_name ring --input_tsn 18 --input_avb 42 --osaco_K 10 --osaco_P 0.6 # Number of input streams: 60

  go run main.go --topology_name ring --input_tsn 3 --input_avb 7 --osaco_K 10 --osaco_P 0.5 # Number of input streams: 10
  go run main.go --topology_name ring --input_tsn 6 --input_avb 14 --osaco_K 10 --osaco_P 0.5 # Number of input streams: 20
  go run main.go --topology_name ring --input_tsn 9 --input_avb 21 --osaco_K 10 --osaco_P 0.5 # Number of input streams: 30
  go run main.go --topology_name ring --input_tsn 12 --input_avb 28 --osaco_K 10 --osaco_P 0.5 # Number of input streams: 40
  go run main.go --topology_name ring --input_tsn 15 --input_avb 35 --osaco_K 10 --osaco_P 0.5 # Number of input streams: 50
  go run main.go --topology_name ring --input_tsn 18 --input_avb 42 --osaco_K 10 --osaco_P 0.5 # Number of input streams: 60

  # K=15
  go run main.go --topology_name ring --input_tsn 3 --input_avb 7 --osaco_K 15 --osaco_P 0.8 # Number of input streams: 10
  go run main.go --topology_name ring --input_tsn 6 --input_avb 14 --osaco_K 15 --osaco_P 0.8 # Number of input streams: 20
  go run main.go --topology_name ring --input_tsn 9 --input_avb 21 --osaco_K 15 --osaco_P 0.8 # Number of input streams: 30
  go run main.go --topology_name ring --input_tsn 12 --input_avb 28 --osaco_K 15 --osaco_P 0.8 # Number of input streams: 40
  go run main.go --topology_name ring --input_tsn 15 --input_avb 35 --osaco_K 15 --osaco_P 0.8 # Number of input streams: 50
  go run main.go --topology_name ring --input_tsn 18 --input_avb 42 --osaco_K 15 --osaco_P 0.8 # Number of input streams: 60

  go run main.go --topology_name ring --input_tsn 3 --input_avb 7 --osaco_K 15 --osaco_P 0.7 # Number of input streams: 10
  go run main.go --topology_name ring --input_tsn 6 --input_avb 14 --osaco_K 15 --osaco_P 0.7 # Number of input streams: 20
  go run main.go --topology_name ring --input_tsn 9 --input_avb 21 --osaco_K 15 --osaco_P 0.7 # Number of input streams: 30
  go run main.go --topology_name ring --input_tsn 12 --input_avb 28 --osaco_K 15 --osaco_P 0.7 # Number of input streams: 40
  go run main.go --topology_name ring --input_tsn 15 --input_avb 35 --osaco_K 15 --osaco_P 0.7 # Number of input streams: 50
  go run main.go --topology_name ring --input_tsn 18 --input_avb 42 --osaco_K 15 --osaco_P 0.7 # Number of input streams: 60

  go run main.go --topology_name ring --input_tsn 3 --input_avb 7 --osaco_K 15 --osaco_P 0.6 # Number of input streams: 10
  go run main.go --topology_name ring --input_tsn 6 --input_avb 14 --osaco_K 15 --osaco_P 0.6 # Number of input streams: 20
  go run main.go --topology_name ring --input_tsn 9 --input_avb 21 --osaco_K 15 --osaco_P 0.6 # Number of input streams: 30
  go run main.go --topology_name ring --input_tsn 12 --input_avb 28 --osaco_K 15 --osaco_P 0.6 # Number of input streams: 40
  go run main.go --topology_name ring --input_tsn 15 --input_avb 35 --osaco_K 15 --osaco_P 0.6 # Number of input streams: 50
  go run main.go --topology_name ring --input_tsn 18 --input_avb 42 --osaco_K 15 --osaco_P 0.6 # Number of input streams: 60

  go run main.go --topology_name ring --input_tsn 3 --input_avb 7 --osaco_K 15 --osaco_P 0.5 # Number of input streams: 10
  go run main.go --topology_name ring --input_tsn 6 --input_avb 14 --osaco_K 15 --osaco_P 0.5 # Number of input streams: 20
  go run main.go --topology_name ring --input_tsn 9 --input_avb 21 --osaco_K 15 --osaco_P 0.5 # Number of input streams: 30
  go run main.go --topology_name ring --input_tsn 12 --input_avb 28 --osaco_K 15 --osaco_P 0.5 # Number of input streams: 40
  go run main.go --topology_name ring --input_tsn 15 --input_avb 35 --osaco_K 15 --osaco_P 0.5 # Number of input streams: 50
  go run main.go --topology_name ring --input_tsn 18 --input_avb 42 --osaco_K 15 --osaco_P 0.5 # Number of input streams: 60

done


