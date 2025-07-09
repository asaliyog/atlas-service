package cache

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"golang-service/internal/models"
	"github.com/redis/go-redis/v9"
)

// RedisCache handles Redis operations for VM data caching
type RedisCache struct {
	client *redis.Client
}

// NewRedisCache creates a new Redis cache instance
func NewRedisCache(addr string, password string, db int) *RedisCache {
	client := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: password,
		DB:       db,
	})

	return &RedisCache{client: client}
}

// GetVMs retrieves VMs from cache
func (rc *RedisCache) GetVMs(ctx context.Context) ([]models.VM, error) {
	key := "vms:all"
	data, err := rc.client.Get(ctx, key).Result()
	if err != nil {
		if err == redis.Nil {
			return nil, nil // Cache miss
		}
		return nil, fmt.Errorf("failed to get VMs from cache: %w", err)
	}

	var vms []models.VM
	if err := json.Unmarshal([]byte(data), &vms); err != nil {
		return nil, fmt.Errorf("failed to unmarshal VMs from cache: %w", err)
	}

	log.Printf("Retrieved %d VMs from cache", len(vms))
	return vms, nil
}

// SetVMs stores VMs in cache with 24-hour expiry
func (rc *RedisCache) SetVMs(ctx context.Context, vms []models.VM) error {
	key := "vms:all"
	data, err := json.Marshal(vms)
	if err != nil {
		return fmt.Errorf("failed to marshal VMs for cache: %w", err)
	}

	// Set with 24-hour expiry
	expiry := 24 * time.Hour
	err = rc.client.Set(ctx, key, data, expiry).Err()
	if err != nil {
		return fmt.Errorf("failed to set VMs in cache: %w", err)
	}

	log.Printf("Stored %d VMs in cache with %v expiry", len(vms), expiry)
	return nil
}

// InvalidateVMs removes VMs from cache
func (rc *RedisCache) InvalidateVMs(ctx context.Context) error {
	key := "vms:all"
	err := rc.client.Del(ctx, key).Err()
	if err != nil {
		return fmt.Errorf("failed to invalidate VMs cache: %w", err)
	}

	log.Println("Invalidated VMs cache")
	return nil
}

// Close closes the Redis connection
func (rc *RedisCache) Close() error {
	return rc.client.Close()
}

// Ping tests the Redis connection
func (rc *RedisCache) Ping(ctx context.Context) error {
	return rc.client.Ping(ctx).Err()
} 