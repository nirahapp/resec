package redis

import (
	"github.com/go-redis/redis"
	"github.com/jpillora/backoff"
	log "github.com/sirupsen/logrus"

	"github.com/nirahapp/resec/resec/state"
)

type Manager struct {
	backoff        *backoff.Backoff // exponential backoff helper
	logger         *log.Entry       // logging specificall for Redis
	client         *redis.Client    // redis client
	config         *Config          // redis config
	state          *state.Redis     // redis state
	stateCh        chan state.Redis // redis state channel to publish updates to the reconciler
	commandCh      chan Command
	stopCh         chan interface{}
	watcherRunning bool
}

type Config struct {
	Address string // address (IP+Port) used to talk to Redis
}
