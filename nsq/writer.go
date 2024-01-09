/*
Copyright 2023 The gpt4batch Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package nsq

import (
	"context"
	"crypto/tls"
	"errors"
	"gitlab.com/gpt4batch"
	"io"
	llog "log"
	"sync"

	"github.com/asaskevich/govalidator"
	"github.com/nsqio/go-nsq"
)

// NSQConfig is the configuration for NSQ.
type NSQConfig struct {
	Enable      bool
	Address     string
	UserAgent   string
	MaxInFlight int
	Topic       string
}

// NewNSQConfig creates a new NSQ configuration.
func NewNSQConfig() NSQConfig {
	return NSQConfig{
		Enable:      true,
		Address:     "127.0.0.1:4150",
		UserAgent:   "",
		MaxInFlight: 64,
		Topic:       "gpt4api",
	}
}

// Validate validates the NSQ configuration.
func (c NSQConfig) Validate() error {
	if govalidator.IsNull(c.Address) {
		return errors.New("nsq address is required")
	}
	return nil
}

// NSQWriter is an NSQ writer.
type nsqWriter struct {
	logger gpt4batch.Logger
	conf   NSQConfig

	tlsConf  *tls.Config
	connMut  sync.RWMutex
	producer *nsq.Producer
}

// NewNSQWriter creates a new NSQ writer.
func NewNSQWriter(conf NSQConfig, logger gpt4batch.Logger) (Async, error) {
	n := nsqWriter{
		logger: logger.WithField("nsq", "svc"),
		conf:   conf,
	}
	return &n, nil
}

// Connect connects to NSQ.
func (n *nsqWriter) Connect(ctx context.Context) error {
	n.connMut.Lock()
	defer n.connMut.Unlock()

	cfg := nsq.NewConfig()
	cfg.UserAgent = n.conf.UserAgent
	if n.tlsConf != nil {
		cfg.TlsV1 = true
		cfg.TlsConfig = n.tlsConf
	}

	producer, err := nsq.NewProducer(n.conf.Address, cfg)
	if err != nil {
		return err
	}

	producer.SetLogger(
		llog.New(io.Discard, "", llog.Flags()),
		nsq.LogLevelError,
	)

	if err := producer.Ping(); err != nil {
		return err
	}
	n.producer = producer

	n.logger.Info("Sending NSQ messages to address: ", n.conf.Address)
	return nil
}

// WriteBatch writes a message to NSQ.
func (n *nsqWriter) WriteBatch(ctx context.Context, topic string, msg []byte) error {
	n.connMut.RLock()
	prod := n.producer
	n.connMut.RUnlock()

	if prod == nil {
		return errors.New("not connected to target source or sink")
	}

	if len(msg) == 0 {
		return nil
	}
	return prod.Publish(topic, msg)
}

// Close closes the NSQ writer.
func (n *nsqWriter) Close(context.Context) error {
	n.connMut.Lock()
	defer n.connMut.Unlock()

	if n.producer != nil {
		n.producer.Stop()
		n.producer = nil
	}
	return nil
}
