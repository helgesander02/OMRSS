#!/bin/bash

go run main.go --tsn 70 --avb 30 --topology_name layered_ring
go run main.go --tsn 70 --avb 40 --topology_name layered_ring
go run main.go --tsn 70 --avb 50 --topology_name layered_ring
go run main.go --tsn 80 --avb 30 --topology_name layered_ring
go run main.go --tsn 80 --avb 40 --topology_name layered_ring
go run main.go --tsn 80 --avb 50 --topology_name layered_ring
go run main.go --tsn 90 --avb 30 --topology_name layered_ring
go run main.go --tsn 90 --avb 40 --topology_name layered_ring
go run main.go --tsn 90 --avb 50 --topology_name layered_ring
go run main.go --tsn 100 --avb 30 --topology_name layered_ring
go run main.go --tsn 100 --avb 40 --topology_name layered_ring
go run main.go --tsn 100 --avb 50 --topology_name layered_ring

go run main.go --tsn 70 --avb 30 --topology_name ring
go run main.go --tsn 70 --avb 40 --topology_name ring 
go run main.go --tsn 70 --avb 50 --topology_name ring
go run main.go --tsn 80 --avb 30 --topology_name ring
go run main.go --tsn 80 --avb 40 --topology_name ring 
go run main.go --tsn 80 --avb 50 --topology_name ring
go run main.go --tsn 90 --avb 30 --topology_name ring
go run main.go --tsn 90 --avb 40 --topology_name ring
go run main.go --tsn 90 --avb 50 --topology_name ring
go run main.go --tsn 100 --avb 30 --topology_name ring
go run main.go --tsn 100 --avb 40 --topology_name ring
go run main.go --tsn 100 --avb 50 --topology_name ring