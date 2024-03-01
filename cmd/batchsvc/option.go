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
	"errors"
	"github.com/asaskevich/govalidator"
	"gitlab.com/gpt4batch"
	"gitlab.com/gpt4batch/nsq"
	"os"
	"path/filepath"
)

// Option is the option.
type Option struct {
	// Addr is the address.
	// 批量调用服务地址
	URL string
	// UploadURL is the upload url.
	// 文件上传服务地址
	UploadURL string
	// In is the input.
	// 输入文件,读取用户即将批量跑的文件名称，默认是out.jsonl
	In string
	// Out is the output.
	// 输出文件，GPTs对话完成后输出的文件名称.默认当前目录
	Out string
	// goroutine is the goroutine.
	// 设置并发数
	Goroutine int
	// Model is the model.
	// 模型名称,GPT-3,GPT-4,GPT-4-Gizmo
	Model string
	// GizmoId is the gizmo id.
	// GPTs 模型ID
	GizmoId string
	// HistoryAndTrainingDisabled is the enable history.
	// 是否开启历史记录，如果下载文件，则必须设置为false,否则会出现文件下载失败
	HistoryAndTrainingDisabled bool
	// Fix is the fix.
	// 是否开启续跑，只跑错误的题。
	Fix bool
	// AccessToken is the access token file.
	// 访问令牌，访问https://gpt4api.shop/consul。复制该令牌
	AccessToken string
	// QPS is the batch size.
	// QPS 设置1s/次 默认是1s/1次
	QPS int
	// NSQ is the nsq config.
	// NSQ配置将数据存储到NSQ队列里，防止丢失.
	NSQ nsq.NSQConfig
	// EnableDownload is the download.
	// 是否开启下载文件
	EnableDownload bool
	// DownloadDir is the download dir.
	// 下载文件文件夹，不传默认当前文件夹
	DownloadDir string
	// DownloadFilePrefix is the download file prefix.
	// 设置下载文件前缀
	DownloadFilePrefix string
	// EnableRDB whether enable rdb.
	// 是否开启RDB缓存.默认会缓存临时数据.
	EnableRDB bool
	// RDBInterval is the rdb interval.
	RDBInterval int
}

func (o *Option) Validate() error {
	if govalidator.IsNull(o.URL) {
		return errors.New("URL is required")
	}

	if govalidator.IsNull(o.In) {
		return errors.New("in is required")
	}

	if govalidator.IsNull(o.Out) {
		return errors.New("out is required")
	}

	if o.Goroutine < 0 {
		return errors.New("goroutine must be greater than 0")
	}

	if o.NSQ.Enable {
		// todo 设置默认参数.
		o.NSQ = nsq.NewNSQConfig()
		if err := o.NSQ.Validate(); err != nil {
			return err
		}
	}

	// get credentials from access token filepath.
	ak, err := gpt4batch.ParseCredentials(o.AccessToken)
	if err != nil {
		return err
	}
	o.AccessToken = ak

	if o.EnableDownload {
		// DownloadDir is null, use the current directory.
		if govalidator.IsNull(o.DownloadDir) {
			localInDir, err := filepath.Abs(o.In)
			if err != nil {
				return err
			}
			o.DownloadDir = filepath.Dir(localInDir)
		} else {
			if _, err := os.Stat(o.DownloadDir); err != nil {
				return err
			}
		}
	}
	return nil
}
