package common

import (
	"fmt"
)

const (
	FilePath = "C:\\Users\\a08h0\\研究生\\研一\\parral\\simulation\\SDSC-SP2-1998-4.2-cln.swf"
)

var (
	ProcessNum  = uint64(128)
	limit = uint64(0)
	EventStatus = []string{"Submit", "Running"}
)

func Check(e error) {
	if e != nil {
		fmt.Printf("%v\n", e)
		panic(e)
	}
}

func TryAllocate(req uint64,allocated bool) bool{
	if req > ProcessNum || allocated {
		return false
	}
	return true
}

func Allocate(req uint64,allocated bool) bool{
	if req > ProcessNum || allocated {
		return false
	}
	ProcessNum -= req
	return true
}

func Release(alloc uint64, allocated bool) {
	if allocated {
		ProcessNum += alloc
	}
}