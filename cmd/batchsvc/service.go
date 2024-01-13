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

package batchsvc

import (
	"bufio"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"
	"syscall"
	"time"

	"github.com/panjf2000/ants/v2"
	"github.com/schollz/progressbar/v3"

	"github.com/asaskevich/govalidator"
	"gitlab.com/gpt4batch"
	"gitlab.com/gpt4batch/log"
)

// service implements gpt4batch.Service.
type service struct {
	// logger is the service logger.
	logger gpt4batch.Logger
	// stats is the service stats.
	stats *Stats
	// config is the service config.
	config Option
	// client is the service client.
	cc gpt4batch.Client
	// items is the service items.
	items gpt4batch.Ins
	// rdbInterval is the service rdb interval.
	rdbInterval time.Duration
	// progressBar	is the service progress bar.
	progressBar *progressbar.ProgressBar
	// doneChan is the service done channel.
	cancel func()
	// doneChan is the service done channel.
	doneChan <-chan struct{}
	// pool is the service pool.
	pool *ants.Pool
}

// NewService returns a new gpt4batch.Service.
func NewService(config Option, cc gpt4batch.Client, items gpt4batch.Ins, stats *Stats) gpt4batch.Service {
	svc := &service{
		logger:      log.New(log.InfoLevel),
		config:      config,
		stats:       stats,
		cc:          cc,
		items:       items,
		rdbInterval: time.Duration(config.RDBInterval) * time.Minute,
		progressBar: progressbar.Default(int64(len(items))),
	}

	// pool is the service pool.
	// if the pool is null, return an error.
	pool, err := ants.NewPool(config.Goroutine)
	if err != nil {
		panic(err)
	}

	// pool is the service pool.
	svc.pool = pool
	return svc
}

// Done is shutdown.
func (s *service) Done() <-chan struct{} {
	return s.doneChan
}

// Open opens the service.
func (s *service) Open(ctx context.Context) error {
	ctx, s.cancel = context.WithCancel(ctx)
	s.doneChan = ctx.Done()

	s.logger.Info("Start")

	// if the items is null, return an error.
	// if the items is not null, do the work.
	go s.doWork(ctx)

	// rdb is the cache in local file.
	// if the rdb is null, return an error.
	go s.rdb(ctx)

	return nil
}

func (s *service) doWork(ctx context.Context) {

	// items
	for _, item := range s.items {
		if s.config.Fix && item.IErr == nil {
			s.stats.IncrSuccessCount()
			s.updateProgressBar(ctx)
			// quickly skip fields that have already succeeded.
			continue
		}

		// if the item is not null, do the work.
		if err := s.pool.Submit(func() {
			if err := s.Chat(ctx, item); err != nil {
				s.stats.IncrFailedCount()
				s.updateProgressBar(ctx)
				item.IErr = &gpt4batch.IErr{
					Code:    http.StatusNotImplemented,
					Message: err.Error(),
				}
			} else {
				s.stats.IncrSuccessCount()
				s.updateProgressBar(ctx)
			}
		}); err != nil {
			s.logger.
				WithField("event", "submit").
				Error(err)
			item.IErr = &gpt4batch.IErr{
				Code:    http.StatusNotImplemented,
				Message: err.Error(),
			}
			s.stats.IncrSuccessCount()
			s.updateProgressBar(ctx)
			continue
		}
	}
}

