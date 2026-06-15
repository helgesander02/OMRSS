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

// OsroMemorizer keeps every aggregate (cumulative → average → CSV →
// result.txt) keyed by encap method so OSRO can report per-method
// results side by side. Methods carries the canonical iteration order
// (which matches plan.OSRO.Methods), so output rendering is
// deterministic across runs.
type OsroMemorizer struct {
	methods []string

	// SP baseline objective per method.
	average_obj_sp_by_method  map[string]*[4]float64
	average_time_sp_by_method map[string]*time.Duration

	// OSACO objective per method × per timeout. 5 timeouts:
	// 200ms / 400ms / 600ms / 800ms / 1000ms.
	average_objs_osaco_by_method map[string]*[5][4]float64
	average_time_osaco_by_method map[string]*[5]time.Duration
}

func newOsroMemorizer() *OsroMemorizer {
	return &OsroMemorizer{
		average_obj_sp_by_method:     make(map[string]*[4]float64),
		average_time_sp_by_method:    make(map[string]*time.Duration),
		average_objs_osaco_by_method: make(map[string]*[5][4]float64),
		average_time_osaco_by_method: make(map[string]*[5]time.Duration),
	}
}

// ensureMethod lazily creates zero-valued aggregates for `m`. Called by
// MCumulative on first sighting so OSRO doesn't have to pre-populate
// the memorizer with method names.
func (OS *OsroMemorizer) ensureMethod(m string) {
	if _, ok := OS.average_obj_sp_by_method[m]; ok {
		return
	}
	OS.methods = append(OS.methods, m)
	OS.average_obj_sp_by_method[m] = &[4]float64{}
	OS.average_time_sp_by_method[m] = new(time.Duration)
	OS.average_objs_osaco_by_method[m] = &[5][4]float64{}
	OS.average_time_osaco_by_method[m] = &[5]time.Duration{}
}
