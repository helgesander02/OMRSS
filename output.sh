#!/bin/bash
#Online Stream-Aware Routing for TSN-Based Industrial Control Systems: Figure 4

for i in $(seq 1 3); do
    go run main.go  --tsn 3 --avb 7 # Number of input streams: 10
    go run main.go  --tsn 6 --avb 14 # Number of input streams: 20
    go run main.go  --tsn 9 --avb 21 # Number of input streams: 30
    go run main.go  --tsn 12 --avb 28 # Number of input streams: 40
    go run main.go  --tsn 15 --avb 35 # Number of input streams: 50
    go run main.go  --tsn 18 --avb 42 # Number of input streams: 60

    go run main.go --topology_name layered_ring --tsn 3 --avb 7 # Number of input streams: 10
    go run main.go --topology_name layered_ring --tsn 6 --avb 14 # Number of input streams: 20
    go run main.go --topology_name layered_ring --tsn 9 --avb 21 # Number of input streams: 30
    go run main.go --topology_name layered_ring --tsn 12 --avb 28 # Number of input streams: 40
    go run main.go --topology_name layered_ring --tsn 15 --avb 35 # Number of input streams: 50
    go run main.go --topology_name layered_ring --tsn 18 --avb 42 # Number of input streams: 60


    go run main.go --topology_name ring --tsn 3 --avb 7 # Number of input streams: 10
    go run main.go --topology_name ring --tsn 6 --avb 14 # Number of input streams: 20
    go run main.go --topology_name ring --tsn 9 --avb 21 # Number of input streams: 30
    go run main.go --topology_name ring --tsn 12 --avb 28 # Number of input streams: 40
    go run main.go --topology_name ring --tsn 15 --avb 35 # Number of input streams: 50
    go run main.go --topology_name ring --tsn 18 --avb 42 # Number of input streams: 60

done