# MST-ACOvsOther
Simulate TSN and AVB data streams to find the best path in the topology.

## Installation
Clone this repo by `git clone https://github.com/helgesander02/MST-ACOvsOther`<br />

## Running
Quickstart: `go run main.go`<br />
More options:
| Option | Description |
| -------- | ---- | 
| --test_case | Conducting n experiments |
| --tsn | Number of TSN flows |
| --avb | Number of AVB flows |
| --HyperPeriod | Greatest Common Divisor of Simulated Time LCM |
| --show_topology | Display all topology information |
| --show_flows | Display all Flows information |


## Reference
[geeksforgeeks, "Steiner Tree Problem"](https://www.geeksforgeeks.org/steiner-tree/)<br />
[知乎, "viterbi"](https://www.zhihu.com/question/20136144)<br />
[Cisco, "OSPF"](https://www.cisco.com/c/zh_tw/support/docs/ip/open-shortest-path-first-ospf/7039-1.html)
[geeksforgeeks, "Dijkstra"](https://www.geeksforgeeks.org/dijkstras-shortest-path-algorithm-greedy-algo-7/)

