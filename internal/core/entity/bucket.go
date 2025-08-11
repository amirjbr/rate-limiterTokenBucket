package entity

// talk about this
type Bucket struct {
	Key               string
	IntervalPerPermit int
	RefillTime        int
	BurstTokens       int
	Limit             int
	Interval          int
}
