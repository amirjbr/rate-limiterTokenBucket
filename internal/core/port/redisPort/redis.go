package redisPort

import (
	"TokenBucketRateLimiter/internal/core/entity"
	"context"
)

type RedisRepo interface {
	GetRule(ctx context.Context, ruleKey string) (map[string]string, error)
	EvaluateScript(ctx context.Context, script string, keys []string, rule entity.RateLimitRule) (int, error)
}
