package main

import (
	"container/heap"

	"simulation/object"
	"simulation/reader"
	"simulation/queue"
	"simulation/common"
	"simulation/framework"
)



func main() {
	/*
	readFileAndInit()
	framework.FCFS()
	readFileAndInit()
	framework.Preempt()
	*/
	readFileAndInit()
	framework.Backfill()
}

func readFileAndInit() {
	result := reader.ReadFile(common.FilePath)
	submitAndFinishQueue := queue.GetEventsQueue()

	// assign event to event queue
	var firstSubmitTime string
	for idx, meta := range result {
		if idx == 0 {
			firstSubmitTime = meta.Submit
		}
		event := object.NewEvent(&meta, firstSubmitTime)

		if idx == 0 {
			common.SetSystemClock(event.GetTimeStamp())
		}
		heap.Push(submitAndFinishQueue, event)
	}
}