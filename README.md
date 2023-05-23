# MST-ACOvsOther
Simulate TSN and AVB data streams to find the best path in the topology. <br />
Taking advantage of the distinctions between the Schedulability, AVB Worst Case Delay, and Computation Time statistical algorithms. <br />

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
[wikipedia, "Steiner tree problem"](https://en.wikipedia.org/wiki/Steiner_tree_problem)<br />
[geeksforgeeks, "Steiner tree problem"](https://www.geeksforgeeks.org/steiner-tree/)<br />
["Steiner tree problem in graphs"](http://sunmoon-template.blogspot.com/2017/04/steiner-tree-problem-in-graphs.html)<br />
["Steiner Minimal Trees"](https://www.csie.ntu.edu.tw/~kmchao/tree10spr/Steiner.pdf)<br />
["Dijkstra’s Algorithm"](https://medium.com/%E6%8A%80%E8%A1%93%E7%AD%86%E8%A8%98/%E5%9F%BA%E7%A4%8E%E6%BC%94%E7%AE%97%E6%B3%95%E7%B3%BB%E5%88%97-graph-%E8%B3%87%E6%96%99%E7%B5%90%E6%A7%8B%E8%88%87dijkstras-algorithm-6134f62c1fc2)<br />


