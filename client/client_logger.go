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
	"fmt"
	"path/filepath"
	"time"

	"gitlab.com/gpt4batch"
)

// clientLogger is a client that logs requests and responses.
type clientLogger struct {
	logger gpt4batch.Logger
	svc    gpt4batch.Client
}

func (c clientLogger) Close(ctx context.Context) error {
	return nil
}

// Upload uploads a file to the server.
func (c clientLogger) Upload(ctx context.Context, req *gpt4batch.UploadRequest) (resp *gpt4batch.UploadResponse, err error) {
	logg := c.logger.
		WithField("url", req.UploadURL).
		WithField("filepath", req.UploadPath).
		WithField("id", req.ID).
		WithField("pid", req.Pid).
		WithField("conversation_id", req.ConversationId)

	defer func(start time.Time) {
		if err != nil {
			logg.
				WithField("took", time.Since(start)).
				Error(fmt.Sprintf("Failed to upload: %s", err))
			return
		}
		logg.
			WithField("took", time.Since(start)).
			Info("Upload")
	}(time.Now())
	return c.svc.Upload(ctx, req)
}

// Chat sends a message to the server.
func (c clientLogger) Chat(ctx context.Context, req *gpt4batch.ChatRequest) (resp *gpt4batch.ChatResponse, err error) {
	logg := c.logger.
		WithField("url", req.URL).
		WithField("model", req.Model).
		WithField("gizmo_id", req.GizmoId).
		WithField("id", req.ID).
		WithField("pid", req.Pid).
		WithField("conversation_id", req.ConversationID)

	defer func(start time.Time) {
		if err != nil {
			logg.
				WithField("took", time.Since(start)).
				Error(fmt.Sprintf("Failed to chat: %s", err))
			return
		}
		logg.
			WithField("end_turn", resp.EndTurn).
			WithField("download", len(resp.Downloads)).
			WithField("took", time.Since(start)).
			Info("Chat")
	}(time.Now())
	return c.svc.Chat(ctx, req)
}

// Download downloads a file from the server.
func (c clientLogger) Download(ctx context.Context, req *gpt4batch.DownloadRequest) (err error) {
	logg := c.logger.
		WithField("url", req.URL).
		WithField("id", req.ID).
		WithField("pid", req.Pid).
		WithField("localPath", filepath.Join(req.LocalDir, req.LocalFileName))

	defer func(start time.Time) {
		if err != nil {
			logg.
				WithField("took", time.Since(start)).
				Error(fmt.Sprintf("Failed to download: %s", err))
			return
		}
		logg.
			WithField("took", time.Since(start)).
			Info("Download")
	}(time.Now())
	return c.svc.Download(ctx, req)
}

// NewClientLogger returns a new client that logs requests and responses.
func NewClientLogger(logger gpt4batch.Logger, svc gpt4batch.Client) gpt4batch.Client {
	return &clientLogger{
		logger: logger,
		svc:    svc,
	}
}
