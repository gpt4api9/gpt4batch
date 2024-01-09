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

package client

import (
	"context"
	"encoding/json"
	"github.com/asaskevich/govalidator"
	"gitlab.com/gpt4batch"
	"gitlab.com/gpt4batch/nsq"
)

// clientNSQ is a client that uses NSQ.
type clientNSQ struct {
	// logger is the logger.
	logger gpt4batch.Logger
	// async is the async client.
	async nsq.Async
	// svc is the service client.
	svc gpt4batch.Client
	// topic is the topic of the client.
	Topic string
}

// Upload uploads a file to the server.
func (c clientNSQ) Upload(ctx context.Context, req *gpt4batch.UploadRequest) (*gpt4batch.UploadResponse, error) {
	return c.svc.Upload(ctx, req)
}

// nsqChatMessage is a message that is sent to the server.
type nsqChatMessage struct {
	// id is the id of the message.
	Ask interface{} `json:"ask"`
	// answer is the answer to the message.
	Answer interface{} `json:"answer"`
	// dir is the directory of the message.
	Dir string `json:"dir"`
}

// Chat sends a message to the server.
func (c clientNSQ) Chat(ctx context.Context, req *gpt4batch.ChatRequest) (*gpt4batch.ChatResponse, error) {
	resp, err := c.svc.Chat(ctx, req)
	if err != nil {
		return nil, err
	}

	if c.async != nil {
		obj := &nsqChatMessage{
			Ask:    req.Message,
			Answer: resp,
			Dir:    req.Dir,
		}

		body, _ := json.Marshal(obj)
		if err := c.async.WriteBatch(ctx, c.Topic, body); err != nil {
			c.logger.Warn("NSQ: failed to write batch: %s", err)
		}
	}
	return resp, nil
}

// Download downloads a file from the server.
func (c clientNSQ) Download(ctx context.Context, req *gpt4batch.DownloadRequest) error {
	return c.svc.Download(ctx, req)
}

// Close closes the client.
func (c clientNSQ) Close(ctx context.Context) error {
	if c.async != nil {
		return c.async.Close(ctx)
	}
	return nil
}

// NewClientNSQ returns a new client that uses NSQ.
func NewClientNSQ(logger gpt4batch.Logger, async nsq.Async, svc gpt4batch.Client, topic string) gpt4batch.Client {
	cc := &clientNSQ{
		logger: logger,
		async:  async,
		svc:    svc,
		Topic:  topic,
	}

	if govalidator.IsNull(cc.Topic) {
		cc.Topic = "gpt4api"
	}
	return cc
}
