package cache

import (
	"io"
	"github.com/go-redis/redis/v8"
)

// Cache is a closer interface
type Cache interface {
	io.Closer
}

type cache struct {
	conn *redis.Client
}

func (c *cache) Close() error {
	return c.conn.Close()
}
