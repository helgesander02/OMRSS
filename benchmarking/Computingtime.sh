#!/bin/bash
#Online Stream-Aware Routing for TSN-Based Industrial Control Systems: Figure 4


go run main.go --topology_name typical_complex --tsn 6 --avb 14 # Number of input streams: 10
go run main.go --topology_name typical_complex --tsn 12 --avb 28 # Number of input streams: 20
go run main.go --topology_name typical_complex --tsn 18 --avb 42 # Number of input streams: 30
go run main.go --topology_name typical_complex --tsn 24 --avb 56 # Number of input streams: 40
go run main.go --topology_name typical_complex --tsn 30 --avb 70 # Number of input streams: 50
go run main.go --topology_name typical_complex --tsn 36 --avb 84 # Number of input streams: 60

go run main.go --topology_name layered_ring --tsn 6 --avb 14 # Number of input streams: 10
go run main.go --topology_name layered_ring --tsn 12 --avb 28 # Number of input streams: 20
go run main.go --topology_name layered_ring --tsn 18 --avb 42 # Number of input streams: 30
go run main.go --topology_name layered_ring --tsn 24 --avb 56 # Number of input streams: 40
go run main.go --topology_name layered_ring --tsn 30 --avb 70 # Number of input streams: 50
go run main.go --topology_name layered_ring --tsn 36 --avb 84 # Number of input streams: 60

go run main.go --topology_name ring --tsn 6 --avb 14 # Number of input streams: 10
go run main.go --topology_name ring --tsn 12 --avb 28 # Number of input streams: 20
go run main.go --topology_name ring --tsn 18 --avb 42 # Number of input streams: 30
go run main.go --topology_name ring --tsn 24 --avb 56 # Number of input streams: 40
go run main.go --topology_name ring --tsn 30 --avb 70 # Number of input streams: 50
go run main.go --topology_name ring --tsn 36 --avb 84 # Number of input streams: 60
