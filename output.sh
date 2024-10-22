#!/bin/bash

go run main.go  --topology_name industrial --input_tsn 3 --input_avb 7  # Number of input streams: 10
go run main.go  --topology_name industrial --input_tsn 6 --input_avb 14  # Number of input streams: 20
go run main.go  --topology_name industrial --input_tsn 9 --input_avb 21  # Number of input streams: 30
go run main.go  --topology_name industrial --input_tsn 12 --input_avb 28  # Number of input streams: 40
go run main.go  --topology_name industrial --input_tsn 15 --input_avb 35 # Number of input streams: 50
go run main.go  --topology_name industrial --input_tsn 18 --input_avb 42  # Number of input streams: 60

go run main.go  --topology_name layered_ring --input_tsn 3 --input_avb 7  # Number of input streams: 10
go run main.go  --topology_name layered_ring --input_tsn 6 --input_avb 14  # Number of input streams: 20
go run main.go  --topology_name layered_ring --input_tsn 9 --input_avb 21  # Number of input streams: 30
go run main.go  --topology_name layered_ring --input_tsn 12 --input_avb 28  # Number of input streams: 40
go run main.go  --topology_name layered_ring --input_tsn 15 --input_avb 35 # Number of input streams: 50
go run main.go  --topology_name layered_ring --input_tsn 18 --input_avb 42  # Number of input streams: 60

go run main.go  --topology_name typical_complex --input_tsn 3 --input_avb 7  # Number of input streams: 10
go run main.go  --topology_name typical_complex --input_tsn 6 --input_avb 14  # Number of input streams: 20
go run main.go  --topology_name typical_complex --input_tsn 9 --input_avb 21  # Number of input streams: 30
go run main.go  --topology_name typical_complex --input_tsn 12 --input_avb 28  # Number of input streams: 40
go run main.go  --topology_name typical_complex --input_tsn 15 --input_avb 35 # Number of input streams: 50
go run main.go  --topology_name typical_complex --input_tsn 18 --input_avb 42  # Number of input streams: 60