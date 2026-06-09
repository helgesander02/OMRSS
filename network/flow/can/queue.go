package can

import "sort"

type Queue struct {
	Frames []*Frame
}

func newQueue() *Queue {
	return &Queue{}
}

func (q *Queue) appendQueue(frames []*Frame) {
	q.Frames = append(q.Frames, frames...)
}

func (q *Queue) popQueueByIdx(idx int) {
	if idx < 0 || idx >= len(q.Frames) {
		return
	}
	q.Frames = append(q.Frames[:idx], q.Frames[idx+1:]...)
}

func (q *Queue) popQueueByHead(head int) {
	q.Frames = q.Frames[head:]
}

func (q *Queue) checkDrop(currentTime int) int {
	o1Drop := 0
	for drop := len(q.Frames) - 1; drop >= 0; drop-- {
		if currentTime > q.Frames[drop].FinishTime {
			o1Drop++
			q.popQueueByIdx(drop)
		}
	}

	return o1Drop
}

func (q *Queue) sortQueue(method string, currentTime int) {
	switch method {
	case "fifo":
		// arrival time (small → large)
		sort.Slice(q.Frames, func(i, j int) bool {
			return q.Frames[i].ArrivalTime < q.Frames[j].ArrivalTime
		})

	case "priority":
		// deadline (small → large)
		sort.Slice(q.Frames, func(i, j int) bool {
			return q.Frames[i].Deadline < q.Frames[j].Deadline
		})

	case "wst":
		// finish time - current time (small → large)
		sort.Slice(q.Frames, func(i, j int) bool {
			ti := q.Frames[i].FinishTime - currentTime
			tj := q.Frames[j].FinishTime - currentTime
			return ti < tj
		})

	case "wst_mar":
		// MAR (Multi-tier Adaptive Release) uses EDF order so the head of
		// the queue is always the most urgent frame. We sort by FinishTime
		// directly because currentTime is uniform within a step, so subtracting
		// it from every key would not change the order.
		sort.Slice(q.Frames, func(i, j int) bool {
			return q.Frames[i].FinishTime < q.Frames[j].FinishTime
		})

	case "mao":
		// deadline (small → large)
		sort.Slice(q.Frames, func(i, j int) bool {
			return q.Frames[i].Deadline < q.Frames[j].Deadline
		})

	default:
		// default: FIFO
		sort.Slice(q.Frames, func(i, j int) bool {
			return q.Frames[i].ArrivalTime < q.Frames[j].ArrivalTime
		})
	}
}
