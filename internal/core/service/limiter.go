package service

import (
	"TokenBucketRateLimiter/internal/core/entity"
	"context"
	"fmt"
	"github.com/redis/go-redis/v9"
	"log"
	"strconv"
	"time"
)

const rateLimiterLua = `
local key, intervalPerPermit, refillTime, burstTokens = KEYS[1], tonumber(ARGV[1]), tonumber(ARGV[2]), tonumber(ARGV[3])
local limit, interval = tonumber(ARGV[4]), tonumber(ARGV[5])
local bucket = redis.call('hgetall', key)

local currentTokens

if table.maxn(bucket) == 0 then
    currentTokens = burstTokens
    redis.call('hset', key, 'lastRefillTime', refillTime)
elseif table.maxn(bucket) == 4 then
    local lastRefillTime, tokensRemaining = tonumber(bucket[2]), tonumber(bucket[4])
    if refillTime > lastRefillTime then
        local intervalSinceLast = refillTime - lastRefillTime
        if intervalSinceLast > interval then
            currentTokens = burstTokens
            redis.call('hset', key, 'lastRefillTime', refillTime)
        else
            local grantedTokens = math.floor(intervalSinceLast / intervalPerPermit)
            if grantedTokens > 0 then
                local padMillis = math.fmod(intervalSinceLast, intervalPerPermit) //tttttttttttt
                redis.call('hset', key, 'lastRefillTime', refillTime - padMillis)
            end
            currentTokens = math.min(grantedTokens + tokensRemaining, limit)
        end
    else
        currentTokens = tokensRemaining
    end
end

////////tttttttttttttttttt
assert(currentTokens >= 0)

if currentTokens == 0 then
    redis.call('hset', key, 'tokensRemaining', currentTokens)
    return 0
else
    redis.call('hset', key, 'tokensRemaining', currentTokens - 1)
    return 1
end
`

type Service struct {
	redis *redis.Client
}

func NewLimiterService(rdb *redis.Client) *Service {
	return &Service{
		redis: rdb,
	}
}

func (s *Service) Limit(ctx context.Context, userID, ip, destinationService, method string) bool {
	//LoadRules from redis and define bucket limit and interval based on entity then send it to redis and see is it allow or not
	var rule entity.RateLimitRule
	ruleKey := fmt.Sprintf("%s:%s", destinationService, method)

	res, err := s.redis.HGetAll(ctx, ruleKey).Result()
	if err != nil {
		log.Fatal(err)
	}
	// need to make this operation as a function and handle errors and validations there
	rule.Limit, err = strconv.Atoi(res["limit"])
	rule.BurstTokens, err = strconv.Atoi(res["burstTokens"])
	rule.Interval, err = strconv.Atoi(res["interval"])
	rule.RefillTime = int(time.Now().UnixMilli())
	rule.IntervalPerPermit, err = strconv.Atoi(res["intervalPerPermit"])

	bucketKey := fmt.Sprintf("userID:%s:ip:%s:service:%s:method:%s", userID, ip, destinationService, method)
	keys := []string{bucketKey}

	rateLimitResult, err := s.redis.Eval(ctx, rateLimiterLua, keys, rule.IntervalPerPermit, rule.RefillTime, rule.BurstTokens, rule.Limit, rule.Interval).Result()
	if err != nil {
		log.Fatal(err)
	}

	if rateLimitResult == int64(1) {
		return true
	}

	return false

}
