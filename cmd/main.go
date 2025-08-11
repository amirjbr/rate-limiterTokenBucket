package main

import (
	"TokenBucketRateLimiter/config"
	"TokenBucketRateLimiter/internal/adapter/redisClient"
	"TokenBucketRateLimiter/internal/adapter/redisClient/repository"
	"TokenBucketRateLimiter/internal/app/httpserver"
	"TokenBucketRateLimiter/internal/core/service"
	"TokenBucketRateLimiter/pkg/observeutil"
	"context"
	"fmt"
	"github.com/joho/godotenv"
	"log"
)

var ctx = context.Background()

// Lua script from your example
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

func main() {
	err := godotenv.Load("./.env")
	if err != nil {
		log.Println("No .env file found, using system env vars instead")
	}

	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(cfg)

	//url.JoinPath()
	//addr := fmt.Sprintf("%v:%v", cfg.Redis.Host, cfg.Redis.Port)
	//fmt.Println(addr)
	//fmt.Println(cfg.Redis.Password)
	//rdb := redis.NewClient(&redis.Options{
	//	Addr: "localhost:6379",
	//	//Password: cfg.Redis.Password,
	//})

	rdb := redisClient.NewRedisClient(cfg)
	redisRepo := repository.NewRedisImpl(rdb)

	fmt.Println("*************DEBUG")
	fmt.Printf("%+v\n", rdb)

	fmt.Println("step 1")
	limiterService := service.NewLimiterService(redisRepo)
	fmt.Println("step 2")
	metrics := observeutil.Setup()
	handlers := httpserver.NewHandler(limiterService)
	fmt.Println("step 3")
	server := httpserver.NewHttpServer(handlers, &metrics)
	fmt.Println("step 4")

	server.Engine.Run()

	//key := "user:123:rate_limit"
	//
	// total refill interval
	//
	//result, err := rdb.Eval(ctx, rateLimiterLua, []string{key},
	//	intervalPerPermit,
	//	refillTime,
	//	burstTokens,
	//	limit,
	//	interval,
	//).Result()
	//
	//if err != nil {
	//	log.Fatalf("Redis Eval failed: %v", err)
	//}
	//
	//if result.(int64) == 1 {
	//	fmt.Println(" Allowed")
	//} else {
	//	fmt.Println(" Rate limited")
	//}
}
