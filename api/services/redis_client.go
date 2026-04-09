package services

import (
	"context"
	"encoding/json"
	"strconv"
	"time"

	"api/globals"
	"api/log"

	redis "github.com/redis/go-redis/v9"
)

var (
	redisClient *redis.Client
	redisCtx    = context.Background()
)

// InitRedis initializes the Redis client if enabled via configuration.
func InitRedis() {
	if !globals.AppConfig.RedisEnabled {
		log.Info("Redis is disabled via configuration; using in-memory cache only")
		return
	}

	addr := globals.AppConfig.RedisHost
	if globals.AppConfig.RedisPort != "" {
		addr = addr + ":" + globals.AppConfig.RedisPort
	}

	redisClient = redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: globals.AppConfig.RedisPassword,
		DB:       globals.AppConfig.RedisDB,
	})

	if err := redisClient.Ping(redisCtx).Err(); err != nil {
  log.WithField("error", err).Warn("Failed to connect to Redis; falling back to in-memory cache")
		redisClient = nil
		return
	}

	log.WithField("addr", addr).Info("Connected to Redis")
}

// RedisAvailable returns true if the Redis client is initialized and healthy.
func RedisAvailable() bool {
	return redisClient != nil
}

// RedisGetFloat tries to get a float64 value from Redis.
// Returns (value, true) if found in Redis; otherwise (0, false).
func RedisGetFloat(key string) (float64, bool) {
	if redisClient == nil {
		return 0, false
	}
	val, err := redisClient.Get(redisCtx, key).Result()
	if err != nil {
		return 0, false
	}
	// Parse as float64
	f, err := parseFloat(val)
	if err != nil {
		return 0, false
	}
	return f, true
}

// RedisSetFloat sets a float64 value in Redis with the given TTL.
func RedisSetFloat(key string, value float64, ttl time.Duration) {
	if redisClient == nil {
		return
	}
	// Errors are logged but not fatal to flow
	if err := redisClient.Set(redisCtx, key, value, ttl).Err(); err != nil {
		log.WithField("error", err).Warn("Failed to set Redis key")
	}
}

// RedisGetJSON gets JSON value from Redis and unmarshals into v. Returns true on success.
func RedisGetJSON(key string, v interface{}) bool {
	if redisClient == nil {
		return false
	}
	bs, err := redisClient.Get(redisCtx, key).Bytes()
	if err != nil {
		return false
	}
	if err := json.Unmarshal(bs, v); err != nil {
		return false
	}
	return true
}

// RedisSetJSON marshals v and stores in Redis with TTL. Logs on error.
func RedisSetJSON(key string, v interface{}, ttl time.Duration) {
	if redisClient == nil {
		return
	}
	bs, err := json.Marshal(v)
	if err != nil {
		log.WithField("error", err).Warn("Failed to marshal JSON for Redis")
		return
	}
	if err := redisClient.Set(redisCtx, key, bs, ttl).Err(); err != nil {
		log.WithField("error", err).Warn("Failed to set Redis key")
	}
}

// parseFloat converts a string to float64 with minimal allocation.
func parseFloat(s string) (float64, error) { return strconv.ParseFloat(s, 64) }
