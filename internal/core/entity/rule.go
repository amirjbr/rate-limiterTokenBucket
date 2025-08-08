package entity

type RateLimitRule struct {
	IntervalPerPermit int
	RefillTime        int
	BurstTokens       int
	Limit             int
	Interval          int
}
