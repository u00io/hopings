package system

type ResultHop struct {
	IP     string
	TimeMs int64
}

type Result struct {
	Status  string
	IP      string
	Country string
	Hops    []*ResultHop
}

func NewResult() *Result {
	var c Result
	return &c
}
