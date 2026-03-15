package algo

import (
	"src/plan/algo_timer"
	"src/plan/routes"
)

type SMT struct {
	Trees      *routes.TreesSet
	InputTrees *routes.TreesSet
	BGTrees    *routes.TreesSet
	Objs_smt   [4]float64
	Timer      *algo_timer.Timer
}

type MDTC struct {
	Trees      *routes.TreesSet
	InputTrees *routes.TreesSet
	BGTrees    *routes.TreesSet
	Objs_mdtc  [4]float64
	Timer      *algo_timer.Timer
}

type OSACO struct {
	Timeout       int
	K             int
	P             float64
	KTrees        *routes.KTreesSet
	VB            *Visibility
	PRM           *Pheromone
	InputTrees    *routes.TreesSet
	BGTrees       *routes.TreesSet
	Objs_osaco    [5][4]float64        // 200ms{o1, o2, o3, o4} 400ms{o1, o2, o3, o4} 600ms{o1, o2, o3, o4}, 800ms{o1, o2, o3, o4}, 1000ms{o1, o2, o3, o4}
	Timer         [5]*algo_timer.Timer // 200ms{time} 400ms{time} 600ms{time}, 800ms{time}, 1000ms{time}
	Method_Number int                  // 0: TOP K minimum weight 1: Increasing Arithmetic Sequence 2: Average Arithmetic Sequence
}

type Visibility struct {
	TSN_VB [][]float64
	AVB_VB [][]float64
}

type Pheromone struct {
	TSN_PRM     [][]float64
	AVB_PRM     [][]float64
	CAN2TSN_PRM [][]float64
}

// OSRO structures (Path-based routing)
type SP struct {
	Paths      *routes.PathsSet
	InputPaths *routes.PathsSet
	BGPaths    *routes.PathsSet
	Objs_sp    [4]float64
	Timer      *algo_timer.Timer
}

type OSACO_Path struct {
	Timeout    int
	K          int
	P          float64
	KPaths     *routes.KPathsSet
	VB_Path    *VisibilityPath
	PRM_Path   *PheromonePath
	InputPaths *routes.PathsSet
	BGPaths    *routes.PathsSet
	Objs_osaco [5][4]float64
	Timer      [5]*algo_timer.Timer
}

type VisibilityPath struct {
	TSN_VB     [][]float64
	AVB_VB     [][]float64
	CAN2TSN_VB [][]float64
}

type PheromonePath struct {
	TSN_PRM     [][]float64
	AVB_PRM     [][]float64
	CAN2TSN_PRM [][]float64
}
