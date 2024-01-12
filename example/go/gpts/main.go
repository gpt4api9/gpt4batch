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

package main

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"gitlab.com/gpt4batch"
	"gitlab.com/gpt4batch/client"
)

func main() {
	cc := client.NewClient()

	// 上传文件
	fileResp, err := cc.Upload(context.Background(), &gpt4batch.UploadRequest{
		Source: &gpt4batch.Source{
			ID:          uuid.New().String(),
			UploadURL:   "https://beta.gpt4api.plus/standard/uploaded",
			Name:        "<随意取一个name>",
			AccessToken: "<填写控制台access_token>",
		},
		UploadPath:     "./example.txt",
		ConversationId: "",
		UploadType:     gpt4batch.MyFiles,
	})
	if err != nil {
		panic(err)
	}

	fmt.Println("<返回数据：>", fileResp)

	// 上传图片
	imageResp, err := cc.Upload(context.Background(), &gpt4batch.UploadRequest{
		Source: &gpt4batch.Source{
			ID:          uuid.New().String(),
			UploadURL:   "https://beta.gpt4api.plus/standard/uploaded",
			Name:        "<随意取一个name>",
			AccessToken: "<填写控制台access_token>",
		},
		UploadPath:     "",
		ConversationId: "",
		UploadType:     gpt4batch.Multimodal,
	})
	if err != nil {
		panic(err)
	}

	fmt.Println("<返回数据：>", imageResp)
}
