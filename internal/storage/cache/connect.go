package cache

import (
	"context"
	"time"

	"github.com/namsral/flag"

	"github.com/go-redis/redis/v8"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

var (
	// fs           = flag.NewFlagSet("apigateway", flag.ExitOnError)
	redisAddr    = flag.String("redis-addr", "localhost:6379", "")
	redisSecret  = flag.String("redis-secret", "localhost:6379", "")
	redisDB      = flag.Int64("redis-db", 0, "")
	redisTimeout = flag.Int64("redis-timeout-ms", 2000, "")
)

func init() {

}

// Connect makes a new redis Connection.
func connect() (*redis.Client, error) {
	conn := redis.NewClient(&redis.Options{
		Addr:     *redisAddr,
		Password: *redisSecret,
		DB:       int(*redisDB),
	})

	// check if the redis is running
	if err := waitForDB(conn); err != nil {
		return nil, err
	}

	return conn, nil
}

// New creates a new databse
func New() (Cache, error) {

	conn, err := connect()
	if err != nil {
		return nil, err
	}

	redis := &cache{conn: conn}
	return redis, nil
}

func waitForDB(conn *redis.Client) error {

	var ctx = context.Background()
	ready := make(chan struct{})
	go func() {
		for {
			if _, err := conn.Ping(ctx).Result(); err == nil {
				logrus.Debug("Cache Connected")
				close(ready)
				return
			}
			time.Sleep(1 * time.Millisecond)
		}
	}()

	select {
	case <-ready:
		return nil
	case <-time.After(time.Duration(*redisTimeout) * time.Millisecond):
		return errors.New("redis not ready")
	}
}
