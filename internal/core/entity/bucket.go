package entity

type Bucket struct {
	Key               string
	IntervalPerPermit int
	RefillTime        int
	BurstTokens       int
	Limit             int
	Interval          int
}
