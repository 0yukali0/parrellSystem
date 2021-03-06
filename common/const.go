package common

const (
	FilePath = "C:\\Users\\a08h0\\研究生\\研一\\parral\\simulation\\backfilling\\parrellSystem\\SDSC-SP2-1998-4.2-cln.swf"
)

const (
	DefaultCPUNum = uint64(128)
	DefaultTimeStart = uint64(0)
	DefaultTimeLimit = ^uint64(0)
	PreemptActive = true
	BackfillActive = true
)

var (
	ProcessNum  = DefaultCPUNum
	SystemClock = DefaultTimeStart
	BaseSubmitTime uint64
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

func Allocate(req uint64, allocated bool) bool{
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

func GetSystemCapcity() uint64 {
	return DefaultCPUNum
}

func GetSystemClock() uint64 {
	return SystemClock
}

func SetSystemClock(currentTime uint64) {
	SystemClock = currentTime
}