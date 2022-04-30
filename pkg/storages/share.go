package storages

import "github.com/go-redis/redis"

type CollectionUnmarshalFunc func(jsonInput string) interface{}

const (
	EmptyResult = redis.Nil
)

type ConnectionError struct {
	Message string
}

func (ce *ConnectionError) Error() string {
	return ce.Message
}
