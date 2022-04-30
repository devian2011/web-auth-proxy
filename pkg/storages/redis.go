package storages

import (
	"encoding/json"
	"fmt"
	"github.com/go-redis/redis"
	"github.com/golang/glog"
	"time"
)

const (
	checkRedisConnectionTimeout = 5
)

type RedisConfig struct {
	Host                 string        `json:"host"`
	Port                 string        `json:"port"`
	Password             string        `json:"password"`
	ReconnectionTimeout  time.Duration `json:"reconnection_timeout"`
	ReconnectionAttempts int           `json:"reconnection_attempts"`
}

type RedisStorage struct {
	config      *RedisConfig
	client      *redis.Client
	isConnected bool
	shutdown    bool
}

func InitRedis(config *RedisConfig) *RedisStorage {
	storage := &RedisStorage{config: config}
	storage.initRedisClient()
	if !isRedisAvailable(storage.client) {
		glog.Exit("Application cannot work without redis")
	}
	storage.isConnected = true
	storage.shutdown = false
	go storage.checkStorageConnection()

	return storage
}

func (s *RedisStorage) IsClientConnected() bool {
	return s.isConnected
}

func (s *RedisStorage) GetAll(collection string, callback CollectionUnmarshalFunc) (map[string]interface{}, error) {
	if s.isConnected {
		result, err := s.client.HGetAll(collection).Result()
		if err != nil {
			return nil, err
		}

		out := make(map[string]interface{}, 0)
		for id, entity := range result {
			out[id] = callback(entity)
		}

		return out, nil
	}

	return nil, &ConnectionError{Message: "Redis is not connected"}
}

func (s *RedisStorage) Get(collection string, id string, value interface{}) error {
	if s.isConnected {
		result, err := s.client.HGet(collection, id).Result()
		if err != nil {
			return err
		}
		err = json.Unmarshal([]byte(result), &value)
		if err != nil {
			return err
		}

		return nil
	}

	return &ConnectionError{Message: "Redis is not connected"}
}

func (s *RedisStorage) Add(collection string, id string, value interface{}) error {
	if s.isConnected {
		marshal, err := json.Marshal(value)
		if err != nil {
			return err
		}
		return s.client.HSet(collection, id, marshal).Err()
	}
	return &ConnectionError{Message: "Redis is not connected"}
}

func (s *RedisStorage) Update(collection string, id string, value interface{}) error {
	return s.Add(collection, id, value)
}

func (s *RedisStorage) Remove(collection string, id string) error {
	if s.isConnected {
		return s.client.HDel(collection, id).Err()
	}

	return &ConnectionError{Message: "Redis is not connected"}
}

func (s *RedisStorage) initRedisClient() {
	s.client = redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%s", s.config.Host, s.config.Port),
		Password: s.config.Password,
	})
}

func isRedisAvailable(cl *redis.Client) bool {
	_, err := cl.Ping().Result()
	if err != nil {
		glog.Errorf("Redis is not available. Reason %s \n", err.Error())
	}
	return err == nil
}

func (s *RedisStorage) checkStorageConnection() {
	for {
		//Close goroutine if stop calls application stop.
		if s.shutdown == true {
			return
		}
		time.Sleep(time.Second * checkRedisConnectionTimeout)
		if !isRedisAvailable(s.client) {
			s.isConnected = false
			for counter := 1; counter <= s.config.ReconnectionAttempts; counter++ {
				time.Sleep(time.Second * s.config.ReconnectionTimeout)
				s.initRedisClient()
				if isRedisAvailable(s.client) {
					s.isConnected = true
				} else {
					s.Stop(false)
				}
			}
		}
	}
}

func (s *RedisStorage) Stop(shutdown bool) {
	err := s.client.Close()
	s.shutdown = shutdown
	if err != nil {
		glog.Error(err)
	}

}
