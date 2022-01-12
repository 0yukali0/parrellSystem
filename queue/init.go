package queue

import (
	"container/heap"
)

var (
	Events = make(EventPQ, 0)
	Jobs = make(JobPQ, 0)
)

func init() {
	heap.Init(GetJobsQueue())
	heap.Init(GetEventsQueue())
}

func GetEventsQueue() *EventPQ{
	return &Events
}

func GetJobsQueue() *JobPQ{
	return &Jobs
}