package port

import "context"

type Limiter interface {
	Limit(ctx context.Context, userID, ip, destinationService, method string) bool
}
