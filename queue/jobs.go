package queue

import (
	"simulation/object"
)

var (
	Jobs = make(JobPQ, 0)
)

func GetJobsQueue() JobPQ{
	return Jobs
}

type JobPQ []*object.Job

func (pq JobPQ) Len() int { return len(pq) }

func (pq JobPQ) Less(i, j int) bool {
	return pq[i].Submission < pq[j].Submission
}

func (pq JobPQ) IsEmpty() bool {
	if pq.Len() > 0 {
		return false
	}
	return true
}

func (pq JobPQ) Swap(i, j int) {
	pq[i], pq[j] = pq[j], pq[i]
	pq[i].Index = i
	pq[j].Index = j
} 

func (pq *JobPQ) Push(x interface{}) {
	n := len(*pq)
	item := x.(*object.Job)
	item.Index = n
	*pq = append(*pq, item)
}

func (pq *JobPQ) Pop() interface{} {
	old := *pq
	n := len(old)
	item := old[n-1]
	old[n-1] = nil  // avoid memory leak
	item.Index = -1 // for safety
	*pq = old[0 : n-1]
	return item
}