// Chat sends a message to the server and returns the response.
// if the response is not null, append the response.
func (s *service) Chat(ctx context.Context, in *gpt4batch.In) error {
	var (
		// conversationID is the conversation id.
		conversationID string
		// parentMessageID is the parent message id.
		parentMessageID string
		// answers is the answers.
		answers []interface{}
		// attachments is the attachments.
		attachments gpt4batch.Attachments
		// parts is the parts.
		parts gpt4batch.Parts
	)

	if ctx.Err() != nil {
		return errors.New("service canceled")
	}

	logger := s.logger.
		WithField("id", in.ID)

	for _, ask := range in.Asks {
		logger.
			WithField("pid", ask.ID).
			Info()

		// tmpConversationID is the temporary conversation id.
		// if the conversation id is not null, use the conversation id.
		var tmpConversationID string

		// conversationID is the conversation id. if the conversation id is not null, use the conversation id.
		// if the conversation id is null, use the temporary conversation id.
		if !govalidator.IsNull(conversationID) {
			tmpConversationID = conversationID
		}

		// Images is the images. if the images is not null, upload the images.
		// if the images is null, do nothing.
		if len(ask.Images) != 0 {
			for _, image := range ask.Images {
				// resp is the response. if the response is not null, upload the image.
				// if the response is null, do nothing.
				resp, err := s.cc.Upload(ctx, &gpt4batch.UploadRequest{
					Source: &gpt4batch.Source{
						ID:          in.ID,                       // in.ID is the id of the batch.
						URL:         s.config.UploadURL,          // s.config.URL is the url of the server.
						Name:        "UploadImage",               // "Upload" is the name of the upload.
						Pid:         ask.ID,                      // ask.ID is the id of the ask.
						Prefix:      s.config.DownloadFilePrefix, // Prefix  is the download file prefix.
						AccessToken: s.config.AccessToken,        // s.config.AccessToken is the access token of the server.
						Dir:         s.config.DownloadDir,        // s.config.DownloadDir is the download dir of the server.
					},
					ConversationId: tmpConversationID,
					UploadPath:     image,
					UploadType:     gpt4batch.Multimodal,
				})
				if err != nil {
					logger.
						WithField("pid", ask.ID).
						Error(err)
					return err
				}

				// parts is the parts. if the parts is not null, append the parts.
				// if the parts is null, do nothing.
				parts = append(parts, &gpt4batch.Part{
					AssetPointer: resp.Part.AssetPointer,
					SizeBytes:    resp.Part.SizeBytes,
					Width:        resp.Part.Width,
					Height:       resp.Part.Height,
				})
			}
		}

		// Files is the files. if the files is not null, upload the files.
		if len(ask.Files) != 0 {
			for _, file := range ask.Files {
				// resp is the response. if the response is not null, upload the file.
				// if the response is null, do nothing.
				resp, err := s.cc.Upload(ctx, &gpt4batch.UploadRequest{
					Source: &gpt4batch.Source{
						ID:          in.ID,                       // in.ID is the id of the batch.
						URL:         s.config.UploadURL,          // s.config.URL is the url of the server.
						Name:        "UploadFile",                // "Upload" is the name of the upload.
						Pid:         ask.ID,                      // ask.ID is the id of the ask.
						Prefix:      s.config.DownloadFilePrefix, // Prefix  is the download file prefix.
						AccessToken: s.config.AccessToken,        // s.config.AccessToken is the access token of the server.
						Dir:         s.config.DownloadDir,        // s.config.DownloadDir is the download dir of the server.
					},
					ConversationId: tmpConversationID,
					UploadPath:     file,
					UploadType:     gpt4batch.MyFiles,
				})
				if err != nil {
					logger.
						WithField("pid", ask.ID).
						Error(err)
					return err
				}

				// parts is the parts. if the parts is not null, append the parts.
				// if the parts is null, do nothing.
				attachments = append(attachments, &gpt4batch.Attachment{
					Id:            resp.Attachment.Id,
					Name:          resp.Attachment.Name,
					Size:          resp.Attachment.Size,
					FileTokenSize: resp.Attachment.FileTokenSize,
					MimeType:      resp.Attachment.MimeType,
					Width:         resp.Attachment.Width,
					Height:        resp.Attachment.Height,
				})
			}
		}

		// Chat sends a message to the server and returns the response.
		// if the response is not null, append the response.
		resp, err := s.cc.Chat(ctx, &gpt4batch.ChatRequest{
			Source: &gpt4batch.Source{
				ID:          in.ID,                       // in.ID is the id of the batch.
				URL:         s.config.URL,                // s.config.URL is the url of the server.
				Name:        "Chat",                      // "Chat" is the name of the chat.
				Pid:         ask.ID,                      // ask.ID is the id of the ask.
				Prefix:      s.config.DownloadFilePrefix, // Prefix  is the download file prefix.
				AccessToken: s.config.AccessToken,        // s.config.AccessToken is the access token of the server.
				Dir:         s.config.DownloadDir,        // s.config.DownloadDir is the download dir of the server.
			},
			GizmoId:                    s.config.GizmoId,
			Message:                    ask.Content,
			ParentMessageID:            parentMessageID,
			ConversationID:             tmpConversationID,
			Stream:                     false,
			Model:                      s.config.Model,
			Attachments:                attachments,
			Parts:                      parts,
			HistoryAndTrainingDisabled: s.config.HistoryAndTrainingDisabled,
		})
		if err != nil {
			logger.
				WithField("pid", ask.ID).
				Error(err)
			return err
		}

		// answers is the answers. if the answers is not null, append the answers.
		// if the answers is null, do nothing.
		answers = append(answers, resp)
		parentMessageID = resp.MessageID
		conversationID = resp.ConversationID
	}

	// in.Answers is the answers. if the answers is not null, append the answers.
	// if the answers is null, do nothing.
	if len(answers) != 0 {
		in.Answers = answers
		in.IErr = nil
		return nil
	}

	logger.Error("chat answer is required")
	return fmt.Errorf("chat answer is required")
}

