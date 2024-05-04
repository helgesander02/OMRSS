package memorizer

import "time"

type OMACO_Memorizer struct {
	average_obj_smt    [4]float64    // {o1, o2, o3, o4}
	average_obj_mdt    [4]float64    // {o1, o2, o3, o4}
	average_objs_osaco [5][4]float64 // 200ms{o1, o2, o3, o4} 400ms{o1, o2, o3, o4} 600ms{o1, o2, o3, o4}, 800ms{o1, o2, o3, o4}, 1000ms{o1, o2, o3, o4}
	average_time_smt   time.Duration
	average_time_mdt   time.Duration
	average_time_osaco time.Duration
}

func new_OMACO_Memorizer() *OMACO_Memorizer {
	return &OMACO_Memorizer{}
}

//type Plan2_Computer struct {

//}

//func New_plan2_Memorizer() *Memorizer2 {
//	return
//}

//type Plan3_Computer struct {

//}

//func New_plan3_Memorizer() *Memorizer3 {
//	return
//}
