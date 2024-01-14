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
	"context"
	"encoding/json"
	"gitlab.com/gpt4batch/signals"
	"net/http"
	"time"

	"github.com/spf13/cobra"
	"gitlab.com/gpt4batch"
	"gitlab.com/gpt4batch/client"
	"gitlab.com/gpt4batch/log"
	"gitlab.com/gpt4batch/nsq"
	"gitlab.com/gpt4batch/reader"
)

// NewBatchCommand returns a new cobra.Command for launching the batchsvc.
func NewBatchCommand(ctx context.Context) *cobra.Command {
	var (
		option Option
		logger = log.New(log.InfoLevel)
	)

	rootCmd := &cobra.Command{
		Use:   "batchsvc",
		Args:  cobra.NoArgs,
		Short: "please proceed with caution when enabling batch calls to GP Ts scripts.",
		RunE: func(cmd *cobra.Command, args []string) error {
			logg := logger.
				WithField("url", option.URL).
				WithField("model", option.Model).
				WithField("fix", option.Fix).
				WithField("qps", option.QPS).
				WithField("rdb", option.EnableRDB).
				WithField("gizmo_id", option.GizmoId)

			// validate the option. if the option is invalid, return an error.
			if err := option.Validate(); err != nil {
				return err
			}

			var (
				// batchTotal is the total number of batches.
				batchTotal uint64 = 0
				// asks is the gpt4api batch.
				ins = make(gpt4batch.Ins, 0)
				// cc is the client.
				cc = client.NewClientDownloader(option.EnableDownload, client.NewClientLogger(logger, client.NewClient()))
				// todo NewNoop only use to test.
				//cc = client.NewNoop()
			)

			// read the input file. if the input file is invalid, return an error.
			// the input file is a json file. each line is a json object.
			// the json object is a gpt4api batch.
			if err := reader.Reader(option.In, func(le string) error {
				// parse the json object.
				// if the json object is invalid, return an error.
				// the json object is a gpt4api batch.
				in := new(gpt4batch.In)
				if err := json.Unmarshal([]byte(le), in); err != nil {
					return err
				}

				// If the run is not continued and there are no errors
				// set the error message to Resource Not Ready.
				if !option.Fix && in.IErr == nil {
					in.IErr = &gpt4batch.IErr{
						Code:    http.StatusBadRequest,
						Message: "resource is not ready",
					}
				}

				// append the gpt4api batch to the asks.
				// the asks is a gpt4api batch.
				ins = append(ins, in)

				// increase the batch total.
				// the batch total is the total number of batches.
				batchTotal++

				logg.
					WithField("id", in.ID).
					WithField("count", batchTotal).
					Info("Reader")
				return nil
			}); err != nil {
				return err
			}

			// NSQ is enabled. create a new NSQ writer.
			// the NSQ writer is used to send the gpt4api batch to the server.
			if option.NSQ.Enable {
				// create a new NSQ writer.
				async, err := nsq.NewNSQWriter(option.NSQ, logger)
				if err != nil {
					return err
				}
				// connect to NSQ. if the connection is failed, return an error.
				if err = async.Connect(ctx); err != nil {
					return err
				}

				// create a new client that uses NSQ. the client is used to send the gpt4api batch to the server.
				// the client is a wrapper of the NSQ writer.
				cc = client.NewClientNSQ(logger, async, cc, option.NSQ.Topic)
			}

			// create a new service. the service is used to send the gpt4api batch to the server.
			svc := NewService(&option, cc, ins, &Stats{
				BatchTotal:    batchTotal,
				CompleteTotal: 0,
				SuccessTotal:  0,
				FailedTotal:   0,
			})

			// set the logger for the service.
			// the logger is used to log the service.
			svc.WithLogger(logg)

			// open the service. if the service is failed to open, return an error.
			// the service is used to send the gpt4api batch to the server.
			if err := svc.Open(signals.WithStandardSignals(ctx)); err != nil {
				return err
			}

			<-svc.Done()

			// Tear down the launcher, allowing it a few seconds to finish any
			// in-progress requests.
			shutdownCtx, cancel := context.WithTimeout(ctx, 6*time.Minute)
			defer cancel()
			return svc.Close(shutdownCtx)
		},
	}

	rootCmd.Flags().StringVarP(&option.URL, "url", "u", "https://beta.gpt4api.plus/standard/all-tools", "设置批量调用服务地址.普通版：standard 并发版：concurrent")
	rootCmd.Flags().StringVarP(&option.UploadURL, "upload_url", "l", "https://beta.gpt4api.plus/standard/uploaded", "设置批量调用服务地址.普通版：standard 并发版：concurrent")
	rootCmd.Flags().StringVarP(&option.In, "in", "i", "example.jsonl", "输入文件路径，数据格式按照规定格式定义.")
	rootCmd.Flags().StringVarP(&option.Out, "out", "o", "out.jsonl", "输出文件路径，GPTs数据跑完存储数据的文件路径.")
	rootCmd.Flags().IntVarP(&option.Goroutine, "goroutine", "g", 60, "设置最大协程数量.")
	rootCmd.Flags().BoolVarP(&option.HistoryAndTrainingDisabled, "history_and_training_disabled", "s", true, "是否开启历史对话历史记录，默认是关闭的.")
	rootCmd.Flags().StringVarP(&option.Model, "model", "m", "gpt-4-gizmo", "设置调用GPTs的模型.")
	rootCmd.Flags().StringVarP(&option.GizmoId, "gizmo-id", "z", "", "设置GPTs gizmo id的名称.")
	rootCmd.Flags().BoolVarP(&option.Fix, "fix", "f", false, "是否开启续跑模式.")
	rootCmd.Flags().IntVarP(&option.QPS, "qps", "q", 1, "设置QPS并发量.")
	rootCmd.Flags().BoolVarP(&option.NSQ.Enable, "enable_nsq", "n", false, "是否开启NSQ消息队列.")
	rootCmd.Flags().BoolVarP(&option.EnableDownload, "enable-download", "e", true, "是否开启文件下载.")
	rootCmd.Flags().StringVarP(&option.DownloadDir, "download-dir", "d", "", "下载文件夹名称.如果未设置会存在当前文件夹目录.")
	rootCmd.Flags().StringVarP(&option.DownloadFilePrefix, "download-prefix", "p", "GPT4API", "设置文件下载前缀，防止下载文件名冲突覆盖.")
	rootCmd.Flags().BoolVarP(&option.EnableRDB, "rdb", "r", true, "是否开启RDB文件缓存持久化策略.")
	rootCmd.Flags().IntVarP(&option.RDBInterval, "rdb_interval", "v", 60, "RDB缓存时间间隔，默认是60分钟")
	return rootCmd
}
