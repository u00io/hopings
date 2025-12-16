package system

type ResultHop struct {
	IP          string
	CountryName string
	CountryISO  string
	TimeMs      int64
}

type Result struct {
	Status      string
	IP          string
	CountryName string
	CountryISO  string
	Hops        []*ResultHop
}

func NewResult() *Result {
	var c Result
	return &c
}
