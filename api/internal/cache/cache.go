package cache

import (
	"github.com/go-redis/redis"
	"fmt"
	"encoding/json"
	"time"
	"github.com/gambarini/cabapi/api/internal/model"
	"log"
)

type (
	Cache struct {
		rClient *redis.Client
	}
)

func NewCache() (cache *Cache, err error) {

	rClient := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})

	pong, err := rClient.Ping().Result()

	if err != nil {
		return cache, fmt.Errorf("failed to connect to redis, %s", err)
	}

	log.Printf("Ping Redis %s", pong)

	return &Cache{
		rClient: rClient,
	}, nil

}

func (cache *Cache) Close() {
	cache.rClient.Close()
}

func (cache *Cache) Set(data model.Trips) error {

	key := data.Medallion + data.Date

	log.Printf("Set cache, %s", key)

	dataJSON, err := json.Marshal(&data)

	if err != nil {
		return fmt.Errorf("failed to marshal to cache, %s", err)
	}

	err = cache.rClient.Set(key, dataJSON, 0).Err()

	if err != nil {
		return fmt.Errorf("failed to cache, %s", err)
	}

	return nil
}

func (cache *Cache) Get(medallion string, date time.Time) (data model.Trips, err error) {

	key := medallion + fmt.Sprintf("%d-%d-%d", date.Year(), date.Month(), date.Day())

	log.Printf("Get cache, %s", key)

	val, err := cache.rClient.Get(key).Result()

	if err != nil {
		return data, fmt.Errorf("failed to get cache, %s", err)
	}

	err = json.Unmarshal([]byte(val), &data)

	if err != nil {
		return data, fmt.Errorf("failed to unmarshal when get cache, %s", err)
	}

	return data, nil
}
