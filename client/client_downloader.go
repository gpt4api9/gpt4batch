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

	"github.com/asaskevich/govalidator"

	"gitlab.com/gpt4batch"
)

// clientDownloader is a client that downloads files from the server.
type clientDownloader struct {
	svc gpt4batch.Client
	// enable is a flag that enables the client.
	enable bool
}

// Upload uploads a file to the server.
func (c clientDownloader) Upload(ctx context.Context, req *gpt4batch.UploadRequest) (*gpt4batch.UploadResponse, error) {
	return c.svc.Upload(ctx, req)
}

// Chat sends a message to the server.
func (c clientDownloader) Chat(ctx context.Context, req *gpt4batch.ChatRequest) (*gpt4batch.ChatResponse, error) {
	resp, err := c.svc.Chat(ctx, req)
	if err != nil {
		return nil, err
	}

	if c.enable {
		if len(resp.Downloads) != 0 {
			for _, download := range resp.Downloads {
				dreq := &downloadReq{
					Id:     req.ID,
					Pid:    req.Pid,
					URL:    download,
					Prefix: req.Prefix,
				}

				localFileName := DownloadUrlPath(dreq)
				if !govalidator.IsNull(localFileName) {
					resp.SpecDownloads = append(resp.SpecDownloads, &gpt4batch.SpecDownload{
						Origin: download,
						Local:  localFileName,
					})

					// todo add a new method to download a file from the server.
					go c.Download(ctx, &gpt4batch.DownloadRequest{
						Source: &gpt4batch.Source{
							URL: download,
						},
						LocalDir:      req.Dir,
						LocalFileName: localFileName,
					})
				}
			}
		}
	}
	return resp, nil
}

// Download downloads a file from the server.
func (c clientDownloader) Download(ctx context.Context, req *gpt4batch.DownloadRequest) error {
	return c.svc.Download(ctx, req)
}

func (c clientDownloader) Close(ctx context.Context) error {
	return nil
}

// NewClientDownloader returns a new client that downloads files from the server.
func NewClientDownloader(enable bool, svc gpt4batch.Client) gpt4batch.Client {
	return &clientDownloader{
		enable: enable,
		svc:    svc,
	}
}
