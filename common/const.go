package common

const (
	FilePath = "C:\\Users\\a08h0\\研究生\\研一\\parral\\simulation\\backfilling\\parrellSystem\\SDSC-SP2-1998-4.2-cln.swf"
)

var (
	ProcessNum  = uint64(128)
	SystemClock = uint64(0)
)

func Check(e error) {
	if e != nil {
		panic(e)
	}
}

func TryAllocate(req uint64, allocated bool) bool{
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

func GetCurrentProcessNum() uint64 {
	return ProcessNum
}

func GetSystemClock() uint64 {
	return SystemClock
}

func SetSystemClock(currentTime uint64) {
	SystemClock = currentTime
}