package port

type Limiter interface {
	Limit(userID, ip, destinationService, method string) bool
}