// kill the service.
func (s *service) kill(ctx context.Context) error {
	// kill the process. wait 4 minutes.
	pid := os.Getpid()
	return syscall.Kill(pid, syscall.SIGTERM)
}

// updateProgressBar increases the complete total.
func (s *service) updateProgressBar(ctx context.Context) {
	// incr increases the complete total.
	s.progressBar.Add(1)
	// incr increases the complete total.
	s.stats.IncrCompleteCount()

	// if the complete total is equal to the batch total, kill the service.
	if s.stats.GetCompleteTotal() == s.stats.GetBatchTotal() {
		s.progressBar.Finish()
		s.logger.
			WithField("event", "incr").
			Info("Done")

		// kill the service.
		s.kill(ctx)
	}
}

// rdb is the cache in local file.
func (s *service) rdb(ctx context.Context) {
	t := time.NewTimer(s.rdbInterval)
	defer t.Stop()

	for {
		select {
		case <-ctx.Done():
			s.logger.
				WithField("rdb", "sync").
				Info("Done")
			return
		case <-t.C:
			// tempFilename is the temporary filename.
			// if the temporary filename is not null, use the temporary filename.
			tempFilename := gpt4batch.TempFileName(s.config.In)

			// if the complete total is not equal to the batch total, do nothing.
			// if the complete total is equal to the batch total, write the items to the temporary filename.
			if s.stats.GetCompleteTotal() != s.stats.GetBatchTotal() {
				logger := s.logger.
					WithField("event", "rdb").
					WithField("tempFilename", tempFilename)

				logger.Info("Start")
				if err := s.write(ctx, tempFilename); err != nil {
					logger.Error(err)
					continue
				}

				s.logger.Info("RDB Completed")
				// reset rdb interval.
				t.Reset(s.rdbInterval)
			}
		}
	}
}

// WithLogger sets the logger for the service.
func (s *service) WithLogger(log gpt4batch.Logger) {
	s.logger = log.
		WithField("service", "batchsvc").
		WithField("total", s.stats.GetBatchTotal()).
		WithField("complete", s.stats.GetCompleteTotal()).
		WithField("success", s.stats.GetSuccessTotal()).
		WithField("failed", s.stats.GetFailedTotal())
}

// write writes the items to the file.
func (s *service) write(ctx context.Context, filename string) error {
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	// set batch writing, each batch writes 10000
	writer := bufio.NewWriter(file)
	// todo batch size 1000. maybe change in the future.
	batchSize := 10000

	// write the items to the temporary filename.
	// if the items is null, do nothing.
	for idx, item := range s.items {
		jsonStr, err := json.Marshal(item)
		if err != nil {
			continue
		}

		// write the items to the temporary filename.
		// if the items is null, do nothing.
		_, err = writer.WriteString(string(jsonStr) + "\n")
		if err != nil {
			continue
		}

		// flush the writer.
		// if the idx is equal to the batch size, flush the writer.
		if idx%batchSize == 0 {
			if err := writer.Flush(); err != nil {
				continue
			}

			// sync the file.
			// if the idx is equal to the batch size, sync the file.
			if err := file.Sync(); err != nil {
				continue
			}
		}
	}

	// finally refresh the disk to prevent data loss.
	// this operation is important to prevent data loss.
	if err = writer.Flush(); err != nil {
		return err
	}

	// sync the file.
	// this operation is important to prevent data loss.
	return file.Sync()
}

// Close closes the service.
func (s *service) Close(ctx context.Context) error {
	s.logger.Info("Close")

	// release the pool.
	s.pool.Release()

	// if the cancel is not null, cancel the service.
	if s.cc != nil {
		if err := s.cc.Close(ctx); err != nil {
			return err
		}
	}
	return s.write(ctx, s.config.Out)
}
