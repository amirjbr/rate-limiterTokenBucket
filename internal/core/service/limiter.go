package service

import (
	"TokenBucketRateLimiter/internal/core/entity"
	"TokenBucketRateLimiter/internal/core/port/redisPort"
	"context"
	"errors"
	"fmt"
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
                local padMillis = math.fmod(intervalSinceLast, intervalPerPermit) 
                redis.call('hset', key, 'lastRefillTime', refillTime - padMillis)
            end
            currentTokens = math.min(grantedTokens + tokensRemaining, limit)
        end
    else
        currentTokens = tokensRemaining
    end
end


assert(currentTokens >= 0)

if currentTokens == 0 then
    redis.call('hset', key, 'tokensRemaining', currentTokens)
    return 0
else
    redis.call('hset', key, 'tokensRemaining', currentTokens - 1)
    return 1
end
`

var RateLimitExceededError = "rate limit exceeded"

type Service struct {
	rdb redisPort.RedisRepo
}

func NewLimiterService(rdb redisPort.RedisRepo) *Service {
	return &Service{
		rdb: rdb,
	}
}

func (s *Service) Limit(ctx context.Context, userID, ip, destinationService, method string) (bool, error) {
	//LoadRules from redis and define bucket limit and interval based on entity then send it to redis and see is it allow or not
	var rule entity.RateLimitRule
	ruleKey := fmt.Sprintf("%s:%s", destinationService, method)
	res, err := s.rdb.GetRule(ctx, ruleKey)
	if err != nil {
		return false, err
	}
	if len(res) == 0 {
		return false, errors.New("no rate limit rule")
	}

	validatedRule, err := FillRuleFields(res, rule)

	bucketKey := fmt.Sprintf("userID:%s:ip:%s:service:%s:method:%s", userID, ip, destinationService, method)
	keys := []string{bucketKey}

	evalResult, err := s.rdb.EvaluateScript(ctx, rateLimiterLua, keys, validatedRule)
	if err != nil {
		log.Fatal(err)
	}
	if evalResult == 1 {
		return true, nil
	}
	return false, errors.New(RateLimitExceededError)

}

func FillRuleFields(data map[string]string, rule entity.RateLimitRule) (entity.RateLimitRule, error) {
	var err error
	rule.Limit, err = strconv.Atoi(data["limit"])
	if err != nil {
		return entity.RateLimitRule{}, err
	}
	rule.BurstTokens, err = strconv.Atoi(data["burstTokens"])
	if err != nil {
		return entity.RateLimitRule{}, err
	}
	rule.Interval, err = strconv.Atoi(data["interval"])
	if err != nil {
		return entity.RateLimitRule{}, err
	}
	rule.RefillTime = int(time.Now().UnixMilli())
	rule.IntervalPerPermit, err = strconv.Atoi(data["intervalPerPermit"])
	if err != nil {
		return entity.RateLimitRule{}, err
	}
	err = ValidateRuleFields(rule)
	if err != nil {
		return entity.RateLimitRule{}, errors.New("invalid rule ")
	}

	return rule, nil
}

func ValidateRuleFields(rule entity.RateLimitRule) error {
	fields := []struct {
		name  string
		value int
	}{
		{"Limit", rule.Limit},
		{"BurstTokens", rule.BurstTokens},
		{"Interval", rule.Interval},
		{"RefillTime", rule.RefillTime},
		{"IntervalPerPermit", rule.IntervalPerPermit},
	}

	for _, f := range fields {
		if f.value == 0 {
			return fmt.Errorf("invalid rate limit field: %s", f.name)
		}
	}
	return nil
}
