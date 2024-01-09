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
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"gitlab.com/gpt4batch"
)

func Test_client_Upload(t *testing.T) {
	cc := NewClient()

	resp, err := cc.Upload(context.Background(), &gpt4batch.UploadRequest{
		Source: &gpt4batch.Source{
			ID:          uuid.New().String(),
			UploadURL:   "https://beta.gpt4api.plus/standard/uploaded",
			Name:        "upload-test",
			AccessToken: "",
		},
		UploadPath:     "unescape.go",
		ConversationId: "",
		UploadType:     "my_files",
	})

	assert.NoError(t, err)
	t.Log(resp)
}

func Test_client_Chat(t *testing.T) {
	cc := NewClient()

	resp, err := cc.Chat(context.Background(), &gpt4batch.ChatRequest{
		Source: &gpt4batch.Source{
			ID:          uuid.New().String(),
			URL:         "https://beta.gpt4api.plus/standard/all-tools",
			Name:        "upload-test",
			AccessToken: "",
		},
		Model:   "gpt-4",
		Message: "你是gpt3还是gpt4?",
	})

	assert.NoError(t, err)
	t.Log(resp)

	// 我是基于 GPT-4 架构的人工智能模型。
}
