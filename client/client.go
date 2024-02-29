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
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/go-resty/resty/v2"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"gitlab.com/gpt4batch"
)

// client is a client that logs requests and responses.
type client struct{}

// Upload uploads a file to the server.
func (c *client) Upload(ctx context.Context, req *gpt4batch.UploadRequest) (*gpt4batch.UploadResponse, error) {
	fileBytes, err := os.ReadFile(req.UploadPath)
	if err != nil {
		return nil, err
	}

	resp, err := resty.New().
		SetTimeout(30*time.Second).
		R().
		EnableTrace().
		SetAuthToken(req.AccessToken).
		SetHeader("Content-Type", "multipart/form-data").
		SetFormData(map[string]string{
			"conversation_id": req.ConversationId,
			"type":            req.UploadType,
		}).
		SetFileReader("file", filepath.Base(req.UploadPath), bytes.NewReader(fileBytes)).
		Post(req.UploadURL)

	if err != nil {
		return nil, err
	}

	if resp.StatusCode() != http.StatusOK {
		return nil, fmt.Errorf("failed to upload file: %s", resp.Status())
	}

	var result gpt4batch.UploadResponse
	if err = json.Unmarshal(resp.Body(), &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// Chat sends a message to the server.
func (c *client) Chat(ctx context.Context, req *gpt4batch.ChatRequest) (*gpt4batch.ChatResponse, error) {
	resp, err := resty.New().
		SetTimeout(8*time.Minute).
		R().
		EnableTrace().
		SetAuthToken(req.AccessToken).
		SetHeader("Content-Type", "application/json").
		SetBody(req.Openai()).
		Post(req.URL)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode() != http.StatusOK {
		return nil, fmt.Errorf("failed to chat: %s", resp.Status())
	}

	var result gpt4batch.ChatResponse
	if err := json.Unmarshal(resp.Body(), &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// Download downloads a file from the server.
func (c *client) Download(ctx context.Context, req *gpt4batch.DownloadRequest) error {
	resp, err := resty.New().
		SetTimeout(5 * time.Second).
		R().
		EnableTrace().
		Get(req.URL)
	if err != nil {
		return err
	}

	if resp.StatusCode() != http.StatusOK {
		return fmt.Errorf("failed to download file: %s", resp.Status())
	}

	// TODO: check if the file exists
	localPath := filepath.Join(req.LocalDir, req.LocalFileName)

	file, err := os.Create(localPath)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = io.Copy(file, resp.RawBody())
	if err != nil {
		return err
	}
	return nil
}

// Close closes the client.
func (c *client) Close(ctx context.Context) error {
	return nil
}

// NewClient returns a new client.
func NewClient() gpt4batch.Client {
	return &client{}
}
