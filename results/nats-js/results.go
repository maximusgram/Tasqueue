package nats

import (
	"context"
	"fmt"

	"github.com/nats-io/nats.go"
)

const (
	resultPrefix = "tasqueue-results-"
	kvBucket     = "tasqueue"
)

type Results struct {
	opt  Options
	conn nats.KeyValue
}

type Options struct {
	URL         string
	EnabledAuth bool
	Username    string
	Password    string
}

// New() returns a new instance of nats-jetstream broker.
func New(cfg Options) (*Results, error) {
	opt := []nats.Option{}

	if cfg.EnabledAuth {
		opt = append(opt, nats.UserInfo(cfg.Username, cfg.Password))
	}

	conn, err := nats.Connect(cfg.URL, opt...)
	if err != nil {
		return nil, fmt.Errorf("error connecting to nats : %w", err)
	}

	// Get jet stream context
	js, err := conn.JetStream()
	if err != nil {
		return nil, fmt.Errorf("error creating jetstream context : %w", err)
	}

	kv, err := js.KeyValue(kvBucket)
	if err != nil {
		return nil, fmt.Errorf("error creating key/value bucket : %w", err)
	}

	return &Results{
		opt:  cfg,
		conn: kv,
	}, nil
}

func (r *Results) Get(ctx context.Context, uuid string) ([]byte, error) {
	rs, err := r.conn.Get(resultPrefix + uuid)
	if err != nil {
		return nil, err
	}

	return rs.Value(), nil
}

func (r *Results) Set(ctx context.Context, uuid string, b []byte) error {
	if _, err := r.conn.Put(resultPrefix+uuid, b); err != nil {
		return err
	}
	return nil
}