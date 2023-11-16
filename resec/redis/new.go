package redis

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"os"
	"time"

	"github.com/go-redis/redis"
	"github.com/jpillora/backoff"
	log "github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"

	"github.com/nirahapp/resec/resec/state"
)

func NewConnection(m *cli.Context) (*Manager, error) {
	redisConfig := &Config{
		Address: m.String("redis-addr"),
	}
	var redisClient *redis.Client

	if m.Bool("tls") {
		cert, err := tls.LoadX509KeyPair(m.Path("tls-cert-file"), m.Path("tls-key-file"))
		if err != nil {
			return nil, fmt.Errorf("failed to load TLS key pair: %s", err)
		}
		caCert, err := os.ReadFile(m.Path("tls-ca-cert-file"))
		if err != nil {
			return nil, fmt.Errorf("failed to read CA certificate: %s", err)
		}
		caCertPool := x509.NewCertPool()
		caCertPool.AppendCertsFromPEM(caCert)
		tlsConfig := &tls.Config{
			Certificates:       []tls.Certificate{cert},
			RootCAs:            caCertPool,
			InsecureSkipVerify: m.Bool("tls-insecure-skip-verify"),
		}

		redisClient = redis.NewClient(&redis.Options{
			Addr:         redisConfig.Address,
			DialTimeout:  1 * time.Second,
			Password:     m.String("redis-password"),
			ReadTimeout:  1 * time.Second,
			WriteTimeout: 1 * time.Second,
			TLSConfig:    tlsConfig,
		})
	} else {
		redisClient = redis.NewClient(&redis.Options{
			Addr:         redisConfig.Address,
			DialTimeout:  1 * time.Second,
			Password:     m.String("redis-password"),
			ReadTimeout:  1 * time.Second,
			WriteTimeout: 1 * time.Second,
		})
	}

	instance := &Manager{
		client: redisClient,
		config: redisConfig,
		logger: log.WithFields(log.Fields{
			"system":     "redis",
			"redis_addr": m.String("redis-addr"),
		}),
		state:     &state.Redis{},
		stateCh:   make(chan state.Redis, 10),
		commandCh: make(chan Command, 10),
		stopCh:    make(chan interface{}, 1),
		backoff: &backoff.Backoff{
			Min:    50 * time.Millisecond,
			Max:    10 * time.Second,
			Factor: 1.5,
			Jitter: false,
		},
	}

	return instance, nil
}
