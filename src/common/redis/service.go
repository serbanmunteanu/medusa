package redis

import (
	"encoding/json"
	"github.com/go-redis/redis/v8"
	"golang.org/x/net/context"
	"time"
)

type Service struct {
	rds *redis.Client
}

func NewRedisService(redisUrl string) (*Service, error) {
	options, err := redis.ParseURL(redisUrl)
	if err != nil {
		return nil, err
	}
	return &Service{
		rds: redis.NewClient(options),
	}, nil
}

func (s *Service) HasKey(key string) (bool, error) {
	cmd, err := s.rds.Do(context.Background(), "EXISTS", key).Bool()
	if err != nil {
		return false, err
	}
	return cmd, nil
}

func (s *Service) Get(key string, value interface{}) error {
	result, err := s.rds.Get(context.Background(), key).Bytes()
	if err != nil {
		return err
	}
	value = json.Unmarshal(result, &value)
	return nil
}

func (s *Service) Set(key string, value interface{}, ttl time.Duration) error {
	bytes, err := json.Marshal(value)
	if err != nil {
		return err
	}
	return s.rds.Set(context.Background(), key, bytes, ttl).Err()
}
