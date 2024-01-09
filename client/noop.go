/*
Copyright 2024 The gpt4batch Authors.

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
	"github.com/go-faker/faker/v4"
	"gitlab.com/gpt4batch"
	"math/rand"
	"time"
)

type noop struct{}

func (s *noop) Close(ctx context.Context) error {
	return nil
}

func (s *noop) Upload(ctx context.Context, req *gpt4batch.UploadRequest) (*gpt4batch.UploadResponse, error) {
	sleepSeconds := rand.Intn(10) + 2
	time.Sleep(time.Duration(sleepSeconds) * time.Second)
	return &gpt4batch.UploadResponse{
		Attachment: &gpt4batch.Attachment{
			Id:            faker.Username(),
			Name:          faker.Username(),
			Size:          100,
			FileTokenSize: 10,
			MimeType:      faker.Username(),
			Width:         10,
			Height:        10,
		},
		Part: nil,
	}, nil
}

func (s *noop) Chat(ctx context.Context, req *gpt4batch.ChatRequest) (*gpt4batch.ChatResponse, error) {
	sleepSeconds := rand.Intn(10) + 2
	time.Sleep(time.Duration(sleepSeconds) * time.Second)

	return &gpt4batch.ChatResponse{
		Created:        time.Now().Unix(),
		MessageID:      faker.Name(),
		ConversationID: faker.Username(),
		EndTurn:        false,
		Contents:       []interface{}{faker.Username()},
	}, nil
}

func (s *noop) Download(ctx context.Context, req *gpt4batch.DownloadRequest) error {
	return nil
}

func NewNoop() gpt4batch.Client {
	return &noop{}
}
