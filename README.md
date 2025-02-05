# Online Multicast Routing Simulation System
Simulate TSN and AVB data streams to find the best path in the topology. <br />
<br />
Taking advantage of the distinctions between the Schedulability, AVB Worst Case Delay, and Computation Time statistical algorithms. <br />
<br />
System design from: https://miro.com/welcomeonboard/ZDI1dEJWWGtsTEhCbWMwSm9oTTVhMTk5Y3BaZE83U0hZVDA1S3ZacHJCUmFVaEhhTzhSV0dNTWNMUWU3Mk11YXwzNDU4NzY0NTc0MzMzMTc0Mjg2fDI=?share_link_id=247038481052 <br />

## Installation
* Clone this repo by `git clone https://github.com/helgesander02/OMRSS`
* Env setting `docker build -t 'omrss' .`  `docker run omrss`


## Simulation Settings
We consider a typical topology for TSN-based industrial factories <br />

| typical_complex | typical_simple | 
| --- | --- |
|![alt text](https://github.com/helgesander02/OMRSS/blob/main/img/typical1.png)|![alt text](https://github.com/helgesander02/OMRSS/blob/main/img/typical2.png)|

| layered_ring | ring |
| --- | --- |
|![alt text](https://github.com/helgesander02/OMRSS/blob/main/img/layeredring.png)|![alt text](https://github.com/helgesander02/OMRSS/blob/main/img/ring.png)|

industrial(IN) |  |
| --- | --- |
|![alt text](https://github.com/helgesander02/OMRSS/blob/main/img/industrial.png)||



## Running
Quickstart: `sh output.sh` <br />
<br />
More options:
| Option | Description |
| -------- | ---- | 
| --test_case | Conducting n experiments |
| --topology_name | Topology architecture has typical_complex, typical_simple, ring, layered_ring and industrial |
| --bg_tsn | Number of TSN BG flows |
| --bg_avb | Number of AVB BG flows |
| --input_tsn | Number of TSN Input flows |
| --input_avb | Number of AVB Input flows |
| --HyperPeriod | Greatest Common Divisor of Simulated Time LCM |
| --bandwidth | 1 Gbps |
| --plan_name | The plan comprises OMSACO |
| --osaco_timeout | Timeout in milliseconds |
| --osaco_K | Select K trees with different weights |
| --osaco_P | The pheromone value of each routing path starts to evaporate, where P is the given evaporation coefficient (0 <= p <= 1)|
| --show_network | Present all network information comprehensively |
| --show_plan | Provide a comprehensive display of all plan information |
| --store_data | Store all statistical data |


## Reference
* [Bang Ye Wu, Kun-Mao Chao, "Steiner Minimal Trees"](https://www.csie.ntu.edu.tw/~kmchao/tree10spr/Steiner.pdf)
* [Amal P M, Ajish Kumar K S, "An Algorithm for kth Minimum Spanning Tree"](https://www.sciencedirect.com/science/article/abs/pii/S157106531630083X)
* [Sune Mølgaard Laursen, Paul Pop, Wilfried Steiner, "Routing Optimization of AVB Streams in TSN Networks"](https://backend.orbit.dtu.dk/ws/files/127311642/Sune_Molgaard_Laursen2016aa_Routing_Optimization_of_AVB_St_SIGBED_Review_1.pdf)
* [Ching-Chih Chuang et al., "Online Stream-Aware Routing for TSN-Based Industrial Control Systems"](https://www.researchgate.net/publication/347154804_Online_Stream-Aware_Routing_for_TSN-Based_Industrial_Control_Systems)
* [QINGHAN YU et al., "Online Scheduling for Dynamic VM Migration in Multicast Time-Sensitive Networks"](https://ieeexplore.ieee.org/document/8747398)
* [Jiachen Wen  et al., "Online Updating in Multicast Time-Sensitive Networking"](https://ieeexplore.ieee.org/document/10258186)
* [Xingbo Feng  et al., "Advancing TSN flow scheduling: An efficient framework without flow isolation constraint"](https://www.sciencedirect.com/science/article/pii/S1389128624005206#:~:text=Central%20to%20our%20approach%20is%20a)
* [Mateusz Pawlik, Nikolaus Augsten, "Tree edit distance: Robust and memory-efficient"](https://www.sciencedirect.com/science/article/abs/pii/S0306437915001611)

## TODO
* UnitTesting
* Review Code
