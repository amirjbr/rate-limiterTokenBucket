package repository

import (
	"TokenBucketRateLimiter/internal/core/entity"
	"context"
	"errors"
	"github.com/redis/go-redis/v9"
	"log"
)

type RedisPortImpl struct {
	client *redis.Client
}

func NewRedisImpl(client *redis.Client) *RedisPortImpl {
	return &RedisPortImpl{client: client}
}

func (r *RedisPortImpl) GetRule(ctx context.Context, ruleKey string) (map[string]string, error) {
	if ruleKey == "" {
		return nil, errors.New("empty rule key")
	}

	res, err := r.client.HGetAll(ctx, ruleKey).Result()
	return res, err
}

func (r *RedisPortImpl) EvaluateScript(ctx context.Context, script string, keys []string, rule entity.RateLimitRule) (int, error) {
	// TODO use evalSHA for better performance
	rateLimitResult, err := r.client.Eval(ctx, script, keys, rule.IntervalPerPermit, rule.RefillTime, rule.BurstTokens, rule.Limit, rule.Interval).Result()
	if err != nil {
		log.Fatal(err)
	}

	if rateLimitResult == int64(1) {
		return 1, nil
	}
	return 0, nil
}
