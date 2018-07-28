package redis

import (
	"github.com/YotpoLtd/resec/resec/state"
	"github.com/go-redis/redis"
	log "github.com/sirupsen/logrus"
	cli "gopkg.in/urfave/cli.v1"
)

func NewConnection(m *cli.Context) (*Manager, error) {
	redisConfig := &Config{
		Address: m.String("redis-addr"),
	}

	instance := &Manager{
		client: redis.NewClient(&redis.Options{
			Addr:        redisConfig.Address,
			DialTimeout: m.Duration("healthcheck-timeout"),
			Password:    m.String("redis-password"),
			ReadTimeout: m.Duration("healthcheck-timeout"),
		}),
		config: redisConfig,
		logger: log.WithFields(log.Fields{
			"system":     "redis",
			"redis_addr": m.String("redis-addr"),
		}),
		state:     &state.Redis{},
		stateCh:   make(chan state.Redis, 1),
		commandCh: make(chan Command, 1),
		stopCh:    make(chan interface{}, 1),
	}

	return instance, nil
}
