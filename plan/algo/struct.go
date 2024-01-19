package algo

import "src/plan/routes"

type SMT struct {
	Trees      *routes.Trees_set
	InputTrees *routes.Trees_set
	BGTrees    *routes.Trees_set
	Objs_smt   [4]float64
}

type MDTC struct {
	Trees      *routes.Trees_set
	InputTrees *routes.Trees_set
	BGTrees    *routes.Trees_set
	Objs_mdtc  [4]float64
}

type OSACO struct {
	Timeout    int
	K          int
	P          float64
	KTrees     *routes.KTrees_set
	VB         *Visibility
	PRM        *Pheromone
	InputTrees *routes.Trees_set
	BGTrees    *routes.Trees_set
	Objs_osaco [5][4]float64 // 1000ms{o1, o2, o3, o4} 800ms{o1, o2, o3, o4} 600ms{o1, o2, o3, o4}, 400ms{o1, o2, o3, o4}, 200ms{o1, o2, o3, o4}
}

type Visibility struct {
	TSN_VB [][]float64
	AVB_VB [][]float64
}

type Pheromone struct {
	TSN_PRM [][]float64
	AVB_PRM [][]float64
}
