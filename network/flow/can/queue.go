package can

import "sort"

type Queue struct {
	Streams []*Stream
}

func new_Queue() *Queue {
	return &Queue{}
}

func (q *Queue) appendQueue(streams []*Stream) {
	q.Streams = append(q.Streams, streams...)
}

func (q *Queue) popQueueByIdx(idx int) {
	if idx < 0 || idx >= len(q.Streams) {
		return
	}
	q.Streams = append(q.Streams[:idx], q.Streams[idx+1:]...)
}

func (q *Queue) popQueueByHead(head int) {
	q.Streams = q.Streams[head:]
}

func (q *Queue) checkDrop(current_time int) int {
	o1_drop := 0
	for drop := len(q.Streams) - 1; drop >= 0; drop-- {
		if current_time > q.Streams[drop].FinishTime {
			o1_drop++
			q.popQueueByIdx(drop)
		}
	}

	return o1_drop
}

func (q *Queue) sortQueue(method string, current_time int) {
	switch method {
	case "fifo":
		// arrival time (small → large)
		sort.Slice(q.Streams, func(i, j int) bool {
			return q.Streams[i].ArrivalTime < q.Streams[j].ArrivalTime
		})

	case "priority":
		// deadline (small → large)
		sort.Slice(q.Streams, func(i, j int) bool {
			return q.Streams[i].Deadline < q.Streams[j].Deadline
		})

	case "wst":
		// finish time - current time (small → large)
		sort.Slice(q.Streams, func(i, j int) bool {
			ti := q.Streams[i].FinishTime - current_time
			tj := q.Streams[j].FinishTime - current_time
			return ti < tj
		})

	case "mao":
		// deadline (small → large)
		sort.Slice(q.Streams, func(i, j int) bool {
			return q.Streams[i].Deadline < q.Streams[j].Deadline
		})

	default:
		// default: FIFO
		sort.Slice(q.Streams, func(i, j int) bool {
			return q.Streams[i].ArrivalTime < q.Streams[j].ArrivalTime
		})
	}
}

func (q *Queue) hasImminent(now, safe int) bool {
	for _, s := range q.Streams {
		if s.FinishTime-now <= safe {
			return true
		}
	}
	return false
}
