# MST-ACOvsOther
Simulate TSN and AVB data streams to find the best path in the topology. <br />
<br />
Taking advantage of the distinctions between the Schedulability, AVB Worst Case Delay, and Computation Time statistical algorithms. <br />

## Installation
Clone this repo by `git clone https://github.com/helgesander02/MST-ACOvsOther`<br />

## Running
Quickstart: `cd MST-ACOvsOther` `go run main.go` <br />
<br />
More options:
| Option | Description |
| -------- | ---- | 
| --test_case | Conducting n experiments |
| --tsn | Number of TSN flows |
| --avb | Number of AVB flows |
| --HyperPeriod | Greatest Common Divisor of Simulated Time LCM |
| --bandwidth | 1 Gbps ==> bytes/hyperperiod  |
| --K | finds kth minimum spanning tree |
| --show_topology | Display all topology information |
| --show_flows | Display all flows information |
| --show_routes | Display all routes information |


## Reference
[Sean Chou, "Dijkstra’s Algorithm"](https://medium.com/%E6%8A%80%E8%A1%93%E7%AD%86%E8%A8%98/%E5%9F%BA%E7%A4%8E%E6%BC%94%E7%AE%97%E6%B3%95%E7%B3%BB%E5%88%97-graph-%E8%B3%87%E6%96%99%E7%B5%90%E6%A7%8B%E8%88%87dijkstras-algorithm-6134f62c1fc2)<br />
[Bang Ye Wu, Kun-Mao Chao, "Steiner Minimal Trees"](https://www.csie.ntu.edu.tw/~kmchao/tree10spr/Steiner.pdf)<br />
[Amal P M, Ajish Kumar K S, "An Algorithm for kth Minimum Spanning Tree"](https://www.sciencedirect.com/science/article/abs/pii/S157106531630083X)<br />
[Ching-Chih Chuang et al., "Online Stream-Aware Routing for TSN-Based Industrial Control Systems"](https://www.researchgate.net/publication/347154804_Online_Stream-Aware_Routing_for_TSN-Based_Industrial_Control_Systems)<br />
[Sune Mølgaard Laursen, Paul Pop, Wilfried Steiner, "Routing Optimization of AVB Streams in TSN Networks"](https://backend.orbit.dtu.dk/ws/files/127311642/Sune_Molgaard_Laursen2016aa_Routing_Optimization_of_AVB_St_SIGBED_Review_1.pdf)<br />



