#!/bin/bash

go run main.go --topology_name typical_complex --tsn 6 --avb 14 # Number of input streams: 10
go run main.go --topology_name typical_complex --tsn 12 --avb 28 # Number of input streams: 20
go run main.go --topology_name typical_complex --tsn 18 --avb 42 # Number of input streams: 30
go run main.go --topology_name typical_complex --tsn 24 --avb 56 # Number of input streams: 40
go run main.go --topology_name typical_complex --tsn 30 --avb 70 # Number of input streams: 50
go run main.go --topology_name typical_complex --tsn 36 --avb 84 # Number of input streams: 60
