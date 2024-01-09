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
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"gitlab.com/gpt4batch"
	"testing"
)

func Test_client_Upload(t *testing.T) {
	cc := NewClient()

	resp, err := cc.Upload(context.Background(), &gpt4batch.UploadRequest{
		Source: &gpt4batch.Source{
			ID:          uuid.New().String(),
			UploadURL:   "https://beta.gpt4api.plus/standard/uploaded",
			Name:        "upload-test",
			AccessToken: "v2.local.VbhT_pdpjICiOqmEHJkEtC__5CVY9BmjsrV2ORiYPTDysqINzv2IyJW5fLaeNIn_uXQoq_5zfmSFSbpasAaJQm6vx8IjMdGgxjiF3QKc24T2os_3rUXRiw-LvJ8G2eiPeNtOJ6RquMrLEgZuEwLoaaN0t2RkyYEMH6bVSXYvIzSyY3gvw1FIUWpYZDwu9edVH6IMcit0HHVZ3LoPfZp8r82g5KlOBoSnHKhbCOZsY1MrqianvsCi3v8nNFjxdAb9110vPn4og-qtwZqQHTU3ldqZ7170sq4Vb4OhosX5iTCaPQ89uzh7SR97oDOtEu3oBXMW3APAsWXKXQ.bnVsbA",
		},
		UploadPath:     "github.com/gpt4batch/client/unescape.go",
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
			AccessToken: "v2.local.VbhT_pdpjICiOqmEHJkEtC__5CVY9BmjsrV2ORiYPTDysqINzv2IyJW5fLaeNIn_uXQoq_5zfmSFSbpasAaJQm6vx8IjMdGgxjiF3QKc24T2os_3rUXRiw-LvJ8G2eiPeNtOJ6RquMrLEgZuEwLoaaN0t2RkyYEMH6bVSXYvIzSyY3gvw1FIUWpYZDwu9edVH6IMcit0HHVZ3LoPfZp8r82g5KlOBoSnHKhbCOZsY1MrqianvsCi3v8nNFjxdAb9110vPn4og-qtwZqQHTU3ldqZ7170sq4Vb4OhosX5iTCaPQ89uzh7SR97oDOtEu3oBXMW3APAsWXKXQ.bnVsbA",
		},
		Model:   "gpt-4",
		Message: "你是gpt3还是gpt4?",
	})

	assert.NoError(t, err)
	t.Log(resp)

	// 我是基于 GPT-4 架构的人工智能模型。
}
