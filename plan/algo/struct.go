package algo

import (
	"src/plan/algo_timer"
	"src/plan/routes"
)

type SMT struct {
	Routes      *routes.RouteSet
	InputRoutes *routes.RouteSet
	BGRoutes    *routes.RouteSet
	Objs_smt    [4]float64
	Timer       *algo_timer.Timer
}

type MDTC struct {
	Routes      *routes.RouteSet
	InputRoutes *routes.RouteSet
	BGRoutes    *routes.RouteSet
	Objs_mdtc   [4]float64
	Timer       *algo_timer.Timer
}

type OSACO struct {
	Timeout       int
	K             int
	P             float64
	KRoutes       *routes.KRouteSet
	VB            *Visibility
	PRM           *Pheromone
	InputRoutes   *routes.RouteSet
	BGRoutes      *routes.RouteSet
	Objs_osaco    [5][4]float64        // 200ms{o1, o2, o3, o4} 400ms{o1, o2, o3, o4} 600ms{o1, o2, o3, o4}, 800ms{o1, o2, o3, o4}, 1000ms{o1, o2, o3, o4}
	Timer         [5]*algo_timer.Timer // 200ms{time} 400ms{time} 600ms{time}, 800ms{time}, 1000ms{time}
	Method_Number int                  // 0: TOP K minimum weight 1: Increasing Arithmetic Sequence 2: Average Arithmetic Sequence
}

type Visibility struct {
	TSN_VB   [][]float64
	AVB_VB   [][]float64
	CAN2TTVB [][]float64
}

type Pheromone struct {
	TSN_PRM   [][]float64
	AVB_PRM   [][]float64
	CAN2TTPRM [][]float64
}

type SP struct {
	Routes      *routes.RouteSet
	InputRoutes *routes.RouteSet
	BGRoutes    *routes.RouteSet
	Objs_sp     [4]float64
	Timer       *algo_timer.Timer
}
