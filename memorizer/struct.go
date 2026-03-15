package memorizer

import "time"

type OmacoMemorizer struct {
	average_obj_smt          [4]float64    // {o1, o2, o3, o4}
	average_obj_mdt          [4]float64    // {o1, o2, o3, o4}
	average_objs_osaco       [5][4]float64 // 200ms{o1, o2, o3, o4} 400ms{o1, o2, o3, o4} 600ms{o1, o2, o3, o4}, 800ms{o1, o2, o3, o4}, 1000ms{o1, o2, o3, o4}
	average_objs_osaco_apted [5][4]float64 // 200ms{o1, o2, o3, o4} 400ms{o1, o2, o3, o4} 600ms{o1, o2, o3, o4}, 800ms{o1, o2, o3, o4}, 1000ms{o1, o2, o3, o4}
	average_time_mdt         time.Duration
	average_time_osaco       [5]time.Duration // 200ms{time} 400ms{time} 600ms{time}, 800ms{time}, 1000ms{time}
	average_time_osaco_apted [5]time.Duration // 200ms{time} 400ms{time} 600ms{time}, 800ms{time}, 1000ms{time}

}

func newOmacoMemorizer() *OmacoMemorizer {
	return &OmacoMemorizer{}
}

type OsroMemorizer struct {
	average_obj_smt          [4]float64    // {o1, o2, o3, o4}
	average_obj_mdt          [4]float64    // {o1, o2, o3, o4}
	average_objs_osaco       [5][4]float64 // 200ms{o1, o2, o3, o4} 400ms{o1, o2, o3, o4} 600ms{o1, o2, o3, o4}, 800ms{o1, o2, o3, o4}, 1000ms{o1, o2, o3, o4}
	average_objs_osaco_apted [5][4]float64 // 200ms{o1, o2, o3, o4} 400ms{o1, o2, o3, o4} 600ms{o1, o2, o3, o4}, 800ms{o1, o2, o3, o4}, 1000ms{o1, o2, o3, o4}
	average_time_mdt         time.Duration
	average_time_osaco       [5]time.Duration // 200ms{time} 400ms{time} 600ms{time}, 800ms{time}, 1000ms{time}
	average_time_osaco_apted [5]time.Duration // 200ms{time} 400ms{time} 600ms{time}, 800ms{time}, 1000ms{time}

}

func newOsroMemorizer() *OsroMemorizer {
	return &OsroMemorizer{}
}
